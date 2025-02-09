package container

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/davidjspooner/dsfile/pkg/store"
	"github.com/davidjspooner/dshttp/pkg/httpclient"
	"github.com/davidjspooner/dshttp/pkg/middleware"
	"github.com/davidjspooner/dsrepo/internal/repository"
)

type repo struct {
	handler *repository.Handler
	order   int
	client  httpclient.Interface
}

func newRepo(ctx context.Context, config *repository.Config) (*repo, error) {
	repo := &repo{}
	var err error
	repo.handler, err = repository.NewHandler(ctx, config)
	if err != nil {
		return nil, err
	}

	if repo.handler.Upstream != nil {
		repo.client = httpclient.NewClient(http.DefaultClient, &middleware.BearerAuthenticator{})
	}

	return repo, nil
}

func (repo *repo) IsAllowed(parsed *parsedRequest, w http.ResponseWriter, r *http.Request, operation string) bool {
	return true
}

func (repo *repo) getBlobByDigest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "GET") {
		return
	}
	if repo.handler.Local != nil {
		path := parsed.name + "/blob/" + parsed.digest
		if repo.handler.LocalFileExists(path) {
			repo.handler.HandleLocalGet(path, parsed.logger, w, r)
		}
		if repo.handler.Upstream != nil {
			//fetch and cache the blob
			brw := bufferedResponseWriter{}
			repo.ProxyUpstream(parsed, &brw, r)
			if brw.status == http.StatusOK {
				//store the blob
				etag := r.Header.Get("ETag")
				if etag == "" {
					hmac := md5.New()
					hmac.Write(brw.body.Bytes())
					etag = hex.EncodeToString(hmac.Sum(nil))
				}

				info := store.Info{
					Size:      int64(brw.body.Len()),
					Mode:      0644,
					EntityTag: etag,
				}

				rFile, err := repo.handler.Local.Create(path, info.FileInfo())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer rFile.Close()
				_, err = io.Copy(rFile, &brw.body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			for k, v := range brw.headers {
				w.Header()[k] = v
			}
			w.WriteHeader(brw.status)
			io.Copy(w, &brw.body)
			return
		}
	}

	if repo.handler.Upstream != nil {
		repo.ProxyUpstream(parsed, w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) uploadBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) updateBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) deleteBlob(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "DELETE") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) getManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "GET") {
		return
	}
	if repo.handler.Upstream != nil {
		repo.ProxyUpstream(parsed, w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) putManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "PUT") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) deleteManifest(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "DELETE") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) getTags(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "LIST") {
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (repo *repo) ProxyUpstream(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if repo.handler.Upstream == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	proxyRequest, err := http.NewRequest(r.Method, repo.handler.Upstream.String()+r.URL.Path, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	proxyRequest.Header = r.Header
	proxyRequest.Header.Set("Host", r.Host)

	response, err := repo.client.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	wh := w.Header()
	for k, v := range response.Header {
		wh[k] = v
	}
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)

}
