package main

import (
	"fmt"
    "net/http"
    "encoding/json"
    "regexp"
)

var taskIdExp = regexp.MustCompile("^/tasks/[0-9]*$")

//A task to be completed on a ToDo List
type Task struct {
	id int
	Name string
	Desc string
	Completed bool
}

type TaskList []Task

func main() {
	var tasklist TaskList = make([]Task, 0, 10) //start with a capacity of 10

    server := &http.Server{
    	Addr : ":80",
    	Handler : tasklist,
    }

    server.ListenAndServe()
}

/**
	This function provides the framework for the http servers REST API.
	It implements the Handler interface in go's http package
*/
func (tasklist TaskList) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	//Retrieve
	case "GET":
		//view the simple front end
		if req.URL.Path == "/" {
			for i, task := range tasklist {
				rw.Write([]byte(fmt.Sprintf("%d. %s\n", i, task.Desc)))
			}
		} else if req.URL.Path == "/tasks" {
			//convert the data tasklist into byte slice encoded in json
			jsonTasks, err := json.Marshal(tasklist)
			if err == nil {
				rw.Write(jsonTasks)
			} else {
				//if it fails return an error
				rw.Write([]byte("404 server error"))
			}
		} else {
			//check if a specific id is being retrieved
			var id int

			if taskIdExp.MatchString(req.URL.Path) {
				fmt.Sscanf(req.URL.Path, "/tasks/%d", &id)
				rw.Write([]byte(fmt.Sprintf("%v", id)))
			}
		}
	//Edit
	case "PUT":

	//Remove
	case "DELETE":

	//Create
	case "POST":
		if req.URL.Path == "/tasks" {

		}
	}
}