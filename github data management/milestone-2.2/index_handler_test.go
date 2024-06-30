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

func TestIndexHandler(t *testing.T) {

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
	testHTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		t.Logf("via: %#+v\n", via)
		t.Logf("redirect to: %s\n", req.URL)
		t.Log()
		return nil
	}

	resp, err := testHTTPClient.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var expectedResponse bytes.Buffer
	indexHtmlTmpl, err := template.New("indexhtml").Parse(indexHtml)
	if err != nil {
		t.Fatal(err)
	}
	expectedLoginData := userData{
		Login: "test-user-1",
	}
	err = indexHtmlTmpl.Execute(&expectedResponse, expectedLoginData)
	if err != nil {
		t.Fatal(err)
	}

	if string(respBytes) != expectedResponse.String() {
		t.Fatalf("Expected: %s, Got: %s", expectedResponse.String(), string(respBytes))
	}

}
