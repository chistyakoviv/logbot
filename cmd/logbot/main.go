package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("Starting http server")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			fmt.Println(fmt.Errorf("write error: %v", err))
		}
	})
	srv := &http.Server{
		Addr:              ":80",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(fmt.Errorf("http server error: %v", err))
	}
	fmt.Println("http server stopped")
}
