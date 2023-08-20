package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"newsgateway/internal/api"
	"newsgateway/internal/logger"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var log *logger.Logger

// Последние новости
func apiNewsSearch(rw http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	searchQuery := r.URL.Query().Get("query")

	var limit string = fmt.Sprintf("%d", NEWS_PER_PAGE)
	var offset string = fmt.Sprintf("%d", NEWS_PER_PAGE*(page-1))

	apiNews, err := api.New("news")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Problem with news service!")
		log.Error(r.Header.Get("X-Request-Id"), "Problem with news service!")
		return
	}
	var news []api.Post
	err = apiNews.Get(&news, "/api/news/search", map[string]string{"limit": limit, "offset": offset, "query": searchQuery}, r.Header.Get("X-Request-Id"))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Problem with news service: %s\n", err.Error())
		log.Error(r.Header.Get("X-Request-Id"), fmt.Sprintf("Problem with news service: %s\n", err.Error()))
		return
	}

	json_line, err := json.Marshal(news)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

func fetchNews(post *api.Post, id int, wg *sync.WaitGroup, e *error, reqId string) {
	defer wg.Done()
	apiNews, err := api.New("news")
	if err != nil {
		e = &err
		return
	}
	err = apiNews.Get(post, fmt.Sprintf("/api/news/byid/%d", id), map[string]string{}, reqId)
	if err != nil {
		e = &err
		return
	}
}

func fetchComments(comments *[]api.Comment, id int, wg *sync.WaitGroup, e *error, reqId string) {
	defer wg.Done()
	apiComments, err := api.New("comments")
	if err != nil {
		e = &err
		return
	}
	err = apiComments.Get(comments, fmt.Sprintf("/api/comment/getforpost/%d", id), map[string]string{}, reqId)

	if err != nil {
		e = &err
		return
	}
}

// Информация о новости по ID
func apiNewsDetail(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	var wg sync.WaitGroup

	var post api.Post
	var errorPost error
	wg.Add(1)
	go fetchNews(&post, id, &wg, &errorPost, r.Header.Get("X-Request-Id"))

	var comments []api.Comment
	var errorComment error
	wg.Add(1)
	go fetchComments(&comments, id, &wg, &errorComment, r.Header.Get("X-Request-Id"))

	wg.Wait()

	var errorTotal string
	if errorPost != nil {
		errorTotal += errorPost.Error() + "\n"
	}
	if errorComment != nil {
		errorTotal += errorComment.Error() + "\n"
	}

	if errorTotal != "" {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, errorTotal)
		log.Error(r.Header.Get("X-Request-Id"), errorTotal)
		return
	}

	type Answer struct {
		Post     api.Post
		Comments []api.Comment
	}
	var answer Answer = Answer{Post: post, Comments: comments}

	json_line, err := json.Marshal(answer)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

func apiCommentsAdd(rw http.ResponseWriter, r *http.Request) {
	type NewCommData struct {
		Comment string `json:"comment"`
		Id      string `json:"id"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}
	var newcomm NewCommData
	err = json.Unmarshal(reqBody, &newcomm)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}
	if newcomm.Comment == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Comment can't be empty!")
		log.Error(r.Header.Get("X-Request-Id"), "Comment can't be empty!")
		return
	}
	apiComments, err := api.New("comments")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Problem with comments service!")
		log.Error(r.Header.Get("X-Request-Id"), "Problem with comments service!")
		return
	}

	var c api.Comment
	err = apiComments.Post(&c, "/api/comment/add", map[string]string{"comment": newcomm.Comment, "id_post": newcomm.Id}, r.Header.Get("X-Request-Id"))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}
	json_line, err := json.Marshal(c)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
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
	// X-Request-Id
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") == "" {
			var reqId string = uuid.NewString()
			r.Header.Add("X-Request-Id", reqId)
			w.Header().Add("X-Request-Id", reqId)
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	log = logger.New(LOG_PREFIX)
	r := mux.NewRouter()
	r.Use(reqIdHandler)
	r.Use(logHandler)
	r.HandleFunc("/api/news/search", apiNewsSearch).Methods("GET")
	r.HandleFunc("/api/news/detail/{id:[0-9]+}", apiNewsDetail).Methods("GET")
	r.HandleFunc("/api/comments/add", apiCommentsAdd).Methods("POST")
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
