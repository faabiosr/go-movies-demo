package main

import (
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	glog "github.com/labstack/gommon/log"
)

const (
	appAddr = "0.0.0.0:8000"
	appName = "moviez"
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

	// Start server
	e.Logger.Infof("%s service", appName)

	if err := e.Start(appAddr); err != nil {
		e.Logger.Fatal(err)
	}
}
