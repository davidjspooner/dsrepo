package tfprovider

import (
	"context"

	"github.com/davidjspooner/dsfile/pkg/store"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Repo struct {
	local store.Interface
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
