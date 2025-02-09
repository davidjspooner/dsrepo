package container

import (
	"bytes"
	"net/http"
)

type bufferedResponseWriter struct {
	headers http.Header
	status  int
	body    bytes.Buffer
}

var _ http.ResponseWriter = &bufferedResponseWriter{}

func (brw *bufferedResponseWriter) Header() http.Header {
	if brw.headers == nil {
		brw.headers = make(http.Header)
	}
	return brw.headers
}

func (brw *bufferedResponseWriter) Write(chunk []byte) (int, error) {
	if brw.status == 0 {
		brw.status = http.StatusOK
	}
	return brw.body.Write(chunk)
}

func (brw *bufferedResponseWriter) WriteHeader(status int) {
	brw.status = status
}
