package tfprovider

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	lastMux httphandler.Mux
	cache   *repository.Cache[Provider]
}

func init() {
	repository.RegisterFactory("tfprovider", &Factory{
		cache: repository.NewCacheMap[Provider](100),
	})
}

func (f *Factory) ConfigureRepo(config *repository.Config, mux httphandler.Mux) error {
	if f.lastMux != mux {
		f.lastMux = mux

		mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
			//TODO: check permissions
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"providers.v1":"/tf/providers/v1/"}`))
		})
		mux.HandleFunc("/tf/providers/v1/{namespace}/{provider}/versions", func(w http.ResponseWriter, r *http.Request) {
			//TODO: check permissions
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
			name := r.PathValue("provider")
			version := r.PathValue("version")
			os := r.PathValue("os")
			arch := r.PathValue("arch")
			key := namespace + "/" + name

			switch r.Method {
			case "GET":
				slog.Info("provider-download", slog.String("namespace", namespace), slog.String("name", name), slog.String("version", version), slog.String("os", os), slog.String("arch", arch))
				//TODO: check permissions
				f.cache.Use(key, func(key string, cached *Provider, age time.Duration) (*Provider, bool) {
					//TODO: return the actual binary
					return nil, true
				})

				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write([]byte("fake-binary"))
			case "DELETE":
				slog.Info("provider-delete", slog.String("namespace", namespace), slog.String("name", name), slog.String("version", version), slog.String("os", os), slog.String("arch", arch))
				//TODO: check permissions

				f.cache.Use(key, func(key string, cached *Provider, age time.Duration) (*Provider, bool) {
					//TODO: delete the binary
					return nil, true
				})
				w.WriteHeader(http.StatusNoContent)
			case "PUT":
				slog.Info("provider-upload", slog.String("namespace", namespace), slog.String("name", name), slog.String("version", version), slog.String("os", os), slog.String("arch", arch))
				//TODO: check permissions
				f.cache.Use(key, func(key string, cached *Provider, age time.Duration) (*Provider, bool) {
					//todo: save the binary
					return nil, true
				})

				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})

	}
	//todo use config to create a a repo

	return nil
}

var MaxCacheAge = 2 * time.Minute
