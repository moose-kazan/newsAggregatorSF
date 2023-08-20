package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"servicenews/internal/dbaccess"
	"servicenews/internal/env"
	"servicenews/internal/logger"
	"servicenews/internal/rssfetch"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var db *dbaccess.Store
var wg sync.WaitGroup
var log *logger.Logger

func webApiNewsById(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}
	post, err := db.PostGetById(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}
	json_line, err := json.Marshal(post)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)

}

func webApiNewsSearch(rw http.ResponseWriter, r *http.Request) {
	paramsRaw := r.URL.Query()
	var params map[string]int = map[string]int{
		"limit":  20,
		"offset": 0,
	}
	for k := range params {
		intVal, err := strconv.Atoi(paramsRaw.Get(k))
		if err == nil {
			params[k] = intVal
		}

	}
	searchQuery := r.URL.Query().Get("query")

	posts, err := db.PostSearch(params["limit"], params["offset"], searchQuery)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	json_line, err := json.Marshal(posts)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		log.Error(r.Header.Get("X-Request-Id"), err.Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(json_line)
}

func fetchSource(s dbaccess.Source, chanItems chan dbaccess.Post, chanDone chan int, chanErrors chan string) {
	fmt.Printf("Start thread for source %d\n", s.Id)
	wg.Add(1)
	defer wg.Done()
	for {
		fmt.Printf("Fetch source %d\n", s.Id)
		items, err := rssfetch.Fetch(s.Link)
		if err == nil {
			for _, item := range items {
				var post dbaccess.Post
				post.Source = s.Id
				post.Title = item.Title
				post.Link = item.Link
				post.Content = item.Content
				post.Guid = item.Guid
				post.PubTime = item.PubTime.Unix()
				chanItems <- post
				//fmt.Printf("Put item from source %d\n", s.Id)
			}
		} else {
			chanErrors <- err.Error()
		}
		time.Sleep(time.Second * time.Duration(s.DefaultInterval))
		select {
		case <-chanDone:
			{
				return
			}
		default:
			{
			}
		}
	}
}

func processNewItems(chanItems chan dbaccess.Post, chanDone chan int, chanErrors chan string) {
	fmt.Println("Start proccesing new items")

	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case item := <-chanItems:
			{
				//fmt.Printf("Process item from source %d\n", item.Source)
				_, err := db.PostAdd(item)
				if err != nil {
					chanErrors <- err.Error()
				}
			}
		case <-chanDone:
			{
				return
			}
		}
	}
}

func processErrors(chanDone chan int, chanErrors chan string) {
	fmt.Println("Start proccesing errors")

	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case err := <-chanErrors:
			{
				os.Stderr.WriteString(err)
			}
		case <-chanDone:
			{
				return
			}
		default:
			{

			}
		}
	}
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Header.Get("X-Request-Id"), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
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

	sources, err := db.SourceGetActive()
	if err != nil {
		panic(err)
	}

	var chanItems = make(chan dbaccess.Post, 100)
	var chanDone = make(chan int)
	var chanErrors = make(chan string)
	defer close(chanDone)

	go processNewItems(chanItems, chanDone, chanErrors)
	go processErrors(chanDone, chanErrors)
	for _, source := range sources {
		go fetchSource(source, chanItems, chanDone, chanErrors)
	}

	r := mux.NewRouter()

	r.Use(logHandler)

	r.HandleFunc("/api/news/byid/{id:[0-9]+}", webApiNewsById).Methods("GET")
	r.HandleFunc("/api/news/search", webApiNewsSearch).Methods("GET")
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

	wg.Wait()
}
