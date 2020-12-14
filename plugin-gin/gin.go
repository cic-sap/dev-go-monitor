package plugin_gin

import (
	"context"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
)

func Init(r *gin.Engine, options ...monitor.Option) {

	conf := monitor.NewDefaultConf(options)

	m := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	r.Use(Handler(conf, m))
	r.GET(conf.Path, gin.WrapH(promhttp.Handler()))

	monitor.Patch()
}

func Handler(conf *monitor.Conf, m middleware.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := &reporter{c: c}
		h := c.FullPath()
		if conf != nil && conf.GinCustomHandle != nil {
			h = conf.GinCustomHandle(c)
		}
		m.Measure(h, r, func() {
			c.Next()
		})
	}
}

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string { return r.c.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }
