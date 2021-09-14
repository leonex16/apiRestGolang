package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	Id      int    `json:"Id"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

var tasks = allTasks{
	{
		Id:      1,
		Name:    "Name",
		Content: "content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	var reqBody, err = ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Task not valid")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.Id = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var taskId, err = strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for _, task := range tasks {
		if task.Id == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode((task))
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var taskId, err = strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for i, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func indexRoute(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Pruea")
}

func main() {
	var router = mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute).Methods("GET")
	// curl -X GET localhost:1234/tasks
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	// curl -X POST -H "Content-Type: aplication/json" -d '{"Name":"Name2","Content":"Content2"}' localhost:1234/tasks
	router.HandleFunc("/tasks", createTask).Methods("POST")
	// curl -X GET localhost:1234/tasks/1
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	// curl -X DELETE localhost:1234/tasks/1
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":1234", router))
}
