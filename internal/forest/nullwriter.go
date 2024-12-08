package forest

import "net/http"

type nullWriter struct {
	responseWriter http.ResponseWriter
}

func (nw *nullWriter) Write(p []byte) (_ int, _ error) {
	return len(p), nil
}

func (nw *nullWriter) Header() http.Header {
	return nw.responseWriter.Header()
}

func (nw *nullWriter) WriteHeader(statusCode int) {
	nw.responseWriter.WriteHeader(statusCode)
}
