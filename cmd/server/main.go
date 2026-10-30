package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/app"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/database/migrations"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/env"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/handler"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/middleware"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/routine"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/slogfmt"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/static/assets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func startServer(s *http.Server, serverErr chan error) {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		serverErr <- err
	}
}

func setupAssetsRoutes(mux *http.ServeMux) {
	development := os.Getenv("SHADCN_TEMPL_DEV") == "true" || os.Getenv("TEMPL_DEV_MODE") == "true"
	files := http.FileServer(http.FS(assets.Assets))
	if development {
		files = http.FileServer(http.Dir("./internal/web/static/assets"))
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if development {
			w.Header().Set("Cache-Control", "no-store")
		} else if strings.HasPrefix(strings.TrimPrefix(r.URL.Path, "/"), "js/shadcn-templ-") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", handler))
}

func main() {
	// global context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load env variables
	godotenv.Load()
	databaseUrl := env.GetString("DATABASE_URL", "")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL env variable empty, check .env.example\n")
	}

	port := env.GetString("PORT", "8080")

	// app arguments
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		err := migrations.Up(databaseUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Migration fail bohoo: %v\n", err)
			os.Exit(1)
			return
		}

		fmt.Println("Migration success wahoo!\n")
		os.Exit(0)
		return
	}

	// TODO: we now use slog for logging, refactor all instances of using the old "github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	// or update the go-toolkit to use the new custom logger, check logfmt/logfmt.go
	slogfmt.Init()

	// init app config
	cfg := app.Config{
		Port:      port,
		ProxyPort: "7331",
		Dsn:       databaseUrl, // TODO: test this change later
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

	// create main server mux
	mux := http.NewServeMux()

	// serve static files
	setupAssetsRoutes(mux)

	// subFS, err := fs.Sub(assets.Assets, ".")
	// if err != nil {
	// 	log.Fatalf("Failed to located web assets FS: %v", err)
	// }
	// assetsFS := http.FS(subFS)
	// mux.Handle("/assets/", utils.DisableCacheInDevMode(http.StripPrefix("/assets/", http.FileServer(assetsFS))))

	// init services
	authService := service.NewAuthService(application.Db)
	tenantService := service.NewTenantService(queries)
	leaseService := service.NewLeaseService(queries)
	propertyService := service.NewPropertyService(queries)
	transactionService := service.NewTransactionService(queries)
	registrationService := service.NewRegistrationService(application.Db, queries)

	// init handlers
	authHandler := handler.NewAuthHandler(authService)
	tenantHandler := handler.NewTenantHandler(tenantService, registrationService)
	propertyHandler := handler.NewPropertyHandler(propertyService)
	tradeHandler := handler.NewTradeHandler(transactionService)
	pageHandler := handler.NewPageHandler(propertyService, leaseService)

	// TODO: refactor this later, find a more elegant way to write this to handler multiple middlewares
	// auth middleware func
	protect := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(http.HandlerFunc(h))
	}

	// auth handlers
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	// TODO: find a better way to write, .Handle() and HandleFunc(), they both serve the same purpose, only difference is that .Handle() needs you to convert the handler functions you pass
	// page handlers
	mux.HandleFunc("GET /login", pageHandler.LoginPage)
	mux.HandleFunc("GET /register", pageHandler.RegisterPage)
	mux.Handle("GET /dashboard", protect(pageHandler.DashboardPage))
	mux.Handle("GET /property/create", protect(pageHandler.CreatePropertyPage))
	mux.Handle("GET /tenant/create", protect(pageHandler.CreateTenantPage))
	mux.Handle("GET /trade/create", protect(pageHandler.CreateTradePage))

	// CRUD handlers
	mux.Handle("POST /property/create", protect(propertyHandler.CreateProperty))
	mux.Handle("POST /tenant/create", protect(tenantHandler.CreateTenant))
	mux.Handle("DELETE /tenant/delete", protect(tenantHandler.DeleteTenant))
	mux.Handle("POST /trade/create", protect(tradeHandler.CreateTrade)) // TODO

	// Datastar handlers
	mux.Handle("POST /ds/tenant/create/compute_next_due_date", protect(tenantHandler.ComputeNextDueDate))

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
