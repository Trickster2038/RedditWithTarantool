package main

import (
	"log"
	"net/http"
	"os"
	"redditclone/pkg/comment"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/post"
	"redditclone/pkg/user"
	"redditclone/pkg/vote"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/html/index.html")
}

func main() {
	zapLogger, err := zap.NewProduction()
	defer func() {
		err = zapLogger.Sync()
		if err != nil {
			log.Panicf("error while logger sync: %v", err)
		}
	}()

	if err != nil {
		log.Panicf("error while logger init: %v", err)
	}

	logger := zapLogger.Sugar()

	var postsRepo post.PostRepo
	var userRepo user.UserRepo
	var commentRepo comment.CommentRepo
	var voteRepo vote.VoteRepo
	args := os.Args
	if (len(args) > 1) && (args[1] == "testmode") {
		logger.Info("Test mode run")
		postsRepo = post.NewMockMemoryRepo()
		userRepo = user.NewMockMemoryRepo()
		commentRepo = comment.NewMockMemoryRepo()
		voteRepo = vote.NewMockMemoryRepo()
	} else {
		logger.Info("Production mode run")
		postsRepo = post.NewMemoryRepo()
		userRepo = user.NewMemoryRepo()
		commentRepo = comment.NewMemoryRepo()
		voteRepo = vote.NewMemoryRepo()
	}

	userHandler := &handlers.UserHandler{
		UserRepo: userRepo,
		Logger:   logger,
	}

	postsHandler := &handlers.PostsHandler{
		Logger:      logger,
		PostRepo:    postsRepo,
		CommentRepo: commentRepo,
		VoteRepo:    voteRepo,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/a/{path:.*}", index).Methods("GET") // no 404 on refresh

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static",
		http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/api/posts/", postsHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Authorize).Methods("POST")

	r.HandleFunc("/api/posts", postsHandler.Add).Methods("POST")

	/*
		ERROR IN redditclone.md:
		this is true route for "GET /api/funny/{CATEGORY_NAME}"
	*/
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postsHandler.GetAllInCategory).Methods("GET")

	r.HandleFunc("/api/post/{POST_ID}", postsHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postsHandler.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postsHandler.DeleteComment).
		Methods("DELETE")
	r.HandleFunc("/api/post/{POST_ID}/upvote", postsHandler.Upvote).
		Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/downvote", postsHandler.Downvote).
		Methods("GET")

	// NOT MENTIONED IN redditclone.md, BUT NECESSARY
	r.HandleFunc("/api/post/{POST_ID}/unvote", postsHandler.Unvote).
		Methods("GET")

	r.HandleFunc("/api/post/{POST_ID}", postsHandler.Delete).
		Methods("DELETE")
	r.HandleFunc("/api/user/{USER_LOGIN}", postsHandler.GetAllByUser).
		Methods("GET")

	r.Use(middleware.AccessLogMiddleware(logger))
	r.Use(middleware.PanicMiddleware(logger))

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Panicf("error while serving port: %v", err)
	}
}
