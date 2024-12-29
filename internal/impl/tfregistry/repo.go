package tfregistry

import (
	"context"
	"net/http"
	"path"

	"github.com/davidjspooner/dsrepo/internal/repository"
)

type repo struct {
	handler *repository.Handler
	order   int
}

func newRepo(ctx context.Context, config *repository.Config) (*repo, error) {
	repo := &repo{}
	var err error
	repo.handler, err = repository.NewHandler(ctx, config)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *repo) IsAllowed(parsed *parsedRequest, w http.ResponseWriter, r *http.Request, operation string) bool {
	//TODO: check permissions
	return true
}

func (repo *repo) HandleProviderVersions(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "list") {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{
	  "versions": [
		{
		  "version": "0.2.0",
		  "protocols": ["4.0", "5.1"],
		  "platforms": [
			{"os": "darwin", "arch": "amd64"},
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm"},
			{"os": "windows", "arch": "amd64"}
		  ]
		},
		{
		  "version": "0.2.1",
		  "protocols": ["6.0"],
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
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandleGet(target, parsed.logger, w, r)
}

func (repo *repo) Upload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "put") {
		return
	}
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandlePut(target, parsed.logger, w, r)
}

func (repo *repo) Delete(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "delete") {
		return
	}
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandleDelete(target, parsed.logger, w, r)
}
