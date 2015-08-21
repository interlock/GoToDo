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
	Id int
	Name string
	Desc string
	Completed bool
}

type TaskList []Task

func NewTaskList() *TaskList {
	tasklist := make(TaskList, 0, 15) //make tasklist of arbitrary size
	return &tasklist
}

var ErrIDNotFound = errors.New("The id you are looking for doesn't exist in this task list")

func (tl *TaskList) FindById(id int) (Task, error) {
	for _, task := range *tl {
		if task.Id == id {
			return task, nil
		}
	}

	return Task{}, ErrIDNotFound
}

func (tl *TaskList) RemoveById(id int) error {
	for i, task := range *tl {
		if task.Id == id {
			*tl = append((*tl)[:i], (*tl)[i + 1:]...)
			return nil
		}
	}

	return ErrIDNotFound
}

func (tl *TaskList) AddTask(task Task) {
	task.Id = nextID //make sure ids don't conflict
	nextID++

	*tl = append(*tl, task)
}

/**
	This function provides the framework for the http servers REST API.
	It implements the Handler interface in go's http package
*/
func (tasklist *TaskList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
		case path == "/":
			//only needs GET
			if r.Method == "GET" {
				for i, task := range *tasklist {
					w.Write([]byte(fmt.Sprintf("%d. %s\n", i + 1, task.Desc)))
				}
			}
		case path == "/tasks":
			if r.Method == "GET" {
				jsonTasks, err := json.Marshal(tasklist)
				if err == nil {
					w.Write(jsonTasks)
				} else {
					//if it fails return an error
					w.Write([]byte(err.Error()))
				}
			} else if r.Method == "POST" {
				task := Task{}
				dec := json.NewDecoder(r.Body)
				err := dec.Decode(&task)

				if err == nil {
					task.Id = nextID
					nextID++
					*tasklist = append(*tasklist, task)
					encoder := json.NewEncoder(w)
					encoder.Encode(&task)
				}
			}
		//is the request about a specific task
		case taskIdExp.MatchString(path):
			var id int
			fmt.Sscanf(path, "/tasks/%d", &id)

			if r.Method == "GET" {
			task, err := tasklist.FindById(id)
				if err == nil {
					jsonTask, _ := json.MarshalIndent(&task, "", "    ")
					w.Write([]byte(jsonTask))
				} else {
					w.Write([]byte(err.Error()))
				}
			} else if r.Method == "PUT" {

			} else if r.Method == "DELETE" {

			}
	}
	// switch r.Method {
	// //Retrieve
	// case "GET":
	// 	//view the simple front end
	//
	// 	} else {
	// 		//check if a specific id is being retrieved
	// 		var id int

	// 		if taskIdExp.MatchString(r.URL.Path) {
	// 			fmt.Sscanf(r.URL.Path, "/tasks/%d", &id)
	// 			task, err := tasklist.FindById(id)
	// 			if err == nil {
	// 				jsonTask, _ := json.MarshalIndent(task, "", "    ")
	// 				w.Write([]byte(jsonTask))
	// 			} else {
	// 				w.Write([]byte(err.Error()))
	// 			}
	// 		}
	// 	}
	// //Edit
	// case "PUT":

	// //Remove
	// case "DELETE":

	// //Create
	//
	// }
}