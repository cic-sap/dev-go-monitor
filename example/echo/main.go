package main

import (
	"github.com/labstack/echo/v4"
	plugin_echo "github.com/cic-sap/dev-go-monitor/plugin-echo"
	"log"
	"net/http"
)

const (
	srvAddr = "127.0.0.1:8480"
)

func main() {
	// Create our middleware factory with the default settings.

	// Create Echo instance and global middleware.
	e := echo.New()

	// Add our handler.
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.GET("/users/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.String(http.StatusOK, name)
	})

	e.GET("/500", func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "ERR 500")
	})
	plugin_echo.Init(e)

	log.Printf("server listening at %s", srvAddr)
	if err := http.ListenAndServe(srvAddr, e); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
