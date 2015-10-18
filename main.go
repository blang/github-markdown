package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)
import "net/http"

var dir = flag.String("dir", "./", "Directory to serve")
var listen = flag.String("listen", "127.0.0.1:8080", "Address to listen on")

func main() {
	flag.Parse()

	handledir := path.Clean(*dir)
	f, err := os.Open(handledir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening directory: %s\n", err)
		f.Close()
		os.Exit(1)
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening directory: %s\n", err)
		f.Close()
		os.Exit(1)
	}
	if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: No directory given\n")
		f.Close()
		os.Exit(1)
	}
	f.Close()

	handler := &HttpHandler{
		Directory: handledir,
	}

	http.Handle("/", handler)
	fmt.Printf("Listening on %s\n", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
