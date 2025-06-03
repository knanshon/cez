package main

import (
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	ID      int
	Message string
}

func main() {

	data := map[string][]Todo{
		"Todos": {
			Todo{ID: 1, Message: "Buy Milk"},
		},
	}

	todosHandler := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		templ.Execute(w, data)
	}

	addTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		message := r.PostFormValue("message")
		templ := template.Must(template.ParseFiles("index.html"))
		todo := Todo{ID: len(data["Todos"]) + 1, Message: message}
		data["Todos"] = append(data["Todos"], todo)
		templ.ExecuteTemplate(w, "todo-list-element", todo)
	}

	// API handlers
	http.HandleFunc("/", todosHandler)
	http.HandleFunc("/add-todo", addTodoHandler)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":9090", nil))
}
