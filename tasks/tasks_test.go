package tasks

import (
	"testing"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"os"
)

var client = &http.Client{}
var tasklist = NewTaskList()

var tasks = []Task{
		{Name:"Laundry", Desc:"Clothes are dirty need to clean"},
		{Name:"Math Homework", Desc:"Need to finish calculus unit 5"},
		{Name:"Look for job", Desc:"Unemployed need work"},
	}

func TestPost(t *testing.T) {
	for i, task := range tasks {
		jsontask, _ := json.Marshal(task)
		buf := bytes.NewBuffer(jsontask)
		req, _ := http.NewRequest("POST", "http://localhost:8080/tasks", buf)
		r, err := client.Do(req)

		if err == nil {
			var retTask Task

			decoder := json.NewDecoder(r.Body)
			decoder.Decode(&retTask)

			//ids start at 1, index in range will start at 0
			if retTask.Id != i + 1 {
				t.Errorf("Post failed returned id was %d and should have been %d", retTask.Id, i + 1)
			}

			if retTask.Name != task.Name {
				t.Error("Name mismatch")
			}

			if retTask.Desc != task.Desc {
				t.Error("Description mismatch")
			}			
		} else {
			t.Error("http request failed")
		}
	}
}

func TestGet(t *testing.T) {
	for i, task := range tasks {
		url := fmt.Sprintf("http://localhost:8080/tasks/%d", i + 1)
		req, _ := http.NewRequest("GET", url, nil)
		r, err := client.Do(req)

		if err == nil {
			var retTask Task

			decoder := json.NewDecoder(r.Body)
			decoder.Decode(&retTask)

			if retTask.Name != task.Name {
				t.Error("Name mismatch")
			}

			if retTask.Desc != task.Desc {
				t.Error("Description mismatch")
			}
		} else {
			t.Error("GET request failed")
		}
	}

	//Now need to test getting /tasks
	var getTasks []Task
	req, _ := http.NewRequest("GET", "http://localhost:8080/tasks", nil)
	r, err := client.Do(req)

	if err == nil {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&getTasks)
		if err == nil {
			for i, task := range getTasks {
				if tasks[i].Name != task.Name {
					t.Error("Name mismatch")
				}

				if tasks[i].Desc != task.Desc {
					t.Error("Description mismatch")
				}
			}
		} 
	}
}

func TestPut(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

func TestMain(m *testing.M) {
	server := &http.Server{
		Addr : ":8080",
		Handler : tasklist,
	}

	go server.ListenAndServe()

	os.Exit(m.Run())
}