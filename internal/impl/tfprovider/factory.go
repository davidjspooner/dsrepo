package tfprovider

import (
	"log/slog"
	"net/http"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	lastMux *httphandler.ServeMux
}

func init() {
	repository.RegisterFactory("tfprovider", &Factory{})
}

func (f *Factory) ConfigureRepo(config *repository.Config, mux *httphandler.ServeMux) error {
	if f.lastMux != mux {
		f.lastMux = mux

		mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"providers.v1":"/tf/providers/v1/"}`))
		})
		mux.HandleFunc("/tf/providers/v1/{namespace}/{provider}/versions", func(w http.ResponseWriter, r *http.Request) {
			namespace := r.PathValue("namespace")
			provider := r.PathValue("provider")
			slog.Info("provider-versions", slog.String("namespace", namespace), slog.String("provider", provider))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
	  "versions": [
		{
		  "version": "2.0.0",
		  "protocols": ["4.0", "5.1"],
		  "platforms": [
			{"os": "darwin", "arch": "amd64"},
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm"},
			{"os": "windows", "arch": "amd64"}
		  ]
		},
		{
		  "version": "2.0.1",
		  "protocols": ["5.2"],
		  "platforms": [
			{"os": "darwin", "arch": "amd64"},
			{"os": "linux", "arch": "amd64"},
			{"os": "linux", "arch": "arm"},
			{"os": "windows", "arch": "amd64"}
		  ]
		}
	  ]
	}`))
		})

		mux.HandleFunc("/tf/providers/v1/{namespace}/{provider}/{version}/download/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			namespace := r.PathValue("namespace")
			provider := r.PathValue("provider")
			version := r.PathValue("version")
			os := r.PathValue("os")
			arch := r.PathValue("arch")
			slog.Info("provider-download", slog.String("namespace", namespace), slog.String("provider", provider), slog.String("version", version), slog.String("os", os), slog.String("arch", arch))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte("fake-binary"))
		})

	}
	//todo use config to create a a repo

	return nil
}
