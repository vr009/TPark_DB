package main

import (
	"context"
	"fmt"
	"forum/internal/delivery"
	repo2 "forum/internal/repo"
	usecase2 "forum/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	srv := http.Server{Handler: r, Addr: fmt.Sprintf(":%s", "5000")}
	//conn := "postgres://postgres:password@127.0.0.1:5432/docker?pool_max_conns=100"
	conn := "host=127.0.0.1 port=5432 user=docker password=docker dbname=docker"

	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	repo := repo2.NewForumRepo(pool)
	usecase := usecase2.NewForumUsecase(repo)
	handler := delivery.NewForumHandler(usecase)

	forum := r.PathPrefix("/api/forum").Subrouter()
	{
		forum.HandleFunc("/create", handler.CreateForum).Methods("POST")
		forum.HandleFunc("/{slug}/details", handler.GetForum).Methods("GET")
		forum.HandleFunc("/{slug}/create", handler.CreateThread).Methods("POST")
		forum.HandleFunc("/{slug}/users", handler.GetUsers).Methods("GET")
		forum.HandleFunc("/{slug}/threads", handler.GetThreads).Methods("GET")
	}

	post := r.PathPrefix("/api/post").Subrouter()
	{
		post.HandleFunc("/{id}/details", handler.GetThreadInfo).Methods("GET")
		post.HandleFunc("/id/details", handler.UpdateMessage).Methods("POST")
	}

	service := r.PathPrefix("/api/service").Subrouter()
	{
		service.HandleFunc("/clear", handler.DropAllInfo).Methods("POST")
		service.HandleFunc("/status", handler.GetStatus).Methods("GET")
	}

	thread := r.PathPrefix("/api/thread").Subrouter()
	{
		thread.HandleFunc("/{slug_or_id}/create", handler.CreatePosts).Methods("POST")
		thread.HandleFunc("/{slug_or_id}/details", handler.GetThreadInfoBySlug).Methods("GET")
		thread.HandleFunc("/{slug_or_id}/details", handler.UpdateThread).Methods("POST")
		thread.HandleFunc("/{slug_or_id}/posts", handler.GetPosts).Methods("GET")
		thread.HandleFunc("/{slug_or_id}/vote", handler.VoteForThread).Methods("POST")
	}

	user := r.PathPrefix("/api/user").Subrouter()
	{
		user.HandleFunc("/{nickname}/create", handler.CreateProfile).Methods("POST")
		user.HandleFunc("/{nickname}/profile", handler.GetProfile).Methods("GET")
		user.HandleFunc("/{nickname}/profile", handler.UpdateProfile).Methods("POST")
	}

	http.Handle("/", r)
	log.Print("main running on: ", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
