package monitor

import (
	req "github.com/imroc/req"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/**
trace http client
*/

var (
	defaultBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	//	prometheus.NewHistogramVec(prometheus.HistogramOpts{
	//Namespace: cfg.Prefix,
	//Subsystem: "http",
	//Name:      "request_duration_seconds",
	//Help:      "The latency of the HTTP requests.",
	//Buckets:   cfg.DurationBuckets,
	//}, []string{cfg.ServiceLabel, cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),
	//
	httpClientRequestDurHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "client_request_duration_seconds",
		Help:    "Number of http_client_request",
		Buckets: defaultBuckets,
	}, []string{"host", "path", "code", "method"})
)

type PrometheusTransport struct {
	originalTransport http.RoundTripper
}

func newTransport() *PrometheusTransport {
	return &PrometheusTransport{
		originalTransport: http.DefaultTransport,
	}
}

func (c *PrometheusTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	t1 := time.Now()

	labels := prometheus.Labels{
		"host":   r.URL.Host,
		"method": r.Method,
		"path":   r.URL.Path,
		"code":   "",
	}

	resp, err := c.originalTransport.RoundTrip(r)
	sec := time.Now().Sub(t1).Seconds()
	if err != nil {
		labels["code"] = "0"
		httpClientRequestDurHistogram.With(labels).Observe(sec)
		return nil, err
	}
	code := strconv.Itoa(resp.StatusCode)
	labels["code"] = code
	httpClientRequestDurHistogram.With(labels).Observe(sec)

	return resp, nil
}

var rt http.RoundTripper

var o sync.Once

func WarpTransport(org http.RoundTripper) http.RoundTripper {
	return &PrometheusTransport{
		originalTransport: org,
	}
}

func StartHttpClientTrace() {

	o.Do(func() {
		rt = newTransport()
		http.DefaultTransport = rt
		req.Client().Transport = rt
		_ = prometheus.Register(httpClientRequestDurHistogram)
	})

}
