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
				return
			}

			if retTask.Name != task.Name {
				t.Error("Name mismatch")
				return
			}

			if retTask.Desc != task.Desc {
				t.Error("Description mismatch")
				return
			}			
		} else {
			t.Error("http request failed")
			return
		}
	}

	fmt.Println("finished TestPost")
}

//this test also ensures that post was working correctly
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
				return
			}

			if retTask.Desc != task.Desc {
				t.Error("Description mismatch")
				return
			}
		} else {
			t.Error("GET request failed")
			return
		}
	}

	//Now need to test getting /tasks
	var getTasks []Task
	req, _ := http.NewRequest("GET", "http://localhost:8080/tasks", nil)
	r, err := client.Do(req)

	if err == nil {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&getTasks)
		if len(tasks) != len(getTasks) {
			t.Error("length mismatch")
		}
		if err == nil {
			for i, task := range getTasks {
				if tasks[i].Name != task.Name {
					t.Error("Name mismatch")
					return
				}

				if tasks[i].Desc != task.Desc {
					t.Error("Description mismatch")
					return
				}
			}
		} 
	}

	fmt.Println("finished TestGet")
}

func TestPut(t *testing.T) {
//modify one of the tasks
	mod := make(map[string]interface{})
	mod["desc"] = "Now finished"
	mod["completed"] = true

	//encode to json
	jsonmod, _ := json.Marshal(mod)
	buf := bytes.NewBuffer(jsonmod)
	//now send the modifications
	req, _ := http.NewRequest("PUT", "http://localhost:8080/tasks/2", buf)
	r, err := client.Do(req)

	if err != nil {
		t.Error("PUT request failed")
		return
	}
	//decode the returned modified task
	var task Task
	dec := json.NewDecoder(r.Body)
	dec.Decode(&task)

	//make sure the values have been changed and
	//and those that haven't are what they should be
	if task.Id != 2 {
		t.Error("id mismatch in TestPut")
		return
	}
	//name was not changed
	if task.Name != tasks[1].Name {
		t.Error("name mismatch in TestPut")
		return
	}
	
	//description was changed
	if task.Desc != "Now finished" {
		t.Error("desc mismatch in TestPut")
		return
	}

	//completed status changed
	if task.Completed != true {
		t.Error("status mismatch in TestPut")
		return
	}

	fmt.Println("finished TestPut")
}

func TestDelete(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/tasks/1", nil)

	r, err := client.Do(req)

	if err != nil {
		t.Error("DELETE request failed")
		return
	}

	//it returns "{}\n" in the body
	body := make([]byte, 2)
	r.Body.Read(body)
	if string(body) != "{}" {
		t.Errorf("DELETE returned %s", string(body))
		return
	}
	
	//now must GET all tasks and see if it
	// deleted the right one

	var getTasks []Task // to store the tasks after GET
	req, _ = http.NewRequest("GET", "http://localhost:8080/tasks", nil)

	r, err = client.Do(req)

	if err != nil {
		t.Error("GET request failed in TestDelete")
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.Decode(&getTasks)

	//there should only be 2 tasks left
	if len(getTasks) != 2 {
		t.Error("number of tasks mismatched")
	}

	//deleted task id 1 so the left over should be
	//tasks 2 and 3
	for i := range getTasks {
		if getTasks[i].Id != i + 2 {
			t.Error("id mismatch after DELETE")
		}
	}

	fmt.Println("finished TestDelete")
}

func TestMain(m *testing.M) {
	server := &http.Server{
		Addr : ":8080",
		Handler : tasklist,
	}

	go server.ListenAndServe()

	os.Exit(m.Run())
}