package tfprovider

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/davidjspooner/dsfile/pkg/store"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Repo struct {
	local store.Interface
	order int
}

func newRepo(ctx context.Context, config *repository.Config) (*Repo, error) {
	repo := &Repo{}
	var err error
	repo.local, err = store.Mount(ctx, config.Local.Path, config.Local.Arguments)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *Repo) CheckAccess(parsed *parsedRequest, operation string) bool {
	//TODO: check permissions
	return true
}

func (repo *Repo) HandleProviderVersions(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
	  "versions": [
		{
		  "version": "2.0.0",
		  "protocols": ["4.0", "5.1"],
		  "platforms": [
			{"os": "darwin", "arch": "amd64"},
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm"},
			{"os": "windows", "arch": "amd64"}
		  ]
		},
		{
		  "version": "2.0.1",
		  "protocols": ["5.2"],
		  "platforms": [
			{"os": "darwin", "arch": "amd64"},
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm"},
			{"os": "windows", "arch": "amd64"}
		  ]
		}
	  ]
	}`))
}

func (repo *Repo) HandleProviderDownload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
}

func (repo *Repo) HandleProviderUpload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	//obs.Logger.Info("provider-upload", slog.String("namespace", key.Namespace), slog.String("name", key.Provider), slog.String("version", key.Version), slog.String("os", key.OS), slog.String("arch", key.Arch))
	defer r.Body.Close()
	buffer := bytes.Buffer{}
	panic("implement me")
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
	w.WriteHeader(http.StatusNoContent)
}

func (repo *Repo) HandleProviderDelete(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
