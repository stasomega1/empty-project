package projectname

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"project/inetrnal/app/model"
	"project/inetrnal/app/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router         *mux.Router
	Logger         *logrus.Logger
	ErrLevel       string
	ProjectService services.ServiceI
}

const (
	Debug  = "debug"
	Normal = "normal"
)

func newServer(service services.ServiceI, logger *logrus.Logger, errLevel string) *Server {
	srv := &Server{
		Router:         mux.NewRouter(),
		Logger:         logger,
		ProjectService: service,
		ErrLevel:       errLevel,
	}
	srv.configureRouter()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	//Set CORS middleware for all requests
	s.Router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodOptions}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.ExposedHeaders([]string{"Authorization", "Set-Cookie"})))

	//Public middleware
	//Set request id middleware for all server requests
	s.Router.Use(s.setRequestId)
	//Set logging middleware for all server requests
	s.Router.Use(s.loggingRequests)

	//Common endpoints
	s.Router.Handle("/metrics", promhttp.Handler())

	//Debug
	s.Router.HandleFunc("/getSomeData", s.handleGetSomeData()).Methods(http.MethodGet, http.MethodOptions)
}

const (
	simpleErrCode     = "ERROR"
	simpleSuccessCode = "SUCCESS"
)

type SimpleResponse struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	DebugMessage string `json:"debugMessage,omitempty"`
}

func (s *Server) defaultResponse(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		s.error(w, err)
	} else {
		s.respond(w, http.StatusOK, data)
	}
}

func (s *Server) error(w http.ResponseWriter, err error) {
	response := &SimpleResponse{}
	status := 500

	switch myErr := err.(type) {
	case *model.MyError:
		response = &SimpleResponse{
			Code:    simpleErrCode,
			Message: myErr.UserMessage,
		}
		if s.ErrLevel == Debug {
			response.DebugMessage = myErr.DebugInfo
		}
		status = myErr.HttpStatus
	default:
		response = &SimpleResponse{
			Code:    simpleErrCode,
			Message: "Неизвестная ошибка",
		}
		if s.ErrLevel == Debug {
			response.DebugMessage = myErr.Error()
		}
	}

	s.respond(w, status, response)
}

func (s *Server) success(w http.ResponseWriter, status int, message string) {
	s.respond(w, status, &SimpleResponse{Code: simpleSuccessCode, Message: message})
}

func (s *Server) respond(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}

func (s *Server) respondFile(w http.ResponseWriter, status int, data []byte, fileName string) {
	if data != nil {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		w.WriteHeader(status)
		//_, _ = w.Write(data)
		r := bytes.NewReader(data)
		_, _ = io.Copy(w, r)
	} else {
		w.WriteHeader(status)
	}
}
