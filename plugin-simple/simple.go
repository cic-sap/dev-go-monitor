package plugin_simple

import (
	"github.com/cic-sap/dev-go-monitor/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"net/http"
)

func Init(mux *http.ServeMux, options ...monitor.Option) http.Handler {
	conf := monitor.NewDefaultConf(options)

	m := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	mux.Handle(conf.Path, promhttp.Handler())
	h := std.Handler("", m, mux)

	monitor.Patch()
	return h
}
