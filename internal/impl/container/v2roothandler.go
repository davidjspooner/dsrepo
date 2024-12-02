package container

import "net/http"

type v2RootHandler struct {
	factory *Factory
}

func (handler *v2RootHandler) get(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
	w.WriteHeader(http.StatusOK)

}
