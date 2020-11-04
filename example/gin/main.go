package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/cic-sap/dev-go-monitor/plugin-gin"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	r := gin.Default()

	plugin_gin.Init(r, monitor.WithPath("/metrics"),
		monitor.WithGinHandle(func(c *gin.Context) string {
			if c.FullPath() == "/hi/:id" {
				return "/hi/" + strings.ToLower(c.Param("id"))
			}
			return c.FullPath()
		}))

	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * time.Duration(rand.Int31n(3)))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/500", func(c *gin.Context) {

		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%d", 2/rand.Intn(2)))
	})
	r.GET("/hi/:id", func(c *gin.Context) {

		log.Println("path:", c.FullPath())

		id := c.Param("id")
		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%s", id))
	})
	r.GET("/info/*any", func(c *gin.Context) {

		log.Println("path:", c.FullPath())

		c.String(http.StatusOK,
			fmt.Sprintf("hello info:%s", c.FullPath()))
	})
	r.Run(":8094")
}
