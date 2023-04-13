package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/fcsgehrke/toggl-deck/docs"
	"github.com/swaggo/http-swagger/v2" // http-swagger middlewarae
)

func Router(service DeckService, log *log.Logger) (http.Handler, error) {
	handler, err := NewDeckHandler(service, log)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route(apiVer(1, "decks"), func(r chi.Router) {
		r.Post("/", handler.CreateDeck)
		r.Get("/{id}", handler.OpenDeck)
		r.Get("/{id}/draw", handler.DrawCard)
	})

	router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/docs/doc.json"), 
  ))

	return router, nil
}

func Start(ctx context.Context, addr string, service DeckService, log *log.Logger) error {

	// Router
	router, err := Router(service, log)
	if err != nil {
		return err
	}

	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: router}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Printf("[INF] - 🚀 Server started at: %s\n", addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}

func apiVer(version int, pattern string) string {
	return fmt.Sprintf("/api/v%d/%s", version, pattern)
}
