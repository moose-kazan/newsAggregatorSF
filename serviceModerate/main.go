package main

import (
	"fmt"
	"net/http"
	"servicemoderate/internal/logger"
	"time"

	"github.com/gorilla/mux"
)

var log *logger.Logger

func webApiModerateComment(rw http.ResponseWriter, r *http.Request) {

}

func reqIdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") != "" {
			w.Header().Add("X-Request-Id", r.Header.Get("X-Request-Id"))
		}
		next.ServeHTTP(w, r)
	})
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Header.Get("X-Request-Id"), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func main() {
	log = logger.New(LOG_PREFIX)

	r := mux.NewRouter()

	r.Use(reqIdHandler)
	r.Use(logHandler)

	r.HandleFunc("/api/moderate/comment", webApiModerateComment).Methods("POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

}
