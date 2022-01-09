package main

import (
	"encoding/json"
	"fmt"
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

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint: return all")
	json.NewEncoder(w).Encode(ToDos)
}

func returnOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]
	//fmt.Fprintf(w, "Key: "+key)
	for _, todo := range ToDos {
		if todo.Id == key {
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
	router.HandleFunc("/", homePage)
	router.HandleFunc("/todos/{Id}", returnOne)
	router.HandleFunc("/todos", returnAll)
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	DefaultTodoCreator()
	handleRequests()
}
