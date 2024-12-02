package repository

import (
	"log/slog"
	"net/http"
	"time"
)

type Response struct {
	statusCode   int
	started      time.Time
	bytesWritten int
	inner        http.ResponseWriter
	Log          *slog.Logger
}

func NewResponse(inner http.ResponseWriter, log *slog.Logger) *Response {
	return &Response{
		statusCode: 200,
		started:    time.Now(),
		inner:      inner,
		Log:        log,
	}
}

func (r *Response) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.inner.WriteHeader(statusCode)
}

func (r *Response) Write(b []byte) (int, error) {
	n, err := r.inner.Write(b)
	r.bytesWritten += n
	return n, err
}

func (r *Response) Header() http.Header {
	return r.inner.Header()
}

func (r *Response) Duration() time.Duration {
	return time.Since(r.started)
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) BytesWritten() int {
	return r.bytesWritten
}
