package main

import (
    "time"
    "net/http"
)

type Task struct {
	Description string
	id int
	Due time.Time
	Complete bool
}

type Tasks []Task

func main() {
	tasks := make(Tasks, 10)

    server := &http.Server{
    	Addr : ":8080",
    	Handler : tasks,
    }

    server.ListenAndServe()
}

func (Tasks) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}