package plugin_echo

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"log"
)

func Init(r *echo.Echo, options ...monitor.Option) {

	conf := monitor.NewDefaultConf(options)

	m := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	r.Use(Handler(m))
	r.GET(conf.Path, echo.WrapHandler(promhttp.Handler()))

}

// Handler returns a Echo measuring middleware.
func Handler(m middleware.Middleware) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {

		return echo.HandlerFunc(func(c echo.Context) error {
			log.Println("get path:", c.Path())
			r := &reporter{c: c}
			var err error
			m.Measure(c.Path(), r, func() {
				err = h(c)
			})
			return err
		})
	}
}

type reporter struct {
	c echo.Context
}

func (r *reporter) Method() string { return r.c.Request().Method }

func (r *reporter) Context() context.Context { return r.c.Request().Context() }

func (r *reporter) URLPath() string { return r.c.Request().URL.Path }

func (r *reporter) StatusCode() int { return r.c.Response().Status }

func (r *reporter) BytesWritten() int64 { return r.c.Response().Size }
