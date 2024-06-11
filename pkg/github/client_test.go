package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

const (
	// baseURLPath is a common prefix for all API requests just for testing.
	baseURLPath = "/api-v3"
)

func TestClient_DownloadRelease(t *testing.T) {
	client, server, _, teardown := setup()
	defer teardown()

	releases := buildTestReleases(t)

	server.HandleFunc("/repos/{owner}/{repo}/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars["owner"]
		repo := vars["repo"]
		release := releases[owner+"/"+repo+"/latest"]
		if release == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"message":"Not Found"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(release)
		fmt.Fprintf(w, "%s", b)
	})

	server.HandleFunc("/repos/{owner}/{repo}/releases/tags/{tag}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars["owner"]
		repo := vars["repo"]
		tag := vars["tag"]
		release := releases[owner+"/"+repo+"/"+tag]
		if release == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"message":"Not Found"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(release)
		fmt.Fprintf(w, "%s", b)
	})
	server.HandleFunc("/repos/{owner}/{repo}/releases/assets/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars["owner"]
		repo := vars["repo"]
		id, _ := strconv.Atoi(vars["id"])

		for k, release := range releases {
			if strings.HasPrefix(k, owner+"/"+repo) {
				for _, asset := range release.Assets {
					if asset.Id == id {
						w.WriteHeader(http.StatusOK)
						fmt.Fprintf(w, "this is file content %s/%s/%s", owner, repo, asset.Name)
						return
					}
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"message":"Not Found"}`)
	})

	testcases := []struct {
		name                   string
		owner, repo, tag, file string
		err                    string
		expected               string
	}{
		{
			name:  "latest release",
			owner: "owner1",
			repo:  "repo1",
			tag:   "",
			file:  "example.zip",
		},
		{
			name:  "latest release with file not found",
			owner: "owner1",
			repo:  "repo2",
			tag:   "latest",
			file:  "example.zip",
			err:   `find asset id (404): {"message":"Not Found"}`,
		},
		{
			name:  "specific release",
			owner: "owner1",
			repo:  "repo1",
			tag:   "v0.1.0",
			file:  "example.tar.gz",
		},
		{
			name:  "specific release with file not found",
			owner: "owner1",
			repo:  "repo1",
			tag:   "v0.1.0",
			file:  "example.txt",
			err:   "asset not found",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// create buffer file to store the downloaded file
			f, err := os.CreateTemp("", tt.file)
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(f.Name())
			if err = client.DownloadRelease(tt.owner, tt.repo, tt.tag, tt.file, f); err != nil {
				if tt.err == "" {
					t.Errorf("DownloadRelease() error = %v", err)
					return
				}
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}
			content := fmt.Sprintf("this is file content %s/%s/%s", tt.owner, tt.repo, tt.file)
			b, err := os.ReadFile(f.Name())
			if err != nil {
				t.Fatalf("failed to read file content: %v", err)
			}
			if string(b) != content {
				t.Errorf("unexpected file content: %s", string(b))
			}
		})
	}
}

func buildTestReleases(t *testing.T) map[string]*Release {
	b, err := os.ReadFile("./testdata/releases.json")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}
	ret := make(map[string]*Release)
	if err := json.Unmarshal(b, &ret); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	return ret
}

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, router *mux.Router, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	router = mux.NewRouter()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, router))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("path:", req)
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = NewClient()
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.url = url
	client.token = "ghp_token"
	return client, router, server.URL, server.Close
}
