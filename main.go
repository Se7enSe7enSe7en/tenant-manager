package main

import (
	"context"
	"log"
	"net/http"

	web "github.com/Se7enSe7enSe7en/tenant-manager/web/templ"
)

const port = "8080"

func main() {
	// create server mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./web/static/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.MainPage().Render(context.Background(), w)
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
