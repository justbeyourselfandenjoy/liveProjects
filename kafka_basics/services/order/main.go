package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	startTime := time.Now()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Got %v signal. Exiting. Uptime %v\n", sig.String(), time.Since(startTime).String())
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /order/{$}", orderHandler)

	log.Println("Starting server at " + serverIP + ":" + serverPort + " ... ")
	log.Panicln(http.ListenAndServe(serverIP+":"+serverPort, mux))
}
