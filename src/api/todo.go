package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ToDo struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Content  string `json:"content"`
	Complete bool   `json:"complete"`
}

var ToDos []ToDo

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got DELETE")
	vars := mux.Vars(r)
	id := vars["Id"]

	for index, todo := range ToDos {
		if todo.Id == id {
			ToDos = append(ToDos[:index], ToDos[index+1:]...)
		}
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got PUT")
	vars := mux.Vars(r)
	id := vars["Id"]

	for index, todo := range ToDos {
		if todo.Id == id {
			requestBody, _ := ioutil.ReadAll(r.Body)
			var tempTodo ToDo
			json.Unmarshal(requestBody, &tempTodo)
			ToDos[index] = tempTodo
		}
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got POST")
	requestBody, _ := ioutil.ReadAll(r.Body)
	var todo ToDo
	json.Unmarshal(requestBody, &todo)
	ToDos = append(ToDos, todo)
	json.NewEncoder(w).Encode(todo)
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint: return all")
	json.NewEncoder(w).Encode(ToDos)
}

func returnOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]
	for _, todo := range ToDos {
		if todo.Id == key {
			fmt.Println("returned", todo, " for", r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
			json.NewEncoder(w).Encode(todo)
		}
	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<b>HomePage")
	fmt.Println("endpoint: homepage")
}

func DefaultTodoCreator() {
	ToDos = []ToDo{
		{Id: "1", Name: "test1", Content: "i need to test this", Complete: false},
		{Id: "2", Name: "test2", Content: "pog test 2", Complete: true},
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)                               // home page
	router.HandleFunc("/todos", returnAll)                         // return all todos
	router.HandleFunc("/todo", createTodo).Methods("POST")         // return all todos
	router.HandleFunc("/todos/{Id}", deleteTodo).Methods("DELETE") // delete specific todo
	router.HandleFunc("/todos/{Id}", updateTodo).Methods("PUT")    // delete specific todo
	router.HandleFunc("/todos/{Id}", returnOne)                    // return specific todo based on Id
	log.Fatal(http.ListenAndServe(":10000", router))               // start server
}

func main() {
	DefaultTodoCreator()
	handleRequests()
}

// curl -i -X POST -H "Content-Type: application/json" -d '{"Id" : "3", "Name" : "added todo", "Content" : "this was added with a POST request", "Complete" : false}' http://127.0.0.1:10000/todo
// curl -i -X DELETE http://127.0.0.1:10000/todos/2
// curl -i -X UPDATE -H "Content-Type: application/json" -d '{"Id" : "1", "Name" : "updated todo", "Content" : "this was updated with a PUT", "Complete" : false}' http://127.0.0.1:10000/todos/1
