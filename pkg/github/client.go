package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
)

var ErrAssetNotFound = fmt.Errorf("asset not found")

const (
	GitHubApi = "https://api.github.com"
)

type Client struct {
	token  string
	silent bool
	url    *url.URL
	client *http.Client
}

func NewClient(opt ...Option) *Client {
	client := &Client{client: http.DefaultClient}
	client.url, _ = url.Parse(GitHubApi)
	for _, o := range opt {
		o(client)
	}
	return client
}

func (c *Client) DownloadRelease(owner, repo, tag, file string, out *os.File) error {
	assetId, err := c.findassetId(owner, repo, tag, file)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("repos/%s/%s/releases/assets/%d", owner, repo, assetId)
	req, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	resp := doRequest(c.client, req, c.silent)
	if resp.code != http.StatusOK {
		return fmt.Errorf("failed to download release (%d): %w", resp.code, resp.err)
	}

	if out == nil {
		return nil
	}

	if _, err := io.Copy(out, resp.body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Client) findassetId(owner, repo, tag, file string) (int, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/releases", owner, repo)
	if tag == "" || tag == "latest" {
		endpoint = path.Join(endpoint, "latest")
	} else {
		endpoint = path.Join(endpoint, "tags", tag)
	}

	req, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return -1, fmt.Errorf("new request: %w", err)
	}

	resp := doRequest(c.client, req, true)
	if resp.code != http.StatusOK {
		return -1, fmt.Errorf("find asset id (%d): %w", resp.code, resp.err)
	}

	var release Release
	if err := json.NewDecoder(resp.body).Decode(&release); err != nil {
		return -1, fmt.Errorf("unmarshal response: %w", err)
	}

	for _, asset := range release.Assets {
		if asset.Name == file {
			return asset.Id, nil
		}
	}

	return -1, ErrAssetNotFound
}

func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.url.JoinPath(endpoint).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	}
	return req, nil
}

type response struct {
	code        int
	err         error
	contentType string
	body        io.ReadCloser
}

func doRequest(client *http.Client, req *http.Request, silent bool) response {
	var ret response
	resp, err := client.Do(req)
	if err != nil {
		ret.code = http.StatusInternalServerError
		ret.err = fmt.Errorf("failed to do request: %w", err)
		return ret
	}
	defer resp.Body.Close()

	var pgFn func(maxBytes int64, prefix ...string) *progressbar.ProgressBar
	if silent {
		pgFn = progressbar.DefaultBytesSilent
	} else {
		pgFn = progressbar.DefaultBytes
	}

	var buf bytes.Buffer

	bar := pgFn(resp.ContentLength, "Downloading...")
	io.Copy(io.MultiWriter(&buf, bar), resp.Body)

	ret.body = io.NopCloser(&buf)
	ret.code = resp.StatusCode
	ret.contentType = resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		b := buf.Bytes()
		if len(b) > 0 {
			ret.err = fmt.Errorf("%s", b)
		} else {
			ret.err = fmt.Errorf("failed to do request: %s", resp.Status)
		}
	}
	return ret
}
