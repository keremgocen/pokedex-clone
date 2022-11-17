package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pokedex-clone/pkg/api"
	"pokedex-clone/pkg/pokemon"
	"pokedex-clone/pkg/storage"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ctxTimeout         = 5 * time.Second
	serverTimeout      = 3 * time.Second
	pokeAPIURL         = "https://pokeapi.co/api/v2/pokemon-species/"
	translationsAPIURL = "https://api.funtranslations.com/translate/"
)

func main() {
	storageAPI := storage.NewStore()
	pokeAPIClient := api.NewClient(pokeAPIURL, serverTimeout)
	pokeAPI := api.Poke{
		Client: pokeAPIClient,
	}
	translationsAPIClient := api.NewClient(translationsAPIURL, serverTimeout)
	translationsAPI := api.Translations{
		Client: translationsAPIClient,
	}
	service := pokemon.NewService(storageAPI, pokeAPI, translationsAPI)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/pokemon/:name", service.Get)
	router.GET("/pokemon/translated/:name", service.GetTranslated)

	httpServer := &http.Server{
		Addr:              ":5000",
		Handler:           router,
		ReadHeaderTimeout: serverTimeout,
		WriteTimeout:      serverTimeout,
		ReadTimeout:       serverTimeout,
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
