package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servicemoderate/internal/filter"
	"servicemoderate/internal/logger"
	"time"

	"github.com/gorilla/mux"
)

var log *logger.Logger

func Send503(rw http.ResponseWriter, r *http.Request, msg string) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(rw, msg)
	log.Error(r.Header.Get("X-Request-Id"), msg)
}

func webApiModerateBadWords(rw http.ResponseWriter, r *http.Request) {
	s := r.FormValue("text")
	haveBadWords := filter.BadWords(s)
	var answer filter.Result = filter.Result{Filtered: haveBadWords}
	json_line, err := json.Marshal(answer)
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}
	if haveBadWords {
		rw.WriteHeader(http.StatusBadRequest) // 400
	} else {
		rw.WriteHeader(http.StatusOK) // 200
	}
	rw.Write(json_line)
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

	r.HandleFunc("/api/moderate/badwords", webApiModerateBadWords).Methods("POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

}
