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

func (tl *TaskList) FindById(id int) (*Task, error) {
	for i := range *tl {
		if (*tl)[i].Id == id {
			return &(*tl)[i], nil
		}
	}

	return nil, ErrIDNotFound
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

func (tl *TaskList) AddTask(task *Task) {
	task.Id = nextID //make sure ids don't conflict
	nextID++

	*tl = append(*tl, *task)
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
					w.Write([]byte(fmt.Sprintf("%d. %s : %s\n", i + 1, task.Name, task.Desc)))
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
				task := &Task{}
				dec := json.NewDecoder(r.Body)
				err := dec.Decode(task)

				if err == nil {
					tasklist.AddTask(task)
					encoder := json.NewEncoder(w)
					encoder.Encode(task)
				}
			}
		//is the request about a specific task
		case taskIdExp.MatchString(path):
			var id int
			fmt.Sscanf(path, "/tasks/%d", &id)

			if r.Method == "GET" {
				task, err := tasklist.FindById(id)
				if err == nil {
					encoder := json.NewEncoder(w)
					encoder.Encode(task)
				} else {
					w.Write([]byte(err.Error()))
				}
			} else if r.Method == "PUT" {
				//used for editing tasks
				task, err := tasklist.FindById(id)
				if err == nil {
					jsonParsed := make(map[string]interface{})
					decoder := json.NewDecoder(r.Body)
					decoder.Decode(&jsonParsed)

					if val, ok := jsonParsed["name"]; ok {
						task.Name, _ = val.(string)
					}

					if val, ok := jsonParsed["desc"]; ok {
						task.Desc, _ = val.(string)
					}

					if val, ok := jsonParsed["completed"]; ok {
						task.Completed, _ = val.(bool)
					}

					encoder := json.NewEncoder(w)
					encoder.Encode(task)
				} else {
					w.Write([]byte(err.Error()))
				}
			} else if r.Method == "DELETE" {
				//remove a task
				err := tasklist.RemoveById(id)
				if err == nil {
						w.Write([]byte("{}\n")) //empty json object for success ?
					} else {
						w.Write([]byte(err.Error())) //return error on failure
					}
			}
	}
}