package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/app"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/env"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/handler"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/middleware"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/routine"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func startServer(s *http.Server, serverErr chan error) {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		serverErr <- err
	}
}

func main() {
	// global context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load env variables
	godotenv.Load()

	// init app config
	cfg := app.Config{
		Port:      "8080",
		ProxyPort: "7331",
		Dsn:       env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=db sslmode=disable"),
	}

	// connect to database
	pool, err := pgxpool.New(ctx, cfg.Dsn)
	if err != nil {
		log.Fatalf("failed to create a pool and connect to DB: %v", err)
	}
	defer pool.Close()
	logger.Debug("connected to database: %v", cfg.Dsn)

	// check if db container is running
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("[EXIT] failed to reach database (is the docker container for the db running?): %v", err)
	}

	// init application variables
	application := app.Application{
		Db: pool,
	}

	// init db queries
	queries := repo.New(application.Db)

	// create server mux
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./internal/web/static/assets"))
	mux.Handle("/assets/", utils.DisableCacheInDevMode(http.StripPrefix("/assets/", fs)))

	// init services and handlers
	tenantService := service.NewTenantService(queries)
	tenantHandler := handler.NewTenantHandler(tenantService)

	propertyService := service.NewPropertyService(queries)
	propertyHandler := handler.NewPropertyHandler(propertyService)

	authService := service.NewAuthService(application.Db)
	authHandler := handler.NewAuthHandler(authService)

	pageHandler := handler.NewPageHandler(handler.PageHandler{
		PropertyService: propertyService,
		TenantService:   tenantService,
	})

	// TODO: refactor this later, find a more elegant way to write this to handler multiple middlewares
	// auth middleware func
	protect := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(http.HandlerFunc(h))
	}

	// TODO: find a better way to write, .Handle() and HandleFunc(), they both serve the same purpose, only difference is that .Handle() needs you to convert the handler functions you pass
	// page handlers
	mux.HandleFunc("GET /login", pageHandler.LoginPage)
	mux.HandleFunc("GET /register", pageHandler.RegisterPage)
	mux.Handle("GET /dashboard", protect(pageHandler.DashboardPage))
	mux.Handle("GET /property/create", protect(pageHandler.CreatePropertyPage))
	mux.Handle("GET /tenant/create", protect(pageHandler.CreateTenantPage))

	// handlers
	mux.Handle("POST /property/create", protect(propertyHandler.CreateProperty))
	mux.Handle("POST /tenant/create", protect(tenantHandler.CreateTenant))
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	// init server
	s := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: middleware.Chain(mux,
			middleware.AttachUser(authService),
			// middleware.RequestLogger, // tmp: this is just an example of adding more middlewares later
		),
	}

	// init routines
	go routine.DeleteExpiredSessions(ctx, queries)

	// init server error channel
	serverErr := make(chan error, 1)
	defer close(serverErr)

	// start server asynchronously
	go startServer(s, serverErr)
	logger.Debug("For Production, open the actual port: http://localhost:%v", cfg.Port)
	logger.Debug("For Development, open the proxy port (for templ hot reload): http://localhost:%v", cfg.ProxyPort)

	// block here until either: (a) a signal fires ctx.Done, or (b) the server crashes
	select {
	case <-ctx.Done():
		logger.Debug("shutdown signal received")
	case err := <-serverErr:
		logger.Error("server failed: ", err)
	}

	// give in-flight requests up to 30s to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error: ", err)
	}

	logger.Debug("server stopped cleanly")
	// defers fire as main returns: stop(), pool.Close()
}
