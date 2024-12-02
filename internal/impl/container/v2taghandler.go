package container

import "net/http"

type v2TagHandler struct {
	factory *Factory
}

func (handler *v2TagHandler) get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not implemented v2TagHandler.get", http.StatusNotImplemented)
}
