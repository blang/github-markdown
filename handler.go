package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var css_body = `
	.markdown-body {
        min-width: 200px;
        max-width: 790px;
        margin: 0 auto;
        padding: 30px;
    }
`

type HttpHandler struct {
	Directory string
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := path.Join(h.Directory, r.RequestURI)
	f, err := os.Open(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File error: %s\n", err)
		return
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File error: %s\n", err)
		return
	}
	m := make(map[string]string)
	m["mode"] = "markdown"
	m["text"] = string(b)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	resp, err := http.Post("https://api.github.com/markdown", "application/json", &buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	fmt.Fprintf(w, "<html><head><title>%s</title><link rel=\"stylesheet\" href=\"http://sindresorhus.com/github-markdown-css/github-markdown.css\"><style>%s</style></head><body><article class=\"markdown-body\">\n", name, css_body)
	fmt.Fprintf(w, "%s", b)
	fmt.Fprintf(w, "</article></body></html>")
	fmt.Printf("200 - Request: %q Local: %q\n", r.RequestURI, name)
}
