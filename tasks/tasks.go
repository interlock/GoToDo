package tasks

import (
	"errors"
	"net/http"
	"regexp"
	"fmt"
	"encoding/json"
)

//simple way to keep track of ids
var nextID int = 1

//regexp for seeing if the url path is of the form /tasks/id
var taskIdExp = regexp.MustCompile("^/tasks/[0-9]+$")

//A task to be completed on a ToDo List
type Task struct {
	id int
	name string
	desc string
	completed bool
}

type TaskList []Task

var ErrIDNotFound = errors.New("The id you are looking for doesn't exist in this task list")

//when refactoring into own module consider adding errors
func (tasklist TaskList) FindById(id int) (Task, error) {
	for _, task := range tasklist {
		if task.id == id {
			return task, nil
		}
	}

	return Task{}, ErrIDNotFound
}


/**
	This function provides the framework for the http servers REST API.
	It implements the Handler interface in go's http package
*/
func (tasklist *TaskList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//Retrieve
	case "GET":
		//view the simple front end
		if r.URL.Path == "/" {
			for i, task := range *tasklist {
				w.Write([]byte(fmt.Sprintf("%d. %s\n", i, task.desc)))
			}
		} else if r.URL.Path == "/tasks" {
			//convert the data tasklist into byte slice encoded in json
			jsonTasks, err := json.Marshal(tasklist)
			if err == nil {
				w.Write(jsonTasks)
			} else {
				//if it fails return an error
				w.Write([]byte("404 server error"))
			}
		} else {
			//check if a specific id is being retrieved
			var id int

			if taskIdExp.MatchString(r.URL.Path) {
				fmt.Sscanf(r.URL.Path, "/tasks/%d", &id)
				task, err := tasklist.FindById(id)
				if err == nil {
					jsonTask, _ := json.MarshalIndent(task, "", "    ")
					w.Write([]byte(jsonTask))
				} else {
					w.Write([]byte(err.Error()))
				}
			}
		}
	//Edit
	case "PUT":

	//Remove
	case "DELETE":

	//Create
	case "POST":
		if r.URL.Path == "/tasks" {
			task := Task{}
			dec := json.NewDecoder(r.Body)
			err := dec.Decode(&task)
			fmt.Println(task)
			if err == nil {
				task.id = nextID
				nextID++
				*tasklist = append(*tasklist, task)
				encoder := json.NewEncoder(w)
				encoder.Encode(task)
			}
		}
	}
}