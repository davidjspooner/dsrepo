package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidjspooner/dshttp/pkg/mux"
	"golang.org/x/exp/maps"
)

type Router interface {
	NewRepo(ctx context.Context, config *Config) error
	SetupRoutes(mux mux.Mux) error
}

var routers = make(map[string]Router)

func RegisterRouter(rType string, router Router) {
	routers[rType] = router
}

func SetupRoutes(mux mux.Mux) error {
	for _, router := range routers {
		if err := router.SetupRoutes(mux); err != nil {
			return err
		}
	}
	return nil
}

func NewRepo(ctx context.Context, config *Config) error {
	router, ok := routers[config.Type]
	if !ok {
		types := maps.Keys(routers)
		return fmt.Errorf("unknown tree type: %s is not one of %s", config.Type, strings.Join(types, ", "))
	}
	return router.NewRepo(ctx, config)
}
