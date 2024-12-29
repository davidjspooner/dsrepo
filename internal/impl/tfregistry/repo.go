package tfregistry

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"path"
	"slices"
	"strings"

	"github.com/davidjspooner/dsrepo/internal/repository"
)

type repo struct {
	handler *repository.Handler
	order   int
}

func newRepo(ctx context.Context, config *repository.Config) (*repo, error) {
	repo := &repo{}
	var err error
	repo.handler, err = repository.NewHandler(ctx, config)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *repo) IsAllowed(parsed *parsedRequest, w http.ResponseWriter, r *http.Request, operation string) bool {
	//TODO: check permissions
	return true
}

var allowedArchs = []string{"amd64", "arm", "arm64", "386", "ppc64le", "s390x", "mips64", "mips64le", "riscv64"}
var allowedOSs = []string{"darwin", "linux", "windows", "freebsd", "openbsd", "netbsd", "solaris", "dragonfly", "plan9", "aix", "zos"}

func (repo *repo) HandleProviderVersions(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "list") {
		return
	}

	//read the filesystem to get the versions, os and archs

	index := Index{}

	target := path.Join(parsed.namespace, parsed.providerName) + "/"
	err := fs.WalkDir(repo.handler.Local, target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			parts := strings.Split(path, "/")
			arch := parts[len(parts)-1]
			arch = strings.TrimSuffix(arch, ".json")
			os := parts[len(parts)-2]
			version := parts[len(parts)-3]
			if slices.Contains(allowedArchs, arch) && slices.Contains(allowedOSs, os) {
				_ = version

				found := false
				for _, v := range index.Versions {
					if v.Version == version {
						//add the os and arch to the version
						found = true
						v.Platforms = append(v.Platforms, Platform{OS: os, Arch: arch})
					}
				}
				if !found {
					//add the version
					version := Version{Version: version}
					version.Platforms = append(version.Platforms, Platform{OS: os, Arch: arch})
					index.Versions = append(index.Versions, &version)
					//todo read a json to get the protocols

					f, err := repo.handler.Local.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()
					var data map[string]interface{}
					err = json.NewDecoder(f).Decode(&data)
					if err != nil {
						return err
					}
					protocols := data["protocols"].([]any)
					for _, protocol := range protocols {
						version.Protocols = append(version.Protocols, protocol.(string))
					}

				}

				return nil
			}
			return nil
		}
		return nil
	})
	if err != nil {
		http.Error(w, "could not walk", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(index)
}

func (repo *repo) Download(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "get") {
		return
	}
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandleGet(target, parsed.logger, w, r)
}

func (repo *repo) Upload(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "put") {
		return
	}
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandlePut(target, parsed.logger, w, r)
}

func (repo *repo) Delete(parsed *parsedRequest, w http.ResponseWriter, r *http.Request) {
	if !repo.IsAllowed(parsed, w, r, "delete") {
		return
	}
	target := path.Join(parsed.namespace, parsed.providerName, parsed.version, parsed.os, parsed.arch+".json")
	repo.handler.HandleDelete(target, parsed.logger, w, r)
}
