package repository

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/davidjspooner/dsfile/pkg/store"
)

type Handler struct {
	Local    store.Interface
	Upstream *url.URL
}

func NewHandler(ctx context.Context, config *Config) (*Handler, error) {
	handler := &Handler{}
	var err error
	handler.Local, err = store.Mount(ctx, config.Local.Path, config.Local.Arguments)
	if err != nil {
		return nil, err
	}

	if config.Upstream.Url != "" {
		handler.Upstream, err = url.Parse(config.Upstream.Url)
		if err != nil {
			return nil, err
		}
	}

	return handler, nil
}

func (handler *Handler) HandleGet(target string, logger slog.Logger, w http.ResponseWriter, r *http.Request) error {
	rFile, err := handler.Local.Open(target)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error("file:open", slog.String("target", target), slog.String("error", err.Error()))
		return err
	}
	defer rFile.Close()
	stat, err := rFile.Stat()
	if err != nil {
		logger.Error("file:stat", slog.String("target", target), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(target))
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
	_, err = io.Copy(w, rFile)
	if err != nil {
		logger.Error("file:read", slog.String("target", target), slog.String("error", err.Error()))
	}
	return err
}

func (handler *Handler) HandlePut(target string, logger slog.Logger, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	buffer := bytes.Buffer{}
	readLength, err := io.Copy(&buffer, r.Body)
	if err != nil {
		logger.Error("content:read", slog.String("target", target), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	contentLength := r.Header.Get("Content-Length")
	if contentLength != "" {
		if claimedLength, _ := strconv.Atoi(contentLength); readLength != -1 && int64(claimedLength) != readLength {
			err := fmt.Errorf("content length mismatch: %d != %d", claimedLength, readLength)
			logger.Error("content:validate", slog.String("target", target), slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	etag := r.Header.Get("ETag")
	if etag == "" {
		hmac := md5.New()
		hmac.Write(buffer.Bytes())
		etag = hex.EncodeToString(hmac.Sum(nil))
	}

	info := store.Info{
		Size:      int64(readLength),
		Mode:      0644,
		EntityTag: etag,
	}

	wFile, err := handler.Local.Create(target, info.FileInfo())
	if err != nil {
		logger.Error("file:create", slog.String("target", target), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = wFile.Write(buffer.Bytes())
	if err != nil {
		logger.Error("file:write start", slog.String("target", target), slog.String("error", err.Error()))
		wFile.Close()
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	err = wFile.Close()
	if err != nil {
		logger.Error("file:write finish", slog.String("target", target), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (handler *Handler) HandleDelete(target string, logger slog.Logger, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotImplemented)
	err := fmt.Errorf("not implemented")
	logger.Error("file:deletion", slog.String("target", target), slog.String("error", err.Error()))
	return err
}
