package container

import (
	"context"
	"net/http"

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
	return true
}

func (repo *repo) getBlobByDigest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "GET") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) uploadBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) updateBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) deleteBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "DELETE") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) getManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "GET") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) putManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) deleteManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "DELETE") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) getTags(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "LIST") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}
