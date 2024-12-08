package forest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/davidjspooner/dshttp/pkg/logevent"
	"github.com/davidjspooner/dshttp/pkg/middleware"
	"github.com/davidjspooner/dsrepo/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var inflightRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "http_inflight_requests",
	Help: "Number of inflight requests",
}, []string{"method"})

var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "HTTP request duration",
	Buckets: prometheus.DefBuckets,
}, []string{"method"})

var bytesWritten = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_size_bytes",
		Help: "Total size of the response sent, in bytes.",
	},
	[]string{"method"},
)

type Server struct {
	config Config
	ctx    context.Context
	log    *slog.Logger
	mux    *httphandler.ServeMux
}

type Option func(*Server) error

func NewServer(options ...Option) (*Server, error) {
	loghandler := logevent.NewHandler(&slog.HandlerOptions{})
	group := &Server{
		log: slog.New(loghandler),
		ctx: context.Background(),
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
	return func(s *Server) error {
		s.log = log
		s.ctx = logevent.WithLogger(s.ctx, s.log)
		return nil
	}
}

func (server *Server) initServers() error {

	pipeline := httphandler.MiddlewarePipeline{
		&middleware.Observer{
			BeforeRequest: func(ctx context.Context, req *http.Request, observed *httphandler.Observation) {
				inflightRequests.WithLabelValues(req.Method).Inc()
			},
			AfterRequest: func(ctx context.Context, req *http.Request, observed *httphandler.Observation) {
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
				if len(observed.Attr) > 0 {
					var other []any
					for _, attr := range observed.Attr {
						other = append(other, attr.Key)
					}
					args = append(args, slog.Group("other", other...))
				}
				server.log.Info("handled", args...)
				inflightRequests.WithLabelValues(req.Method).Dec()
				requestDuration.WithLabelValues(req.Method).Observe(duration)
				bytesWritten.WithLabelValues(req.Method).Add(float64(observed.Response.Body.Length))
			},
		},
		&middleware.Recovery{},
		&middleware.HeadMethodHelper{},
	}

	server.mux = httphandler.NewServeMux()

	swp := server.mux.WithPipeline(pipeline)

	swp.Handle("GET /metrics", promhttp.Handler())
	swp.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for _, repoConfig := range server.config.Repositories {
		err := repository.ConfigureRepo(repoConfig, swp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", server.config.Listener.Port)

	if server.config.Listener.CertFile == "" {
		server.log.Info("listening", slog.String("addr", addr))
		err := http.ListenAndServe(addr, server.mux)
		return err
	}

	//get expiration date of cert
	cert, err := tls.LoadX509KeyPair(server.config.Listener.CertFile, server.config.Listener.KeyFile)
	if err != nil {
		return err
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}

	server.log.Info("listening", slog.String("addr", addr), slog.Time("cert_expires", leaf.NotAfter), slog.String("san", strings.Join(leaf.DNSNames, ",")))
	err = http.ListenAndServeTLS(
		addr,
		server.config.Listener.CertFile,
		server.config.Listener.KeyFile,
		server.mux,
	)
	return err
}
