package container

import "net/http"

type v2CatalogHandler struct {
	factory *Factory
}

func (handler *v2CatalogHandler) get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2CatalogHandler.get", http.StatusNotImplemented)
}
