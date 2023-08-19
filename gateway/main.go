package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newsgateway/api"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Последние новости
func apiNewsLatest(rw http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	var limit string = fmt.Sprintf("%d", NEWS_PER_PAGE)
	var offset string = fmt.Sprintf("%d", NEWS_PER_PAGE*(page-1))

	apiNews, err := api.New("news")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Problem with news service!")
		return
	}
	var news []api.Post
	err = apiNews.Get(&news, "/api/news/latest", map[string]string{"limit": limit, "offset": offset})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Problem with news service: %s\n", err.Error())
		return
	}

	json_line, err := json.Marshal(news)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

// Информация о новости по ID
func apiNewsDetail(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	apiNews, err := api.New("news")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Problem with news service!")
		return
	}
	var post api.Post
	err = apiNews.Get(&post, fmt.Sprintf("/api/news/byid/%d", id), map[string]string{})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Problem with news service: %s\n", err.Error())
		return
	}

	json_line, err := json.Marshal(post)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

// Комментарии к новости по её ID
func apiCommentsLast(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Not implemented")
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("WebRequest: %s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(logHandler)
	r.HandleFunc("/api/news/latest", apiNewsLatest)
	r.HandleFunc("/api/news/detail/{id:[0-9]+}", apiNewsDetail)
	r.HandleFunc("/api/comments/last/{id:[0-9]+}", apiCommentsLast)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./webroot/")))
	//r.Handle("/", http.FileServer(http.Dir("./webroot/")))
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
