package forest

import (
	"os"

	"github.com/davidjspooner/dsrepo/internal/repository"
	"gopkg.in/yaml.v3"
)

type ListenerConfig struct {
	Port     int
	CertFile string
	KeyFile  string
}

type Config struct {
	Listener     ListenerConfig
	Repositories []*repository.Config
}

func WithConfigFile(path string) Option {
	return func(s *Server) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		d := yaml.NewDecoder(f)
		d.KnownFields(true)
		err = d.Decode(&s.config)
		return err
	}
}
