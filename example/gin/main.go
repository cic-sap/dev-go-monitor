package main

import (
	"fmt"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/cic-sap/dev-go-monitor/plugin-gin"
	"github.com/gin-gonic/gin"
	req "github.com/imroc/req"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var reg = prometheus.NewRegistry()

var c1 = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "c1",
})

func main() {
	reg.MustRegister(c1)

	r := gin.Default()

	plugin_gin.Init(r, monitor.WithPath("/metrics"),
		monitor.WithGinHandle(func(c *gin.Context) string {
			if c.FullPath() == "/hi/:id" {
				return "/hi/" + strings.ToLower(c.Param("id"))
			}
			return c.FullPath()
		}),
		monitor.WithGatherer(reg),
		monitor.WithGatherer(prometheus.DefaultGatherer),
	)

	r.GET("/", func(c *gin.Context) {
		c1.Add(1.0)
		time.Sleep(time.Second * time.Duration(rand.Int31n(3)))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/500", func(c *gin.Context) {

		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%d", 2/rand.Intn(2)))
	})
	r.GET("/hi/:id", func(c *gin.Context) {

		//log.Println("path:", c.FullPath())

		id := c.Param("id")
		c.String(http.StatusOK,
			fmt.Sprintf("hello world:%s", id))
	})
	r.GET("/info/*any", func(c *gin.Context) {

		//log.Println("path:", c.FullPath(), c.Param("any"))
		req.Get("https://httpbin.org/" + c.Param("any"))
		c.String(http.StatusOK,
			fmt.Sprintf("hello info:%s", c.FullPath()))
	})
	r.Run(":8094")
}
