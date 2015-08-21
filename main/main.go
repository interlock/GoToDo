package main

import (
    "net/http"
    "GoToDo/tasks"
)

func main() {
	tasklist := tasks.NewTaskList()

    server := &http.Server{
    	Addr : ":80",
    	Handler : tasklist,
    }

    server.ListenAndServe()
}