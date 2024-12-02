package forest

import (
	"os"

	"github.com/davidjspooner/dsrepo/internal/repository"
	"gopkg.in/yaml.v3"
)

type ListenerConfig struct {
	Name     string
	Port     int
	CertFile string
	KeyFile  string
	Expose   []string
}

type Config struct {
	Listeners    []*ListenerConfig
	Repositories []*repository.Config
}

func WithConfigFile(path string) Option {
	return func(s *Group) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		d := yaml.NewDecoder(f)
		err = d.Decode(&s.config)
		return err
	}
}
