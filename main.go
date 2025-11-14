package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/benbjohnson/hashfs"
)

const port = "8080"

//go:embed web/static/*
var staticFS embed.FS

// hashed static FS, this is for http caching
var staticSys = hashfs.NewFS(staticFS)

func main() {
	// create server mux
	mux := http.NewServeMux()

	// register handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, _ := fs.Sub(staticSys, "web/static")
		http.FileServer(http.FS(content)).ServeHTTP(w, r)
	})

	// init server
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start server
	log.Printf("Open browser to: http://localhost:%v", port)
	log.Fatalln(s.ListenAndServe())
}
