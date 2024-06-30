package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"milestone-2.2/testutils"

	"golang.org/x/net/publicsuffix"
)

func TestListExportsView(t *testing.T) {

	var err error

	expectedResult := []githubExportListResult{
		{
			ID:        79,
			State:     "pending",
			CreatedAt: "2022-July-6",
		},
		{
			ID:          99,
			State:       "exported",
			CreatedAt:   "2022-July-8",
			DownloadURL: "https://example.com/myfile.tar.gz",
		},
	}

	exportsListHtmlTmpl, err := template.New("github_export_view_html").Parse(githubListExportsHtml)
	if err != nil {
		t.Fatal(err)
	}
	var expectedHtml bytes.Buffer
	err = exportsListHtmlTmpl.Execute(&expectedHtml, expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	registerHandlers(mux)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	initOAuthConfig(testutils.GetenvStub)
	oauthHttpClient = testutils.HttpClientWithGithubStub(ts.URL)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatal(err)
	}

	testHTTPClient := testutils.HttpClientWithGithubStub(ts.URL)
	testHTTPClient.Jar = jar

	// this call sets up the necessary server side session data as well
	// client side cookies so that we satisfy the auth requirements
	_, err = testHTTPClient.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := testHTTPClient.Get(ts.URL + "/github/exports")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status: %v, Got: %v", http.StatusOK, resp.StatusCode)
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body")
	}

	if string(respData) != expectedHtml.String() {
		t.Fatalf("Expected: %s\nGot: %s", expectedHtml.String(), string(respData))
	}
}
