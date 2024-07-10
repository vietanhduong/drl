package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseGitHubDownloadUrl(t *testing.T) {
	testcases := []struct {
		url    string
		expect *release
	}{
		{
			url: "https://github.com/vietanhduong/drl/releases/download/v0.1.0/drl_darwin_amd64",
			expect: &release{
				owner: "vietanhduong",
				repo:  "drl",
				tag:   "v0.1.0",
				asset: "drl_darwin_amd64",
			},
		},
		{
			url: "https://github.com/vietanhduong/drl/releases/latest/download/drl_darwin_amd64",
			expect: &release{
				owner: "vietanhduong",
				repo:  "drl",
				tag:   "latest",
				asset: "drl_darwin_amd64",
			},
		},
		{"nomatch", nil},
	}

	for _, tt := range testcases {
		actual := parseGitHubDownloadUrl(tt.url)
		if tt.expect == nil {
			assert.Nil(t, actual)
			continue
		}
		if actual == nil {
			t.Errorf("parseGitHubDownloadUrl(%s) = nil, want %v", tt.url, *tt.expect)
			continue
		}
		assert.Equal(t, *tt.expect, *actual)
	}
}
