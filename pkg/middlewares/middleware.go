package middlewares

import (
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL) //send this log to kibana
		next.ServeHTTP(w, r)
	})
}
