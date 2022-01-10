package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ToDo struct {
	Id       int    `json:"Id"`
	Name     string `json:"Name"`
	Content  string `json:"Content"`
	Complete bool   `json:"Complete"`
}

type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success int `json:"success"`
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got DELETE")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["Id"])

	db, e := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	if e != nil {
		log.Fatal("unable to open db ", e)
	}
	defer db.Close()

	stmt, e := db.Prepare("delete from todos where id=?")
	if e != nil {
		log.Fatal("unable to make statement", e)
	}
	_, e = stmt.Exec(id)
	if e != nil {
		log.Fatal("unable to remove from db", e)
	}

	var success Success
	success.Success = 1
	json.NewEncoder(w).Encode(success)

}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got PUT")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["Id"])

	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Fatal("no valid json available", e)
	}
	var todo ToDo
	json.Unmarshal(requestBody, &todo)

	db, e := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	if e != nil {
		log.Fatal("unable to open db ", e)
	}
	defer db.Close()

	stmt, e := db.Prepare("update todos set Name=?, Content=?, Complete=? where id=?")
	if e != nil {
		log.Fatal("unable to make statement", e)
	}
	_, e = stmt.Exec(todo.Name, todo.Content, todo.Complete, id)
	if e != nil {
		log.Fatal("unable to update", e)
	}

	var success Success
	success.Success = 1
	json.NewEncoder(w).Encode(success)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got POST")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Fatal("no valid json available", e)
	}
	var todo ToDo
	json.Unmarshal(requestBody, &todo)

	db, e := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	if e != nil {
		log.Fatal("unable to open db ", e)
	}
	defer db.Close()

	stmt, e := db.Prepare("insert into todos(Name, Content, Complete) values (?, ?, ?)")
	if e != nil {
		log.Fatal("unable to make statement", e)
	}
	results, e := stmt.Exec(todo.Name, todo.Content, todo.Complete)
	if e != nil {
		log.Fatal("unable to insert to db", e)
	}
	fmt.Println(results)

	var success Success
	success.Success = 1
	json.NewEncoder(w).Encode(success)
}

func returnAll(w http.ResponseWriter, r *http.Request) {

	var ToDos []ToDo
	db, e := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	if e != nil {
		log.Fatal("unable to open db ", e)
	}

	defer db.Close()

	results, e := db.Query("select * from todos")
	if e != nil {
		log.Fatal("error fetching from db", e)
	}

	for results.Next() {
		var tempTodo ToDo
		e = results.Scan(&tempTodo.Id, &tempTodo.Name, &tempTodo.Content, &tempTodo.Complete)

		if e != nil {
			log.Fatal("unable to parse todo", e)
		}

		ToDos = append(ToDos, tempTodo)
	}

	json.NewEncoder(w).Encode(ToDos)
}

func returnOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["Id"])

	var tempTodo ToDo
	db, e := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	if e != nil {
		log.Fatal("unable to open db ", e)
	}

	defer db.Close()

	db.QueryRow("select * from todos where id=?", key).Scan(&tempTodo.Id, &tempTodo.Name, &tempTodo.Content, &tempTodo.Complete)
	if tempTodo.Id == 0 {
		var err Error
		err.Error = "cannot find todo with id of " + strconv.Itoa(key)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(tempTodo)
	}

}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<b>HomePage")
	fmt.Println("endpoint: homepage")
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
	handleRequests()
}

// curl -i -X POST -H "Content-Type: application/json" -d '{"Id" : "3", "Name" : "added todo", "Content" : "this was added with a POST request", "Complete" : false}' http://127.0.0.1:10000/todo
// curl -i -X DELETE http://127.0.0.1:10000/todos/2
// curl -i -X PUT -H "Content-Type: application/json" -d '{"Id" : "1", "Name" : "updated todo", "Content" : "this was updated with a PUT", "Complete" : false}' http://127.0.0.1:10000/todos/1
