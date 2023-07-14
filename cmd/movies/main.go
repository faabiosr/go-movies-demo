package main

import (
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	glog "github.com/labstack/gommon/log"
	bolt "go.etcd.io/bbolt"

	"github.com/faabiosr/go-movies-demo/internal/movies"
)

const (
	appAddr   = "0.0.0.0:8000"
	appName   = "moviez"
	dbName    = "catalog.db"
	dbPathEnv = "MOVIES_DB_PATH"
)

func main() {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Logger.SetLevel(glog.INFO)

	// Middlewares
	e.Use(mw.Recover())
	e.Pre(mw.RemoveTrailingSlash())
	e.Use(mw.RequestID())
	e.Use(mw.Secure())
	e.Use(mw.Logger())

	if os.Getenv(dbPathEnv) == "" {
		e.Logger.Fatalf("env '%s' was not defined", dbPathEnv)
	}

	dbPath := filepath.Join(os.Getenv(dbPathEnv), dbName)

	// Database connect
	db, err := bolt.Open(dbPath, 0o600, nil) // nolint:gomnd
	if err != nil {
		e.Logger.Fatal(err)
	}

	ds := movies.NewDatasource(db)

	// API Resources
	movies.Routes(e.Group("/movies"), ds)

	// Start server
	e.Logger.Infof("%s service", appName)

	if err := e.Start(appAddr); err != nil {
		e.Logger.Fatal(err)
	}
}
