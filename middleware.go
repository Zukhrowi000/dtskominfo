package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"    // You need to install this package: go get -u github.com/gorilla/mux
	"github.com/urfave/negroni" // You need to install this package: go get -u github.com/urfave/negroni
)

// Task struct represents a task
type Task struct {
	OrderID      int    `json:"order_id"`
	CustomerName string `json:"customer_name"`
	OrderedAt    string `json:"ordered_at"`
}

var tasks []Task

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Route handlers
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{order_id}", getTask).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{order_id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{order_id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/hello", helloHandler).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))

	// Create a new Negroni instance
	n := negroni.New()

	// Add middleware
	// n.Use(negroni.HandlerFunc(customMiddleware))

	// Attach the router to Negroni
	n.UseHandler(router)

	// Start the server
	http.ListenAndServe(":8080", n)
}

// Get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get single task by ID
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through tasks and find one with the given ID
	for _, item := range tasks {
		if strconv.Itoa(item.OrderID) == params["order_id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{})
}

// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.OrderID = len(tasks) + 1
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}

// Update a task
func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	for index, item := range tasks {
		if strconv.Itoa(item.OrderID) == params["order_id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)
			var task Task
			_ = json.NewDecoder(r.Body).Decode(&task)
			task.OrderID, _ = strconv.Atoi(params["order_id"])
			tasks = append(tasks, task)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	json.NewEncoder(w).Encode(tasks)
}

// Delete a task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	for index, item := range tasks {
		if strconv.Itoa(item.OrderID) == params["order_id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tasks)
}

// Handler function for "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
