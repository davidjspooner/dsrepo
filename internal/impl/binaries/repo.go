package binaries

import (
	"context"
	"log/slog"
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
	repo.handler.HandleGet(target, parsed.logger, w, r)
}

func (repo *repo) Upload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "put") {
		return
	}
	target := path.Join(parsed.namespace, parsed.filename)
	err := repo.handler.HandlePut(target, parsed.logger, w, r)
	if err != nil {
		parsed.logger.Error("failed to put file", slog.String("error", err.Error()))
	}
}

func (repo *repo) Delete(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "delete") {
		return
	}
	target := path.Join(parsed.namespace, parsed.filename)
	err := repo.handler.HandleDelete(target, parsed.logger, w, r)
	if err != nil {
		parsed.logger.Error("failed to delete file", slog.String("error", err.Error()))
	}
}
