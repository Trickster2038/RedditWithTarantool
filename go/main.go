package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln("No .env file found")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	serverPort := os.Getenv("SERVER_PORT")
	repo := TarantoolRepo{host, user, pass}

	r := mux.NewRouter()
	r.HandleFunc("/post", wrappedCreatePostHandler(&repo)).Methods("Post")
	r.HandleFunc("/comment", wrappedCreateCommentHandler(&repo)).Methods("Post")
	r.HandleFunc("/posts", wrappedReadAllPostsHandler(&repo)).Methods("Get")
	r.HandleFunc("/comments", wrappedReadPostCommentsHandler(&repo)).Methods("Get")
	r.HandleFunc("/reset", wrappedResetHandler(&repo)).Methods("Post")

	protectedRouter := panicMiddleware(r)
	http.Handle("/", protectedRouter)

	fmt.Printf("Server is listening on %s\n", serverPort)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
