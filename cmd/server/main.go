package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/env"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/handler"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// think of config as constants, variables inside do not change over time, and is initialized only before starting the application
type config struct {
	port       string // our main port (or address) for prod
	proxy_port string // the port will be used for development with templ
	dsn        string // Data Source Name, or database connection string
}

// think of the application as states, variables inside change over time
type application struct {
	db *pgx.Conn // db driver
}

func main() {
	// global context
	ctx := context.Background()

	// load env variables
	godotenv.Load()

	// init app config
	cfg := config{
		port:       "8080",
		proxy_port: "7331",
		dsn:        env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=db sslmode=disable"),
	}

	// connect to database
	conn, err := pgx.Connect(ctx, cfg.dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer conn.Close(ctx)

	logger.Debug("connected to database: %v", cfg.dsn)

	// init application variables
	app := application{
		db: conn,
	}

	// create server mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./internal/web/static/assets"))
	mux.Handle("/assets/", utils.DisableCacheInDevMode(http.StripPrefix("/assets/", fs)))

	// init services and handlers
	tenantService := service.NewTenantService(repo.New(app.db))
	tenantHandler := handler.NewTenantHandler(tenantService)

	propertyService := service.NewPropertyService(repo.New(app.db))
	propertyHandler := handler.NewPropertyHandler(propertyService)

	// page handlers
	mux.HandleFunc("/", tenantHandler.ListTenantPage)
	mux.HandleFunc("/property/create", propertyHandler.CreatePropertyPage)

	// init server
	s := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: mux,
	}

	// start server
	logger.Debug("For Production, open the actual port: http://localhost:%v", cfg.port)
	logger.Debug("For Development, open the proxy port (for templ hot reload): http://localhost:%v", cfg.proxy_port)

	log.Fatalln(s.ListenAndServe())
}
