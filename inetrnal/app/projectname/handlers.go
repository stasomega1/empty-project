package projectname

import (
	"net/http"
	"project/inetrnal/app/model"
	"strconv"
	"strings"
)

func getGETStringValue(r *http.Request, valueKey string) string {
	query := r.URL.Query()
	value := strings.TrimSpace(query.Get(valueKey))
	return value
}

func getGETIntValue(r *http.Request, valueKey string) int {
	query := r.URL.Query()
	value := strings.TrimSpace(query.Get(valueKey))
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return valueInt
}

func (s *Server) handleGetSomeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameter1 := getGETStringValue(r, "parameter1")
		parameter2 := getGETIntValue(r, "parameter2")

		request := model.DbModelRequest{
			Parameter1: parameter1,
			Parameter2: parameter2,
		}

		//SERVICE CALL//
		result, err := s.ProjectService.GetSomeData(request)
		s.defaultResponse(w, result, err)
	}
}
