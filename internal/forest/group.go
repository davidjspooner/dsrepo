package forest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dshttp/pkg/logevent"
	dsmiddleware "github.com/davidjspooner/dshttp/pkg/middleware"
	"github.com/davidjspooner/dsrepo/internal/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Group struct {
	config      Config
	ctx         context.Context
	log         *slog.Logger
	servers     map[string]*server
	promHandler http.Handler
}

type Option func(*Group) error

func NewServerGroup(options ...Option) (*Group, error) {
	loghandler := logevent.NewHandler(&slog.HandlerOptions{})
	group := &Group{
		log:         slog.New(loghandler),
		ctx:         context.Background(),
		servers:     make(map[string]*server),
		promHandler: promhttp.Handler(),
	}
	for _, option := range options {
		if err := option(group); err != nil {
			return nil, err
		}
	}

	err := group.initServers()
	if err != nil {
		return nil, err
	}
	return group, nil
}

func WithLogger(log *slog.Logger) Option {
	return func(s *Group) error {
		s.log = log
		s.ctx = logevent.WithLogger(s.ctx, s.log)
		return nil
	}
}

func (g *Group) initServers() error {

	repositories := make(map[string]*repository.Config)

	for _, repo := range g.config.Repositories {
		if repo.Name == "" {
			return fmt.Errorf("repository name is required")
		}
		if repositories[repo.Name] != nil {
			return fmt.Errorf("duplicate repository name: %s", repo.Name)
		}
		repositories[repo.Name] = repo
	}

	for _, listener := range g.config.Listeners {
		if listener.Name == "" {
			return fmt.Errorf("listener name is required")
		}
		if g.servers[listener.Name] != nil {
			return fmt.Errorf("duplicate listener name: %s", listener.Name)
		}

		pipeline := httphandler.MiddlewarePipeline{
			&dsmiddleware.Observer{
				BeforeRequest: func(ctx context.Context, req *http.Request, observed *dsmiddleware.Observation) {
					inflightRequests.WithLabelValues(listener.Name, req.Method).Inc()
				},
				AfterRequest: func(ctx context.Context, req *http.Request, observed *dsmiddleware.Observation) {
					duration := observed.Response.Duration.Seconds()
					args := []any{
						slog.Group("req",
							slog.Uint64("id", observed.Request.ID),
							slog.String("method", req.Method),
							slog.String("path", req.URL.Path),
							slog.Int("bytes", observed.Request.Body.Length),
							slog.String("remote", req.RemoteAddr),
						),
						slog.Group("res",
							slog.Int("status", observed.Response.Status),
							slog.Int("bytes", observed.Response.Body.Length),
							slog.String("duration", fmt.Sprintf("%.3f", duration)),
						),
					}
					g.log.Info("handled", args...)
					inflightRequests.WithLabelValues(listener.Name, req.Method).Dec()
					requestDuration.WithLabelValues(listener.Name, req.Method).Observe(duration)
					bytesWritten.WithLabelValues(listener.Name, req.Method).Add(float64(observed.Response.Body.Length))
				},
			},
			&dsmiddleware.Recovery{},
			&dsmiddleware.HeadMethodHelper{},
		}

		server := &server{
			listenerConfig: listener,
			group:          g,
			log:            g.log.WithGroup(listener.Name),
			mux:            httphandler.NewServeMux(pipeline),
		}
		g.servers[listener.Name] = server
		for _, repoName := range listener.Expose {
			config := repositories[repoName]
			err := repository.ConfigureRepo(config, server.mux)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Group) ListenAndServe() {
	group := sync.WaitGroup{}
	for _, s := range g.servers {
		group.Add(1)
		go func(s *server) {
			defer group.Done()
			err := s.ListenAndServe()
			if err != nil {
				g.log.Error("Failed to start server", slog.String("error", err.Error()))
			}
		}(s)

	}
	group.Wait()
}

func (g *Group) metricsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodHead:
		promhttp.Handler().ServeHTTP(w, req)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Group) healthHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodHead:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
