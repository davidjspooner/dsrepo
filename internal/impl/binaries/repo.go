package binaries

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/davidjspooner/dsfile/pkg/store"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type repo struct {
	local store.Interface
	order int
}

func newRepo(ctx context.Context, config *repository.Config) (*repo, error) {
	repo := &repo{}
	var err error
	repo.local, err = store.Mount(ctx, config.Local.Path, config.Local.Arguments)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *repo) IsAllowed(parsed *parsedRequest, w http.ResponseWriter, r *http.Request, operation string) bool {
	return true
}

func (repo *repo) List(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "list") {
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) Download(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "get") {
		return
	}

	target := path.Join(parsed.namespace, parsed.filename)
	rFile, err := repo.local.Open(target)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer rFile.Close()
	stat, err := rFile.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+parsed.filename)
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Modified", stat.ModTime().UTC().Format(http.TimeFormat))

	etagged, ok := stat.(store.EntityTagged)
	if ok {
		etag, err := etagged.EntityTag()
		if err != nil {
			w.Header().Set("ETag", etag)
		}
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, rFile)
}

func (repo *repo) Upload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "put") {
		return
	}
	//obs.Logger.Info("provider-upload", slog.String("namespace", key.Namespace), slog.String("name", key.Provider), slog.String("version", key.Version), slog.String("os", key.OS), slog.String("arch", key.Arch))
	defer r.Body.Close()
	buffer := bytes.Buffer{}
	readLength, err := io.Copy(&buffer, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	contentLength := r.Header.Get("Content-Length")
	if contentLength != "" {
		if claimedLength, _ := strconv.Atoi(contentLength); readLength != -1 && int64(claimedLength) != readLength {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	etag := r.Header.Get("ETag")
	if etag == "" {
		hmac := md5.New()
		hmac.Write(buffer.Bytes())
		etag = hex.EncodeToString(hmac.Sum(nil))
	}

	target := path.Join(parsed.namespace, parsed.filename)
	info := store.Info{
		Size:      int64(readLength),
		Mode:      0644,
		EntityTag: etag,
	}

	wFile, err := repo.local.Create(target, info.FileInfo())
	if err != nil {
		parsed.logger.Error("failed to create file", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = wFile.Write(buffer.Bytes())
	if err != nil {
		wFile.Close()
		parsed.logger.Error("failed to start writing file", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = wFile.Close()
	if err != nil {
		parsed.logger.Error("failed to finish writing file", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (repo *repo) Delete(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "delete") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}
