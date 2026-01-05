package main

import (
	"context"
	"log"
	"net/http"

	web "github.com/Se7enSe7enSe7en/tenant-manager/web/components"
)

const port = "8080"
const proxy_port = "7331"

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
	log.Printf("Actual Port: http://localhost:%v", port)
	log.Printf("Proxy Port (Templ hot reload): http://localhost:%v", proxy_port)
	log.Fatalln(s.ListenAndServe())
}
