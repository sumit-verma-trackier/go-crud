package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Task   string  `json:"task"`
	Status *Status `json:"status"`
}

type Status struct {
	Completed bool `json:"completed"`
	Pending   bool `json:"pending"`
}

// using memory for learing
var todos []Todo

func main() {

	addTodo := Todo{
		ID:   1,
		Name: "Jophn",
		Task: "Nlk",
		Status: &Status{
			Completed: true,
			Pending:   false,
		},
	}

	todos = append(todos, addTodo)

	ro := mux.NewRouter()

	ro.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	// ro.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

	// 	msg := make(map[string]string)

	// 	msg["message"] = "pong"
	// 	// msg["two"] = "Second"

	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(msg)
	// })

	// ro.HandleFunc("/books/{title}/{page}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	title := vars["title"]
	// 	page := vars["page"]

	// 	fmt.Fprintf(w, "You've requested the book: %s & page: %s", title, page)
	// })

	// Gettin todo list
	ro.HandleFunc("/todo/list", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if len(todos) == 0 {
			json.NewEncoder(w).Encode(map[string]string{"message": "no data found"})
			return
		}

		err := json.NewEncoder(w).Encode(map[string]any{"status": true, "response": todos})
		if err != nil {
			http.Error(w, "error encoding todos", http.StatusInternalServerError)
			return
		}
	})

	ro.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)

		idStr := params["id"]

		// Convert id string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		for _, todo := range todos {
			if todo.ID == id {
				json.NewEncoder(w).Encode(map[string]any{
					"status":   true,
					"response": todo,
				})
				return
			}
		}

		// If not found
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "todo not found",
		})

	})

	ro.HandleFunc("/todo/update/{id}", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json") // seting the header

		params := mux.Vars(r)
		idStr := params["id"] // gettng the param

		// dummyName := r.URL.Query().Get("dummyName") // reading value from Query Params
		// fmt.Println(dummyName, "----")

		id, err := strconv.Atoi(idStr) // converting into int from string
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updatedTodo Todo

		eerr := json.NewDecoder(r.Body).Decode(&updatedTodo)
		if eerr != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		for index, todo := range todos {
			if todo.ID == id {

				updatedTodo.ID = id
				// todos[index].Name = updatedTodo.Name
				// todos[index].Task = updatedTodo.Task
				// todos[index].Status.Completed = updatedTodo.Status.Completed
				// todos[index].Status.Pending = updatedTodo.Status.Pending
				todos[index] = updatedTodo

				json.NewEncoder(w).Encode(map[string]any{
					"status":   true,
					"message":  "todo updated successfully",
					"response": updatedTodo,
				})
				return
			}
		}

		// If not found
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "todo not found",
		})
	})

	ro.HandleFunc("/todo/add/new", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "invalid method, use POSt instead", http.StatusMethodNotAllowed)
			return
		}

		var payload Todo

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid or empty JSON body", http.StatusBadRequest)
			return
		}

		todos = append(todos, payload)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  true,
			"message": "Todo added successfully",
			"data":    payload,
		})

	})

	http.ListenAndServe(":8001", ro)
}
