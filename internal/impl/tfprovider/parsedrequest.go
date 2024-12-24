package tfprovider

import (
	"log/slog"
	"net/http"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
)

type parsedRequest struct {
	Namespace string
	Provider  string
	Version   string
	OS        string
	Arch      string

	Logger slog.Logger
}

func NewParsedRequest(r *http.Request) *parsedRequest {
	pr := &parsedRequest{
		Namespace: r.PathValue("namespace"),
		Provider:  r.PathValue("provider"),
	}
	obs, _ := httphandler.GetObservation(r)
	if obs != nil {
		pr.Logger = obs.Logger
	}
	return pr
}

func (pr *parsedRequest) ParseVersionOSArch(r *http.Request) {
	pr.Version = r.PathValue("version")
	pr.OS = r.PathValue("os")
	pr.Arch = r.PathValue("arch")
}
