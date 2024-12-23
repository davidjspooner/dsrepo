package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidjspooner/dshttp/pkg/mux"
	"golang.org/x/exp/maps"
)

type Factory interface {
	ConfigureRepo(ctx context.Context, config *Config, mux mux.Mux) error
}

var factories = make(map[string]Factory)

func RegisterFactory(rType string, factory Factory) {
	factories[rType] = factory
}

func ConfigureRepo(ctx context.Context, config *Config, mux mux.Mux) error {
	factory, ok := factories[config.Type]
	if !ok {
		types := maps.Keys(factories)
		return fmt.Errorf("unknown tree type: %s is not one of %s", config.Type, strings.Join(types, ", "))
	}
	return factory.ConfigureRepo(ctx, config, mux)
}
