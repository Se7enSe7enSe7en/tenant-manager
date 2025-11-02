package main

import (
	"log"
	"net/http"
)

const port = "8080"
const fileServerRootPath = "web/static/"

func main() {
	// create server mux
	mux := http.NewServeMux()

	// register handlers
	fsHandler := http.FileServer(http.Dir(fileServerRootPath))
	mux.Handle("/", fsHandler)

	// init server
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start server
	log.Printf("Serving in: http://localhost:%v", port)
	log.Fatalln(s.ListenAndServe())
}
