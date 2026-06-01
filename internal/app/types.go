package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// think of config as constants, variables inside do not change over time, and is initialized only before starting the application
type Config struct {
	Port      string // our main port (or address) for prod
	ProxyPort string // the port will be used for development with templ
	Dsn       string // Data Source Name, or database connection string
}

// think of the application as states, variables inside change over time
type Application struct {
	Db *pgxpool.Pool // db driver
}
