package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// use .env file in production
var (
	host       = "127.0.0.1:3301"
	user       = "admin"
	pass       = "pass"
	accessPort = "8085"
)

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	repo := TarantoolRepo{host, user, pass}

	r := mux.NewRouter()
	r.HandleFunc("/post", wrappedCreatePostHandler(&repo)).Methods("Post")
	r.HandleFunc("/comment", wrappedCreateCommentHandler(&repo)).Methods("Post")
	r.HandleFunc("/posts", wrappedReadAllPostsHandler(&repo)).Methods("Get")
	r.HandleFunc("/comments", wrappedReadPostCommentsHandler(&repo)).Methods("Get")
	r.HandleFunc("/reset", wrappedResetHandler(&repo)).Methods("Post")

	protectedRouter := panicMiddleware(r)
	http.Handle("/", protectedRouter)

	fmt.Printf("Server is listening on %s\n", accessPort)
	err := http.ListenAndServe(":"+accessPort, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
