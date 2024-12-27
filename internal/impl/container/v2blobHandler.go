package container

import "net/http"

type v2BlobHandler struct {
	router *Router
}

func (handler *v2BlobHandler) get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2BlobHandler.get", http.StatusNotImplemented)
}
func (handler *v2BlobHandler) put(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2BlobHandler) put", http.StatusNotImplemented)
}
func (handler *v2BlobHandler) post(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2BlobHandler.post", http.StatusNotImplemented)
}
func (handler *v2BlobHandler) patch(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2BlobHandler.patch", http.StatusNotImplemented)
}
func (handler *v2BlobHandler) delete(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2BlobHandler.delete", http.StatusNotImplemented)
}
