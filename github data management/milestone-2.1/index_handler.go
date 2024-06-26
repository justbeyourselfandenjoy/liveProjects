package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	s, err := getSession(req, sessionCookie)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("Cookie '" + sessionCookie + "' is not set")
		default:
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		stateToken, err := getRandomString()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		githubLoginUrl := oauthConf.AuthCodeURL(stateToken)
		setCookie(w, oauthStateCookie, stateToken, 600)
		http.Redirect(w, req, githubLoginUrl, http.StatusTemporaryRedirect)
		return
	}
	log.Printf("Successfully authorized to access GitHub on your behalf: %s\n", sessionsStore[s.ID].Login)
	fmt.Fprintf(w, "Successfully authorized to access GitHub on your behalf: %s", sessionsStore[s.ID].Login)

	/*	if _, ok := sessionsStore[cookie.Value]; ok {
			userData := sessionsStore[cookie.Value]
			written, err := w.Write([]byte("Successfully authorized to access GitHub on your behalf: " + userData.Login))
			if err != nil {
				log.Println("Can't write to http.ResponseWriter")
				return
			}
			log.Printf("%v bytes sent back \n", written)
		}
	*/
}
