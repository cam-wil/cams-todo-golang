package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ToDo struct {
	Id       int    `json:"Id"`
	Name     string `json:"Name"`
	Content  string `json:"content"`
	Complete bool   `json:"complete"`
}

type Error struct {
	Error string `json:"error"`
}

// func deleteTodo(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("got DELETE")
// 	vars := mux.Vars(r)
// 	id := vars["Id"]

// 	for index, todo := range ToDos {
// 		if todo.Id == id {
// 			ToDos = append(ToDos[:index], ToDos[index+1:]...)
// 		}
// 	}
// }

// func updateTodo(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("got PUT")
// 	vars := mux.Vars(r)
// 	id := vars["Id"]

// 	for index, todo := range ToDos {
// 		if todo.Id == id {
// 			requestBody, _ := ioutil.ReadAll(r.Body)
// 			var tempTodo ToDo
// 			json.Unmarshal(requestBody, &tempTodo)
// 			ToDos[index] = tempTodo
// 		}
// 	}
// }

//func createTodo(w http.ResponseWriter, r *http.Request) {
// fmt.Println("got POST")
// requestBody, _ := ioutil.ReadAll(r.Body)
// var todo ToDo
// json.Unmarshal(requestBody, &todo)
// ToDos = append(ToDos, todo)
// json.NewEncoder(w).Encode(todo)
// db, e := sql.Open("mysql", "todo:todopassword@/todos")
// ErrorCheck(e)
// defer db.Close()
// PingDB(db)
//}

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
	router.HandleFunc("/", homePage)       // home page
	router.HandleFunc("/todos", returnAll) // return all todos
	//router.HandleFunc("/todo", createTodo).Methods("POST") // return all todos
	//router.HandleFunc("/todos/{Id}", deleteTodo).Methods("DELETE") // delete specific todo
	//router.HandleFunc("/todos/{Id}", updateTodo).Methods("PUT")    // delete specific todo
	router.HandleFunc("/todos/{Id}", returnOne)      // return specific todo based on Id
	log.Fatal(http.ListenAndServe(":10000", router)) // start server
}

func main() {
	handleRequests()
}

// curl -i -X POST -H "Content-Type: application/json" -d '{"Id" : "3", "Name" : "added todo", "Content" : "this was added with a POST request", "Complete" : false}' http://127.0.0.1:10000/todo
// curl -i -X DELETE http://127.0.0.1:10000/todos/2
// curl -i -X PUT -H "Content-Type: application/json" -d '{"Id" : "1", "Name" : "updated todo", "Content" : "this was updated with a PUT", "Complete" : false}' http://127.0.0.1:10000/todos/1
