package binaries

import (
	"context"
	"log/slog"
	"net/http"
	"path"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dshttp/pkg/mux"
	"github.com/davidjspooner/dsmatch/pkg/matcher"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	repos matcher.Tree[*repo]
	count int
}

type parsedRequest struct {
	namespace string
	filename  string
	repo      *repo
	logger    slog.Logger
}

func init() {
	repository.RegisterFactory("directory", &Factory{})
}

func (f *Factory) GetRepo(filename string) []matcher.Leaf[*repo] {
	leaves := f.repos.FindLeaves([]byte(filename))
	for n := 0; n < len(leaves); {
		if leaves[n].Payload == nil {
			leaves = append(leaves[:n], leaves[n+1:]...)
		} else {
			n++
		}
	}
	return leaves
}

func (f *Factory) ParseRequest(w http.ResponseWriter, r *http.Request) *parsedRequest {
	namespace := r.PathValue("filename")
	filename := ""

	leaves := f.GetRepo(filename)
	if len(leaves) == 0 {
		namespace, filename = path.Split(namespace)
		leaves = f.GetRepo(namespace)
	}
	if len(leaves) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	if len(leaves) > 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("ambiguous request"))
		return nil
	}

	pr := &parsedRequest{
		namespace: namespace,
		filename:  filename,
		repo:      leaves[0].Payload,
	}

	obs, _ := httphandler.GetObservation(r)
	if obs != nil {
		pr.logger = obs.Logger
	}
	return pr
}

func (f *Factory) ConfigureRepo(ctx context.Context, config *repository.Config, mux mux.Mux) error {
	if f.repos.IsEmpty() {
		mux.HandleFunc("GET /binary/{filename...}", func(w http.ResponseWriter, r *http.Request) {
			parsed := f.ParseRequest(w, r)
			if parsed == nil {
				return
			}
			if parsed.filename == "" {
				parsed.repo.List(parsed, w, r)
				return
			}
			parsed.repo.Download(parsed, w, r)
		})
		mux.HandleFunc("PUT /binary/{filename...}", func(w http.ResponseWriter, r *http.Request) {
			parsed := f.ParseRequest(w, r)
			if parsed == nil {
				return
			}
			parsed.repo.Upload(parsed, w, r)
		})
		mux.HandleFunc("DELETE /binary/{filename...}", func(w http.ResponseWriter, r *http.Request) {
			parsed := f.ParseRequest(w, r)
			if parsed == nil {
				return
			}
			parsed.repo.Delete(parsed, w, r)
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
