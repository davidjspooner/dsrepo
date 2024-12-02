package container

import (
	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	lastMux *httphandler.ServeMux

	v2BlobHandler     *v2BlobHandler
	v2CatalogHandler  *v2CatalogHandler
	v2ManifestHandler *v2ManifestHandler
	v2TagHandler      *v2TagHandler
	v2RootHandler     *v2RootHandler
}

func init() {
	repository.RegisterFactory("container", &Factory{})
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

func (f *Factory) ConfigureRepo(config *repository.Config, mux *httphandler.ServeMux) error {
	if f.v2BlobHandler == nil {
		f.v2BlobHandler = &v2BlobHandler{factory: f}
		f.v2CatalogHandler = &v2CatalogHandler{factory: f}
		f.v2ManifestHandler = &v2ManifestHandler{factory: f}
		f.v2TagHandler = &v2TagHandler{factory: f}
		f.v2RootHandler = &v2RootHandler{factory: f}
	}
	if f.lastMux != mux {
		f.lastMux = mux
		mux.HandleFunc("GET /v2/", f.v2RootHandler.get)
		mux.HandleFunc("GET /v2/_catalog", f.v2CatalogHandler.get)
		mux.HandleFunc("GET /v2/{name...}/blobs/{digest}", f.v2BlobHandler.get)
		mux.HandleFunc("POST /v2/{name...}/blobs/uploads/", f.v2BlobHandler.post)
		mux.HandleFunc("PATCH /v2/{name...}/blobs/uploads/{reference}", f.v2BlobHandler.patch)
		mux.HandleFunc("PUT /v2/{name...}/blobs/uploads/{reference}", f.v2BlobHandler.put)
		mux.HandleFunc("DELETE /v2/{name...}/blobs/", f.v2BlobHandler.delete)
		mux.HandleFunc("GET /v2/{name...}/manifests/{reference}", f.v2ManifestHandler.get)
		mux.HandleFunc("PUT /v2/{name...}/manifests/{reference}", f.v2ManifestHandler.put)
		mux.HandleFunc("DELETE /v2/{name...}/manifests/{reference}", f.v2ManifestHandler.delete)
		mux.HandleFunc("GET /v2/{name...}/tags/list", f.v2TagHandler.get)
	}
	return nil
}
