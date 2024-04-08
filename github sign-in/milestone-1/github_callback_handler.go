package main

import (
	"log"
	"net/http"
)

func githubCallbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("githubCallbackHandler called")

	reqContext := req.Context()

	//checking OAuthState == state
	cookie, err := req.Cookie(oauthStateCookie)
	stateFromURI := req.URL.Query().Get("state")

	if err != nil || stateFromURI != cookie.Value {
		log.Println("State tokens don't match. Ignoring callback.")
		http.Error(w, "Invalid callback request", http.StatusBadRequest)
		return
	}

	codeFromURI := req.URL.Query().Get("code")
	token, err := oauthConf.Exchange(reqContext, codeFromURI)
	if err != nil {
		log.Println("oAuth exchange error: ", err)
		http.Error(w, "Error logging in.", http.StatusInternalServerError)
		return
	}

	session, err := createSession(reqContext, token.AccessToken)
	if err != nil {
		log.Println("Error creating session: ", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	setCookie(w, sessionCookie, session.ID, sessionCookieMaxAge)
	setCookie(w, oauthStateCookie, "", -1)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}
