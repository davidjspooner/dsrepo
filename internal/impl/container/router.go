package container

import (
	"context"
	"log/slog"
	"net/http"

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
	name      string
	digest    string
	reference string
	repo      *repo
	logger    slog.Logger
}

func init() {
	repository.RegisterRouter("container", &Router{})
}

/*
GET /v2/
GET /v2/_catalog
GET /v2/<name>/blobs/<digest>
HEAD /v2/<name>/blobs/<digest>
POST /v2/<name>/blobs/uploads/
PATCH /v2/<name>/blobs/uploads/<reference>
DELETE /v2/<name>/blobs/
PUT /v2/<name>/blobs/uploads/<reference>
GET /v2/<name>/manifests/<reference>
HEAD /v2/<name>/manifests/<reference>
PUT /v2/<name>/manifests/<reference>
DELETE /v2/<name>/manifests/<reference>
GET /v2/<name>/tags/list
*/

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

func (router *Router) ParseRequest(w http.ResponseWriter, r *http.Request) *parsedRequest {
	parsed := &parsedRequest{}
	parsed.name = r.PathValue("name")
	parsed.digest = r.PathValue("digest")
	parsed.reference = r.PathValue("reference")
	leaves := router.repos.FindLeaves([]byte(parsed.name))
	if len(leaves) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	if len(leaves) > 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("ambiguous request"))
		return nil
	}
	parsed.repo = leaves[0].Payload
	obs, _ := httphandler.GetObservation(r)
	if obs != nil {
		parsed.logger = obs.Logger
	}
	return parsed
}

func (router *Router) SetupRoutes(aMux mux.Mux) error {
	if router.repos.IsEmpty() {
		return nil
	}
	aMux.HandleFunc("GET /v2/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
		w.WriteHeader(http.StatusOK)
	})
	aMux.HandleFunc("GET /v2/_catalog", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	aMux.HandleFunc("GET /v2/{name...}/blobs/{digest}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.getBlobByDigest(parsed, w, r)
	})
	aMux.HandleFunc("POST /v2/{name...}/blobs/uploads/{$}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.uploadBlob(parsed, w, r)
	})
	aMux.HandleFunc("PATCH /v2/{name...}/blobs/uploads/{reference}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.updateBlob(parsed, w, r)
	})
	aMux.HandleFunc("PUT /v2/{name...}/blobs/uploads/{reference}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.updateBlob(parsed, w, r)
	})
	aMux.HandleFunc("DELETE /v2/{name...}/blobs/{$}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.deleteBlob(parsed, w, r)
	})
	aMux.HandleFunc("GET /v2/{name...}/manifests/{reference}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.getManifest(parsed, w, r)
	})
	aMux.HandleFunc("PUT /v2/{name...}/manifests/{reference}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.putManifest(parsed, w, r)
	})
	aMux.HandleFunc("DELETE /v2/{name...}/manifests/{reference}", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.deleteManifest(parsed, w, r)
	})
	aMux.HandleFunc("GET /v2/{name...}/tags/list", func(w http.ResponseWriter, r *http.Request) {
		parsed := router.ParseRequest(w, r)
		if parsed == nil {
			return
		}
		parsed.repo.getTags(parsed, w, r)
	})

	// sm, _ := aMux.(*mux.ServeMux)
	// if sm != nil {
	// 	sm.WriteDebug(os.Stdout, 0)
	// }

	return nil
}
