package tfprovider

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/davidjspooner/dshttp/pkg/mux"
	"github.com/davidjspooner/dsmatch/pkg/matcher"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	repos matcher.Tree[*repo]
	count int
}

func init() {
	repository.RegisterFactory("tfprovider", &Factory{})
}

func (f *Factory) lookupRepo(w http.ResponseWriter, parsed *parsedRequest) *repo {
	path := parsed.namespace + "/" + parsed.providerName
	leaves := f.repos.FindLeaves([]byte(path))
	if len(leaves) == 0 {
		parsed.logger.Error("repo not found", slog.String("namespace", parsed.namespace), slog.String("name", parsed.providerName))
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	repo := leaves[0].Payload
	return repo
}

func (f *Factory) ConfigureRepo(ctx context.Context, config *repository.Config, mux mux.Mux) error {
	if f.repos.IsEmpty() {

		mux.HandleFunc("GET /.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"providers.v1":"/tf/providers/v1/"}`))
		})
		mux.HandleFunc("GET /tf/providers/v1/{namespace}/{provider}/versions", func(w http.ResponseWriter, r *http.Request) {
			parsed := NewParsedRequest(r)
			repo := f.lookupRepo(w, parsed)
			if repo == nil {
				return
			}
			parsed.ParseVersionOSArch(r)
			repo.HandleProviderVersions(parsed, w, r)
		})

		mux.HandleFunc("GET /tf/providers/v1/{namespace}/{provider}/{version}/download/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			parsed := NewParsedRequest(r)
			repo := f.lookupRepo(w, parsed)
			if repo == nil {
				return
			}
			parsed.ParseVersionOSArch(r)
			repo.Download(parsed, w, r)
		})
		mux.HandleFunc("PUT /tf/providers/v1/{namespace}/{provider}/{version}/download/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			parsed := NewParsedRequest(r)
			repo := f.lookupRepo(w, parsed)
			if repo == nil {
				return
			}
			parsed.ParseVersionOSArch(r)
			repo.Upload(parsed, w, r)
		})
		mux.HandleFunc("DELETE /tf/providers/v1/{namespace}/{provider}/{version}/download/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			parsed := NewParsedRequest(r)
			repo := f.lookupRepo(w, parsed)
			if repo == nil {
				return
			}
			parsed.ParseVersionOSArch(r)
			repo.Delete(parsed, w, r)
		})
	}

	repo, err := newRepo(ctx, config)
	if err != nil {
		return err
	}
	f.count++
	repo.order = f.count
	for _, item := range config.Items {
		seq, err := repository.NewGlob([]byte(item), '/')
		if err != nil {
			return err
		}
		leaf, err := f.repos.Add(seq)
		if err != nil {
			return err
		}
		leaf.Payload = repo
	}

	return nil
}
