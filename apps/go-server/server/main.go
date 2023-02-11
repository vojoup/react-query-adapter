package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	webrpcHandler := NewExampleServer(&ExampleRPC{})
	r.Handle("/*", webrpcHandler)

	fmt.Println("Listening on 0.0.0.0:4242...")

	return http.ListenAndServe(":4242", r)
}

type ExampleRPC struct {
}

func (s *ExampleRPC) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *ExampleRPC) GetUser(ctx context.Context, userID uint64) (*User, error) {
	if userID == 911 {
		return nil, ErrorNotFound("unknown userID %d", 911)
		// return nil, Errorf(ErrNotFound, "unknown userID %d", 911)
		// return nil, WrapError(ErrNotFound, err, "unknown userID %d", 911)
	}

	return &User{
		ID:       userID,
		Username: "hihi",
		Meta:     map[string]interface{}{"location": "Toronto"},
	}, nil
}

func (s *ExampleRPC) FindUsers(ctx context.Context, q string) (*Page, []*User, error) {
	page := &Page{Num: 1}

	users := []*User{
		{ID: 1, Username: "a", Meta: map[string]interface{}{"location": "Montreal"}},
		{ID: 2, Username: "b", Meta: map[string]interface{}{"age": 10}},
	}

	return page, users, nil
}
