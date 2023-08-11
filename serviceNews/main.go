package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"servicenews/internal/dbaccess"
	"servicenews/internal/rssfetch"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var db *dbaccess.Store
var wg sync.WaitGroup

func webApiNews(rw http.ResponseWriter, r *http.Request) {
	count := 100
	posts, err := db.PostGetLast(count)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
		return
	}

	json_line, err := json.Marshal(posts)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, err.Error())
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

func main() {
	var err error
	db, err = dbaccess.New(DSN)
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

	r.HandleFunc("/api/news/latest", webApiNews)
	srv := &http.Server{
		Handler:      r,
		Addr:         LISTEN_SOCKET,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

	wg.Wait()
}
