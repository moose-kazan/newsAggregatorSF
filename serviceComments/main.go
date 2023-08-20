package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servicecomments/internal/dbaccess"
	"servicecomments/internal/env"
	"servicecomments/internal/logger"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var db *dbaccess.Store
var log *logger.Logger

func Send503(rw http.ResponseWriter, r *http.Request, msg string) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(rw, msg)
	log.Error(r.Header.Get("X-Request-Id"), msg)
}

func webApiCommentAdd(rw http.ResponseWriter, r *http.Request) {
	newcomm := r.FormValue("comment")
	id_post, err := strconv.Atoi(r.FormValue("id_post"))
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}
	if id_post < 1 {
		Send503(rw, r, "IdPost must be > 0")
		return
	}
	if newcomm == "" {
		Send503(rw, r, "Comment can't be empty!")
		return
	}
	c, err := db.Add(id_post, newcomm)
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}

	json_line, err := json.Marshal(c)
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

func webApiCommentGetForPost(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}
	comments, err := db.GetForPost(id)
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}
	json_line, err := json.Marshal(comments)
	if err != nil {
		Send503(rw, r, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)

}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Header.Get("X-Request-Id"), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func reqIdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") != "" {
			w.Header().Add("X-Request-Id", r.Header.Get("X-Request-Id"))
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	log = logger.New(LOG_PREFIX)

	var err error

	var db_port = env.GetInt("DB_PORT", DEFAULT_DB_PORT)
	var db_host = env.GetStr("DB_HOST", DEFAULT_DB_HOST)

	db, err = dbaccess.New(fmt.Sprintf(DSN, db_host, db_port))
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Use(reqIdHandler)
	r.Use(logHandler)

	r.HandleFunc("/api/comment/getforpost/{id:[0-9]+}", webApiCommentGetForPost).Methods("GET")
	r.HandleFunc("/api/comment/add", webApiCommentAdd).Methods("POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
