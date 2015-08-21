package main

import (
    "net/http"
    "GoToDo/tasks"
)

func main() {
	tasklist := make(tasks.TaskList, 0, 10) //start with a capacity of 10

    server := &http.Server{
    	Addr : ":80",
    	Handler : &tasklist,
    }

    server.ListenAndServe()
}