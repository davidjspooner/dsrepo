package container

import (
	"context"

	"github.com/davidjspooner/dshttp/pkg/mux"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Router struct {
	v2BlobHandler     *v2BlobHandler
	v2CatalogHandler  *v2CatalogHandler
	v2ManifestHandler *v2ManifestHandler
	v2TagHandler      *v2TagHandler
	v2RootHandler     *v2RootHandler
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

func (router *Router) NewRepo(ctc context.Context, config *repository.Config) error {
	if router.v2BlobHandler == nil {
		router.v2BlobHandler = &v2BlobHandler{router: router}
		router.v2CatalogHandler = &v2CatalogHandler{router: router}
		router.v2ManifestHandler = &v2ManifestHandler{router: router}
		router.v2TagHandler = &v2TagHandler{router: router}
		router.v2RootHandler = &v2RootHandler{router: router}
	}
	return nil
}

func (router *Router) SetupRoutes(mux mux.Mux) error {
	if router.v2BlobHandler == nil {
		return nil
	}
	mux.HandleFunc("GET /v2/{$}", router.v2RootHandler.get)
	mux.HandleFunc("GET /v2/_catalog", router.v2CatalogHandler.get)
	mux.HandleFunc("GET /v2/{name...}/blobs/{digest}", router.v2BlobHandler.get)
	mux.HandleFunc("POST /v2/{name...}/blobs/uploads/{$}", router.v2BlobHandler.post)
	mux.HandleFunc("PATCH /v2/{name...}/blobs/uploads/{reference}", router.v2BlobHandler.patch)
	mux.HandleFunc("PUT /v2/{name...}/blobs/uploads/{reference}", router.v2BlobHandler.put)
	mux.HandleFunc("DELETE /v2/{name...}/blobs/{$}", router.v2BlobHandler.delete)
	mux.HandleFunc("GET /v2/{name...}/manifests/{reference}", router.v2ManifestHandler.get)
	mux.HandleFunc("PUT /v2/{name...}/manifests/{reference}", router.v2ManifestHandler.put)
	mux.HandleFunc("DELETE /v2/{name...}/manifests/{reference}", router.v2ManifestHandler.delete)
	mux.HandleFunc("GET /v2/{name...}/tags/list", router.v2TagHandler.get)
	return nil
}
