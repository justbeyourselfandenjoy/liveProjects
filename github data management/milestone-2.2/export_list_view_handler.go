package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//go:embed exports_list.html.tmpl
var githubListExportsHtml string

type githubExportListResult struct {
	ID          int64
	CreatedAt   string
	State       string
	DownloadURL string
}

func githubListExportsViewHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("githubListExportsViewHandler has been called")

	var exportListResult []githubExportListResult
	var downloadURL string

	ctx := req.Context()

	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}

	s := initiateLoginIfRequired(w, req)
	if s == nil {
		return
	}

	token := sessionsStore[s.ID].accessToken
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx = context.WithValue(ctx, oauth2.HTTPClient, oauthHttpClient)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	migrations, response, err := ghClient.Migrations.ListUserMigrations(ctx)
	if err != nil {
		log.Println("Error listing user migrations: ", err)
		if response != nil {
			defer response.Body.Close()
			if responseData, err := io.ReadAll(response.Body); err == nil {
				log.Println(string(responseData))
			}
		}
		http.Error(w, "Error listing user migrations", http.StatusInternalServerError)
		return
	}
	for _, m := range migrations {
		log.Printf("Got migration state = %s for ID = %d\n", *m.State, *m.ID)
		if *m.State == "exported" {
			downloadURL, err = ghClient.Migrations.UserMigrationArchiveURL(ctx, *m.ID)
			if err != nil {
				log.Println("Error retrieving migration archive URL", err)
			}
		}
		createdAt, err := time.Parse(time.RFC3339, *m.CreatedAt)
		if err != nil {
			log.Println("Error parsing createdAt time string", *m.CreatedAt)
		}
		createdYear, createdMonth, createdDay := createdAt.Date()
		exportListResult = append(
			exportListResult,
			githubExportListResult{
				ID:          *m.ID,
				CreatedAt:   fmt.Sprintf("%d-%s-%d", createdYear, createdMonth, createdDay),
				State:       *m.State,
				DownloadURL: downloadURL,
			},
		)
	}

	exportsListHtmlTmpl, err := template.New("github_export_view_html").Parse(githubListExportsHtml)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}
	err = exportsListHtmlTmpl.Execute(w, exportListResult)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}
}
