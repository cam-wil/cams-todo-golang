package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var e error

type ToDo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Complete bool   `json:"complete"`
}

func printError(s string, err error) {
	if e != nil {
		log.Fatal(s, err)
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// get id number from end of url
	vars := mux.Vars(r)
	id, e := strconv.Atoi(vars["Id"])
	printError("no valid id provided", e)

	// prepare a statement
	stmt, e := db.Prepare("DELETE from todos WHERE id=?")
	printError("unable to make statement", e)

	// execute prepared statement
	res, e := stmt.Exec(id)
	printError("unable to execute statement", e)

	// if no rows affected, send 400
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(400)
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	// get id number from end of url
	vars := mux.Vars(r)
	id, e := strconv.Atoi(vars["Id"])
	printError("no valid id provided", e)

	// get json from body
	requestBody, e := ioutil.ReadAll(r.Body)
	printError("no valid json available", e)

	// turn json into ToDo struct
	todo := ToDo{}
	e = json.Unmarshal(requestBody, &todo)
	printError("cannot unmarshal", e)

	// create statement for update
	stmt, e := db.Prepare("UPDATE todos set Name=?, Content=?, Complete=? WHERE id=?")
	printError("unable to make statement", e)

	// execute statement
	res, e := stmt.Exec(todo.Name, todo.Content, todo.Complete, id)
	printError("unable to execute statement", e)

	// if now rows affected, throw 400
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(400)
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	// get json from body
	requestBody, e := ioutil.ReadAll(r.Body)
	printError("no valid json available", e)

	// unmarshal json to ToDo struct
	todo := ToDo{}
	e = json.Unmarshal([]byte(requestBody), &todo)
	printError("error unmarshaling to struct", e)

	// prepare a statement
	stmt, e := db.Prepare("INSERT into todos(Name, Content, Complete) VALUES (?, ?, ?)")
	printError("unable to execute statement", e)

	// execute statement
	res, e := stmt.Exec(todo.Name, todo.Content, todo.Complete)
	printError("unable to insert to db", e)

	// if now rows affected, throw 400
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(400)
	}
}

func returnOne(w http.ResponseWriter, r *http.Request) {
	// get int for id
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["Id"])

	// query db for specific id, scan to tempTodo
	tempTodo := ToDo{}
	db.QueryRow("SELECT * FROM todos WHERE id=?", key).Scan(&tempTodo.Id, &tempTodo.Name, &tempTodo.Content, &tempTodo.Complete)

	// if tempTodo is blank, throw 400 | else return todo as json
	if (tempTodo == ToDo{}) {
		w.WriteHeader(400)
	} else {
		json.NewEncoder(w).Encode(tempTodo)
	}
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	ToDos := []ToDo{}

	// query database for all todos
	results, e := db.Query("select * from todos")
	printError("error fetching from db", e)

	// iterate over all results, scan to ToDo slice
	for results.Next() {
		tempTodo := ToDo{}
		e = results.Scan(&tempTodo.Id, &tempTodo.Name, &tempTodo.Content, &tempTodo.Complete)
		printError("unable to parse todo", e)
		ToDos = append(ToDos, tempTodo)
	}

	// if slice is empty return 400, else return slice as json
	if len(ToDos) == 0 {
		w.WriteHeader(400)
	} else {
		json.NewEncoder(w).Encode(ToDos)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/todos", returnAll)                         // return all todos
	router.HandleFunc("/todo", createTodo).Methods("POST")         // create todo
	router.HandleFunc("/todos/{Id}", deleteTodo).Methods("DELETE") // delete specific todo
	router.HandleFunc("/todos/{Id}", updateTodo).Methods("PUT")    // update specific todo
	router.HandleFunc("/todos/{Id}", returnOne)                    // return specific todo based on Id
	log.Fatal(http.ListenAndServe(":10000", router))               // start server
}

func main() {
	db, e = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo")
	printError("unable to open db ", e)

	defer db.Close()
	handleRequests()
}

// docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -ti mysql:latest
// mysql -h localhost -P 3306 --protocol=tcp -u root -p
// curl -i -X POST -H "Content-Type: application/json" -d '{"Id" : "3", "Name" : "added todo", "Content" : "this was added with a POST request", "Complete" : false}' http://127.0.0.1:10000/todo
// curl -i -X DELETE http://127.0.0.1:10000/todos/2
// curl -i -X PUT -H "Content-Type: application/json" -d '{"Id" : "1", "Name" : "updated todo", "Content" : "this was updated with a PUT", "Complete" : false}' http://127.0.0.1:10000/todos/1
