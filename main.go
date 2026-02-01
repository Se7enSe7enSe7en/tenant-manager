package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	web "github.com/Se7enSe7enSe7en/tenant-manager/web/components"
)

const port = "8080"
const proxy_port = "7331"

func main() {
	// create server mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./web/static/assets"))
	mux.Handle("/assets/", utils.DisableCacheInDevMode(http.StripPrefix("/assets/", fs)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.MainPage().Render(context.Background(), w)
	})

	// init server
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start server
	logger.Debug("For Production, open the actual port: http://localhost:%v", port)
	logger.Debug("For Development, open the proxy port (for templ hot reload): http://localhost:%v", proxy_port)

	log.Fatalln(s.ListenAndServe())
}
