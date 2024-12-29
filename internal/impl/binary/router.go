package binary

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

type Router struct {
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
	repository.RegisterRouter("binary", &Router{})
}

func (router *Router) GetRepo(filename string) []matcher.Leaf[*repo] {
	leaves := router.repos.FindLeaves([]byte(filename))
	for n := 0; n < len(leaves); {
		if leaves[n].Payload == nil {
			leaves = append(leaves[:n], leaves[n+1:]...)
		} else {
			n++
		}
	}
	return leaves
}

func (router *Router) ParseRequest(w http.ResponseWriter, r *http.Request) *parsedRequest {
	namespace := r.PathValue("filename")
	filename := ""

	leaves := router.GetRepo(filename)
	if len(leaves) == 0 {
		namespace, filename = path.Split(namespace)
		leaves = router.GetRepo(namespace)
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

func (router *Router) SetupRoutes(mux mux.Mux) error {
	mux.HandleFunc("GET /binary/{filename...}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
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
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.Upload(parsed, w, r)
	})
	mux.HandleFunc("DELETE /binary/{filename...}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.Delete(parsed, w, r)
	})
	return nil
}

func (router *Router) NewRepo(ctx context.Context, config *repository.Config) error {
	repo, err := newRepo(ctx, config)
	if err != nil {
		return err
	}
	router.count++
	repo.order = router.count
	for _, item := range config.Items {
		seq, err := repository.NewGlob([]byte(item), '/')
		if err != nil {
			return err
		}
		leaf, err := router.repos.Add(seq)
		if err != nil {
			return err
		}
		leaf.Payload = repo
	}

	return nil
}
