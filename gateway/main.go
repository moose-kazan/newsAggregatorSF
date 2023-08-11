package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Последние новости
func apiNewsLatest(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Not implemented")
}

// Информация о новости по ID
func apiNewsDetail(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Not implemented")
}

// Комментарии к новости по её ID
func apiCommentsLast(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Not implemented")
}

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./webroot/")))
	r.HandleFunc("/api/news/latest", apiNewsLatest)
	r.HandleFunc("/api/news/detail/{id:[0-9]+}", apiNewsDetail)
	r.HandleFunc("/api/comments/last/{id:[0-9]+}", apiCommentsLast)
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
