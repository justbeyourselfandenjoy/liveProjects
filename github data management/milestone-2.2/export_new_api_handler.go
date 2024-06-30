package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubExportResult struct {
	ID           int64    `json:"id"`
	State        string   `json:"state"`
	Repositories []string `json:"repositories"`
}

func githubNewExportApiHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("githubNewExportApiHandler has been called")

	var responseData []byte
	var repoToExportFullNames, repoExportedFullNames []string

	ctx := req.Context()

	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	s, err := getSession(req, sessionCookie)
	if err != nil {
		http.Error(w, "session cookie invalid or not found", http.StatusUnauthorized)
		return
	}

	token := sessionsStore[s.ID].accessToken
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	ctx = context.WithValue(ctx, oauth2.HTTPClient, oauthHttpClient)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	listOptions := &github.RepositoryListOptions{Affiliation: "owner", ListOptions: github.ListOptions{Page: 1}, Sort: "updated", Direction: "desc"}
	for {
		repos, response, err := ghClient.Repositories.List(ctx, sessionsStore[s.ID].Login, listOptions)
		if err != nil {
			log.Println("Error retrieving user's repositories to export: ", err)
			if response != nil {
				if responseData, err := io.ReadAll(response.Body); err == nil {
					log.Println(string(responseData))
				}
				response.Body.Close()
			}
			http.Error(w, "Error retrieving user's repositories to export", http.StatusInternalServerError)
			return
		}
		for _, repo := range repos {
			if *repo.Fork {
				continue
			}
			repoToExportFullNames = append(repoToExportFullNames, *repo.FullName)
		}

		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage
	}
	log.Println("repoToExportFullNames=", repoToExportFullNames)

	migration, response, err := ghClient.Migrations.StartUserMigration(ctx, repoToExportFullNames, &github.UserMigrationOptions{ExcludeAttachments: true})
	if err != nil {
		log.Println("Error starting user migration: ", err)
		if response != nil {
			if responseData, err := io.ReadAll(response.Body); err == nil {
				log.Println(string(responseData))
			}
			response.Body.Close()
		}
		http.Error(w, "Error starting user repos migration", http.StatusInternalServerError)
		return
	}

	for _, r := range migration.Repositories {
		repoExportedFullNames = append(repoExportedFullNames, *r.FullName)
	}

	exportResult := githubExportResult{
		ID:           *migration.ID,
		State:        *migration.State,
		Repositories: repoExportedFullNames,
	}

	responseData, err = json.Marshal(&exportResult)
	if err != nil {
		log.Println("Error marshalling start migration response: ", err)
		http.Error(w, "Error marshalling start migration response", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, string(responseData))
}
