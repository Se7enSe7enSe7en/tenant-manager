package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/env"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/tenant"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	connStr string
}

type AppConfig struct {
	dbConfig DBConfig
}

func main() {
	// constants
	const port = "8080"
	const proxy_port = "7331"

	// global context
	ctx := context.Background()

	// load env variables
	godotenv.Load()

	// app config
	cfg := AppConfig{
		dbConfig: DBConfig{
			connStr: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=db sslmode=disable"),
		},
	}

	// connect to database
	conn, err := pgx.Connect(ctx, cfg.dbConfig.connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer conn.Close(ctx)

	logger.Debug("connected to database: %v", cfg.dbConfig.connStr)

	// create server mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./web/static/assets"))
	mux.Handle("/assets/", utils.DisableCacheInDevMode(http.StripPrefix("/assets/", fs)))

	// page handlers
	mux.HandleFunc("/", tenant.ListTenantPage)

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
