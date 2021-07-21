package controller

import (
	"fridge/src/response"
	"net/http"
)

func (server *RestServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.MakeJson(w, http.StatusOK, "Health Checker")
}
