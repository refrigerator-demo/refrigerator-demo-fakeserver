package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeJson(w http.ResponseWriter, statusCode int, data interface{}) {

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if nil != err {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func MakeJsonError(w http.ResponseWriter, statusCode int, err error) {
	if nil != err {
		MakeJson(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	MakeJson(w, http.StatusBadRequest, nil)
}
