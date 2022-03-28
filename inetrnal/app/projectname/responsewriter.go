package projectname

import (
	"bytes"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	respBuffer *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.respBuffer == nil {
		rw.respBuffer = new(bytes.Buffer)
	}
	rw.respBuffer.Write(b)
	return rw.ResponseWriter.Write(b)
}
