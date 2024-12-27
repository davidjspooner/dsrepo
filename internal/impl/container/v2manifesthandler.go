package container

import "net/http"

type v2ManifestHandler struct {
	router *Router
}

func (handler *v2ManifestHandler) get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2ManifestHandler.get", http.StatusNotImplemented)
}
func (handler *v2ManifestHandler) put(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2ManifestHandler.put", http.StatusNotImplemented)
}
func (handler *v2ManifestHandler) delete(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2ManifestHandler.delete", http.StatusNotImplemented)
}
