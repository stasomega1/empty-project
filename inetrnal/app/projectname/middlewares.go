package projectname

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	ctxClaimsKey ctxKey = iota
	ctxRequestIdKey
	containerName = "projectname"
)

type ctxKey int8

func (s *Server) loggingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		s.Logger.Infof("started %s %s remote_address=%s request_id=%s", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(ctxRequestIdKey))
		now := time.Now()
		rw := &responseWriter{w, http.StatusOK, nil}
		next.ServeHTTP(rw, r)
		s.Logger.Infof("completed with %d %s in %v remote_address=%s request_id=%s",
			rw.statusCode,
			http.StatusText(rw.statusCode),
			time.Now().Sub(now),
			r.RemoteAddr,
			r.Context().Value(ctxRequestIdKey))
	})
}

func (s *Server) setRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxRequestIdKey, id)))
	})
}
