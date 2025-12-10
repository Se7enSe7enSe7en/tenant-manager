package main

import (
	"context"
	"log"
	"net/http"

	web "github.com/Se7enSe7enSe7en/tenant-manager/web/static"
)

const port = "8080"

func main() {
	// create server mux
	mux := http.NewServeMux()

	// register handlers
	// mux.Handle("/", templ.Handler(web.MainPage()))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.MainPage().Render(context.Background(), w)
	})

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	content, _ := fs.Sub(staticSys, "web/static")
	// 	http.FileServer(http.FS(content)).ServeHTTP(w, r)
	// })

	// init server
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start server
	log.Printf("Open browser to: http://localhost:%v", port)
	log.Fatalln(s.ListenAndServe())
}
