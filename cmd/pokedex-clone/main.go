package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pokedex-clone/internal/pokemon"
	"pokedex-clone/internal/storage"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const (
	ctxTimeout        = 5 * time.Second
	readHeaderTimeout = 3 * time.Second
)

func main() {
	storageAPI := storage.NewStore()
	service := pokemon.NewService(storageAPI)

	r := mux.NewRouter()
	r.Handle("/pokemon/{name}", mw(http.HandlerFunc(service.Get))).Methods(http.MethodGet)

	httpServer := &http.Server{
		Addr:              ":5000",
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)
	go func() {
		log.Printf("pokedex-clone service is starting")
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-signalChan:
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		log.Println("server shutdown initiated")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
}

func mw(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}
