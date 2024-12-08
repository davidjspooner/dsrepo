package repository

import (
	"fmt"
	"strings"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"golang.org/x/exp/maps"
)

type Factory interface {
	ConfigureRepo(config *Config, mux httphandler.Mux) error
}

var factories = make(map[string]Factory)

func RegisterFactory(rType string, factory Factory) {
	factories[rType] = factory
}

func ConfigureRepo(config *Config, mux httphandler.Mux) error {
	factory, ok := factories[config.Type]
	if !ok {
		types := maps.Keys(factories)

		return fmt.Errorf("unknown tree type: %s is not one of %s", config.Type, strings.Join(types, ", "))
	}
	return factory.ConfigureRepo(config, mux)
}
