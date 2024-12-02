package forest

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/davidjspooner/dshttp/pkg/httphandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var inflightRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "http_inflight_requests",
	Help: "Number of inflight requests",
}, []string{"listener", "method"})

var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "HTTP request duration",
	Buckets: prometheus.DefBuckets,
}, []string{"listener", "method"})

var bytesWritten = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_size_bytes",
		Help: "Total size of the response sent, in bytes.",
	},
	[]string{"listener", "method"},
)

type server struct {
	group          *Group
	log            *slog.Logger
	listenerConfig *ListenerConfig
	mux            *httphandler.ServeMux
}

func (s *server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.listenerConfig.Port)

	if s.listenerConfig.CertFile == "" {
		s.log.Info("listening", slog.String("addr", addr))
		err := http.ListenAndServe(addr, s)
		return err
	}

	//get expiration date of cert
	cert, err := tls.LoadX509KeyPair(s.listenerConfig.CertFile, s.listenerConfig.KeyFile)
	if err != nil {
		return err
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}

	s.log.Info("listening", slog.String("addr", addr), slog.Time("cert_expires", leaf.NotAfter), slog.String("san", strings.Join(leaf.DNSNames, ",")))
	err = http.ListenAndServeTLS(
		addr,
		s.listenerConfig.CertFile,
		s.listenerConfig.KeyFile,
		s,
	)
	return err
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		req.Method = http.MethodGet
		w = &nullWriter{responseWriter: w}
	}

	switch req.URL.Path {
	case "/metrics":
		s.group.metricsHandler(w, req)
		return
	case "/health":
		s.group.healthHandler(w, req)
		return
	}
	s.mux.ServeHTTP(w, req)
}
