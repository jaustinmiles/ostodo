package main

import (
	"context"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger()
	_, err := w.Write([]byte("pong\n"))
	if err != nil {
		l.Errorf("couldn't write response to client: %v", err)
	}
	l.Info("responded to client ", r.UserAgent())
}

func main() {
	l := getLogger()
	l.Info("creating server")
	// CORS
	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	serveMutex := mux.NewRouter()

	getRouter := serveMutex.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/ping", defaultHandler)

	server := http.Server{
		Addr:         ":9090",
		Handler:      cors(serveMutex),
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	l.Info("starting server")

	// start server
	go func() {
		l.Info("starting server on port 9090")
		err := server.ListenAndServe()
		if err != nil {
			l.Errorf("Error starting server: %v", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		l.Warnf("error shutting down server, %v", err)
	}
}
