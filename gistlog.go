package gistlog

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	gistApiBaseUrl = "https://api.github.com/gists/"
)

// Gist represents a GitHub's gist.
type gist struct {
	Files map[string]gistFile `json:"files,omitempty"`
}

type gistFile struct {
	Filename string `json:"filename,omitempty"`
	Content  string `json:"content,omitempty"`
}

type log struct {
	gistId string
	token  string
}

type createGistFn func(string) (gist, error)

func NewLog(gistId string, token func() string) *log {
	return &log{
		gistId: gistId,
		token:  token(),
	}
}

func (g *log) Insert(filename string, data []string) error {
	return g.updateFile(filename, g.buildGist(data))
}

func (g *log) InsertAsync(filename string, data []string) {
	go func() {
		_ = g.updateFile(filename, g.buildGist(data))
	}()
}

func (g *log) Read(filename string) ([][]string, error) {
	read, err := g.read(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(strings.NewReader(read))
	return reader.ReadAll()
}

func (g *log) read(filename string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, gistApiBaseUrl+g.gistId, nil)
	if err != nil {
		return "", err
	}
	setHeaders(req, g.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkHTTPResponse(resp); err != nil {
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	var gist gist
	err = decoder.Decode(&gist)
	if err != nil {
		return "", err
	}
	return gist.Files[filename].Content, nil
}

func (g *log) buildGist(data []string) createGistFn {
	return func(filename string) (gist, error) {

		prevContent, err := g.read(filename)
		if err != nil {
			return gist{}, err
		}

		buf := &bytes.Buffer{}
		w := csv.NewWriter(buf)
		err = w.Write(data)
		w.Flush()

		if err != nil {
			return gist{}, err
		}

		content := buf.String()

		return gist{
			Files: map[string]gistFile{filename: {
				Filename: filename,
				Content:  prevContent + "\n" + content,
			}},
		}, nil
	}
}

func (g *log) updateFile(filename string, fn createGistFn) error {

	var buf io.ReadWriter

	buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	gist, err := fn(filename)
	if err != nil {
		return err
	}

	err = enc.Encode(gist)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, gistApiBaseUrl+g.gistId, buf)
	if err != nil {
		return err
	}

	setHeaders(req, g.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkHTTPResponse(resp)
}

func setHeaders(req *http.Request, token string) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+token)
}

func checkHTTPResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("HTTP Status: %d, %s", resp.StatusCode, string(b))
	}
	return nil
}
