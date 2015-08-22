package tasks

import (
	"testing"
	"net/http"
	"encoding/json"
	"bytes"
	"os"
)

var client = &http.Client{}
var tasklist = NewTaskList()

var tasks = []Task{
		{Name:"Laundry", Desc:"Clothes are dirty need to clean"},
		{Name:"Math Homework", Desc:"Need to finish calculus unit 5"},
		{Name:"Look for job", Desc:"Unemployed need work"},
	}

func TestPut(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

func TestPost(t *testing.T) {
	for i, task := range tasks {
		jsontask, _ := json.Marshal(task)
		req, _ := http.NewRequest("POST", "http://localhost:8080/tasks", bytes.NewBuffer(jsontask))
		r, err := client.Do(req)

		if err == nil {
			var retTask Task

			decoder := json.NewDecoder(r.Body)
			decoder.Decode(&retTask)

			//ids start at 1, index in range will start at 0
			if retTask.Id != i + 1 {
				t.Errorf("Post failed returned id was %d and should have been %d", retTask.Id, i + 1)
			}
		} else {
			t.Error("http request failed")
		}
	}
}

func TestMain(m *testing.M) {
	server := &http.Server{
		Addr : ":8080",
		Handler : tasklist,
	}

	go server.ListenAndServe()

	os.Exit(m.Run())
}