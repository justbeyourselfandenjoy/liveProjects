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
	oauthConf       *oauth2.Config
	oauthHttpClient *http.Client
)

func initOAuthConfig(getEnvValue func(string) string) {
	if len(getEnvValue("CLIENT_ID")) == 0 {
		log.Fatal("Must specify your app's CLIENT_ID")
	}

	if len(getEnvValue("CLIENT_SECRET")) == 0 {
		log.Fatal("Must specify your app's CLIENT_SECRET")
	}

	oauthConf = &oauth2.Config{
		ClientID:     getEnvValue("CLIENT_ID"),
		ClientSecret: getEnvValue("CLIENT_SECRET"),
		Scopes:       []string{"repo", "user"},
		Endpoint:     github.Endpoint,
	}
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/github/callback", githubCallbackHandler)
}

func main() {

	initOAuthConfig(os.Getenv)
	oauthHttpClient = &http.Client{}

	mux := http.NewServeMux()
	registerHandlers(mux)

	log.Println("Starting server at " + serverIP + ":" + serverPort + " ... ")
	log.Panicln(http.ListenAndServe(serverIP+":"+serverPort, mux))
}
