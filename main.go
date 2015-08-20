package main

import (
    "time"
    "net/http"
    "fmt"
)

//A task to be completed on a ToDo List
type Task struct {
	id int
	Description string
	Completed bool
}

//The data type holding all info for the application
type ToDoApp struct {

}

func main() {
	todo := new(ToDoApp)

    server := &http.Server{
    	Addr : ":80",
    	Handler : todo,
    }

    server.ListenAndServe()
}

/**
	This function provides the framework for the http servers REST API.
	It implements the Handler interface in go's http package
*/
func (todo *ToDoApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		
	case "PUT":

	case "DELETE":

	case "POST":

	}
}