package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	serverIP   = "localhost"
	serverPort = "8080"
)

var (
	oauthConf *oauth2.Config
)

func initOAuthConfig() {
	if len(os.Getenv("CLIENT_ID")) == 0 {
		log.Fatal("Must specify your app's CLIENT_ID")
	}

	if len(os.Getenv("CLIENT_SECRET")) == 0 {
		log.Fatal("Must specify your app's CLIENT_SECRET")
	}

	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"repo", "user"},
		Endpoint:     github.Endpoint,
	}
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/github/callback", githubCallbackHandler)
}

func main() {

	initOAuthConfig()

	mux := http.NewServeMux()
	registerHandlers(mux)

	log.Println("Starting server at " + serverIP + ":" + serverPort + " ... ")
	log.Fatal(http.ListenAndServe(serverIP+":"+serverPort, mux))
}
