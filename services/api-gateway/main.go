package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /trip/preview", enabledCORS(handleTripPreview))
	mux.HandleFunc("/ws/drivers", handleDriverWebSocket)
	mux.HandleFunc("/ws/riders", handleRiderWebSocket)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Server listening on %s", httpAddr)
		serverError <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("sever error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("server is shutting down due to %v signal ", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Print("could not stop server gracefully ")
			server.Close()
		}

	}

}
