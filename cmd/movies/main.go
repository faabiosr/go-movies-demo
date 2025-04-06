package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
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

const dbPerm = 0o600

const timeout = 10 * time.Second

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
	db, err := bolt.Open(dbPath, dbPerm, nil)
	if err != nil {
		e.Logger.Fatal(err)
	}

	ds := movies.NewDatasource(db)

	// API Resources
	movies.Routes(e.Group("/movies"), ds)

	// Start server
	e.Logger.Infof("%s service", appName)

	go func() {
		if err := start(e, appAddr); err != nil {
			e.Logger.Info("shutting down the service")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func start(e *echo.Echo, host string) error {
	listeners, err := activation.Listeners()
	if err != nil {
		return nil
	}

	if len(listeners) > 0 {
		e.Listener = listeners[0]
	}

	return e.Start(host)
}
