package plugin_rest

import (
	"context"
	"github.com/emicklei/go-restful"
	gorestful "github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics2 "github.com/slok/go-http-metrics/metrics"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/cic-sap/dev-go-monitor/monitor"
	"sync"
)

var record metrics2.Recorder
var once sync.Once

func Init(c *restful.Container, options ...monitor.Option) {

	conf := monitor.NewDefaultConf(options)

	once.Do(func() {
		record = metrics.NewRecorder(metrics.Config{})
	})
	m := middleware.New(middleware.Config{
		Recorder: record,
	})
	c.Filter(Handler(m))
	c.Handle(conf.Path, promhttp.Handler())

}

// Handler returns a gorestful measuring middleware.
func Handler(m middleware.Middleware) gorestful.FilterFunction {
	return func(req *gorestful.Request, resp *gorestful.Response, chain *gorestful.FilterChain) {
		r := &reporter{req: req, resp: resp}
		m.Measure(req.SelectedRoutePath(), r, func() {
			chain.ProcessFilter(req, resp)
		})
	}
}

type reporter struct {
	req  *gorestful.Request
	resp *gorestful.Response
}

func (r *reporter) Method() string { return r.req.Request.Method }

func (r *reporter) Context() context.Context { return r.req.Request.Context() }

func (r *reporter) URLPath() string { return r.req.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.resp.StatusCode() }

func (r *reporter) BytesWritten() int64 { return int64(r.resp.ContentLength()) }
