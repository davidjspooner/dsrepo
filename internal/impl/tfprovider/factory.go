package tfprovider

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/davidjspooner/dshttp/pkg/mux"
	"github.com/davidjspooner/dsmatch/pkg/matcher"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type Factory struct {
	repos matcher.Tree[*Repo]
}

func init() {
	repository.RegisterFactory("tfprovider", &Factory{})
}

func (f *Factory) ConfigureRepo(ctx context.Context, config *repository.Config, mux mux.Mux) error {
	if f.repos.Empty() {

		mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"providers.v1":"/tf/providers/v1/"}`))
		})
		mux.HandleFunc("/tf/providers/v1/{namespace}/{provider}/versions", func(w http.ResponseWriter, r *http.Request) {
			var key key
			key.Namespace = r.PathValue("namespace")
			key.Provider = r.PathValue("provider")
			f.HandleProviderVersions(&key, w, r)
		})

		mux.HandleFunc("GET /tf/providers/v1/{namespace}/{provider}/{version}/download/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			var key key
			key.Namespace = r.PathValue("namespace")
			key.Provider = r.PathValue("provider")
			key.Version = r.PathValue("version")
			key.OS = r.PathValue("os")
			key.Arch = r.PathValue("arch")
			f.HandleProviderDownload(&key, w, r)
		})
		mux.HandleFunc("PUT /tf/providers/v1/{namespace}/{provider}/{version}/upload/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			var key key
			key.Namespace = r.PathValue("namespace")
			key.Provider = r.PathValue("provider")
			key.Version = r.PathValue("version")
			key.OS = r.PathValue("os")
			key.Arch = r.PathValue("arch")
			f.HandleProviderUpload(&key, w, r)
		})
		mux.HandleFunc("DELETE /tf/providers/v1/{namespace}/{provider}/{version}/upload/{os}/{arch}", func(w http.ResponseWriter, r *http.Request) {
			var key key
			key.Namespace = r.PathValue("namespace")
			key.Provider = r.PathValue("provider")
			key.Version = r.PathValue("version")
			key.OS = r.PathValue("os")
			key.Arch = r.PathValue("arch")
			f.HandleProviderDelete(&key, w, r)
		})

	}

	repo, err := newRepo(ctx, config)
	if err != nil {
		return err
	}
	_ = repo //todo parse items list into the tree

	return nil
}

func (repo *Repo) HandleProviderVersions(key *key, w http.ResponseWriter, r *http.Request) {

	//TODO: check permissions
	slog.Info("provider-versions", slog.String("namespace", key.Namespace), slog.String("provider", key.Provider))
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
}

func (f *Factory) HandleProviderVersions(key *key, w http.ResponseWriter, r *http.Request) {
	slog.Info("provider-versions", slog.String("namespace", key.Namespace), slog.String("name", key.Provider))
}

func (f *Factory) HandleProviderDownload(key *key, w http.ResponseWriter, r *http.Request) {
	slog.Info("provider-download", slog.String("namespace", key.Namespace), slog.String("name", key.Provider), slog.String("version", key.Version), slog.String("os", key.OS), slog.String("arch", key.Arch))
}

func (f *Factory) HandleProviderUpload(key *key, w http.ResponseWriter, r *http.Request) {
	slog.Info("provider-upload", slog.String("namespace", key.Namespace), slog.String("name", key.Provider), slog.String("version", key.Version), slog.String("os", key.OS), slog.String("arch", key.Arch))
	defer r.Body.Close()
	io.Copy(io.Discard, r.Body)
	w.WriteHeader(http.StatusNoContent)
}

func (f *Factory) HandleProviderDelete(key *key, w http.ResponseWriter, r *http.Request) {
	slog.Info("provider-delete", slog.String("namespace", key.Namespace), slog.String("name", key.Provider), slog.String("version", key.Version), slog.String("os", key.OS), slog.String("arch", key.Arch))
	w.WriteHeader(http.StatusNoContent)
}
