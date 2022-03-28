package projectname

import (
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type HealthServer struct {
	router     *mux.Router
	liveFunc   http.HandlerFunc
	readyFunc  http.HandlerFunc
	metricFunc http.HandlerFunc
}

func NewHealthServer(health healthcheck.Handler) *HealthServer {
	srv := mux.NewRouter()
	srv.HandleFunc("/ready", health.ReadyEndpoint)
	srv.HandleFunc("/live", health.LiveEndpoint)
	srv.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	return &HealthServer{
		router:     srv,
		liveFunc:   health.LiveEndpoint,
		readyFunc:  health.ReadyEndpoint,
		metricFunc: promhttp.Handler().ServeHTTP,
	}
}

func (h *HealthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
