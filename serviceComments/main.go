package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servicecomments/internal/dbaccess"
	"servicecomments/internal/env"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var db *dbaccess.Store
var wg sync.WaitGroup

func webApiCommentAdd(rw http.ResponseWriter, r *http.Request) {
	newcomm := r.FormValue("comment")
	id_post, err := strconv.Atoi(r.FormValue("id_post"))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}
	if id_post < 1 {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "IdPost must be > 0")
		return
	}
	if newcomm == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Comment can't be empty!")
		return
	}
	c, err := db.Add(id_post, newcomm)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}
	json_line, err := json.Marshal(c)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

func webApiCommentGetForPost(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}
	comments, err := db.GetForPost(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}
	json_line, err := json.Marshal(comments)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)

}

func main() {
	var err error

	var db_port = env.GetInt("DB_PORT", DEFAULT_DB_PORT)
	var db_host = env.GetStr("DB_HOST", DEFAULT_DB_HOST)

	db, err = dbaccess.New(fmt.Sprintf(DSN, db_host, db_port))
	if err != nil {
		panic(err)
	}

	var chanDone = make(chan int)
	defer close(chanDone)

	r := mux.NewRouter()

	r.HandleFunc("/api/comment/getforpost/{id:[0-9]+}", webApiCommentGetForPost).Methods("GET")
	r.HandleFunc("/api/comment/add", webApiCommentAdd).Methods("POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

	wg.Wait()

}
