package middleware

import (
	"errors"
	"net/http"
	"strings"
	"fmt"
	"reflect"

	"fridge/src/auth"
	"fridge/src/response"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		keys := reflect.ValueOf(r.Header).MapKeys()
		strkeys := make([]string, len(keys))
		for i := 0; i < len(keys); i++ {
			strkeys[i] = keys[i].String()
		}
		headres := strings.Join(strkeys, ",") + ", X-Naver-Client-Id,X-Naver-Client-Secret,X-TARGET-URL"
		fmt.Printf("haha")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", headres)
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.VaildToken(r)

		if nil != err {
			response.MakeJsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		next(w, r)
	}
}
