package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"milestone-2.2/testutils"

	"golang.org/x/net/publicsuffix"
)

func TestExportNewApiHandler(t *testing.T) {

	var expectedResponse = githubExportResult{
		ID:           79,
		State:        "pending",
		Repositories: []string{"test-user-1/test-repo-1"},
	}
	var gotResponse githubExportResult
	var err error

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

	// we make a request to the view so that the session store
	// is setup correctly and the expected cookies have been
	// set
	_, err = testHTTPClient.Get(ts.URL + "/github/export/new")
	if err != nil {
		t.Fatal(err)
	}

	emptyBody := strings.NewReader("")
	resp, err := testHTTPClient.Post(ts.URL+"/api/github/export", "application/json", emptyBody)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(respData, &gotResponse)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotResponse, expectedResponse) {
		t.Fatalf("Expected: %#v, Got: %#v\n", expectedResponse, gotResponse)
	}
}
