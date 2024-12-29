package tfregistry

import (
	"log/slog"
	"net/http"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
)

type parsedRequest struct {
	namespace    string
	providerName string
	version      string
	os           string
	arch         string

	logger slog.Logger
}

func NewParsedRequest(r *http.Request) *parsedRequest {
	pr := &parsedRequest{
		namespace:    r.PathValue("namespace"),
		providerName: r.PathValue("provider"),
	}
	obs, _ := httphandler.GetObservation(r)
	if obs != nil {
		pr.logger = obs.Logger
	}
	return pr
}

func (pr *parsedRequest) ParseVersionOSArch(r *http.Request) {
	pr.version = r.PathValue("version")
	pr.os = r.PathValue("os")
	pr.arch = r.PathValue("arch")
}
