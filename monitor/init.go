package monitor

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

type Conf struct {
	Path            string
	GinCustomHandle func(c *gin.Context) string
	Gatherers       prometheus.Gatherers
}

func NewDefaultConf(options []Option) *Conf {
	conf := &Conf{
		Path: "/metrics",
	}
	for _, op := range options {
		op(conf)
	}
	return conf
}

type Option func(conf *Conf)

func WithPath(path string) Option {
	return func(conf *Conf) {
		conf.Path = path
	}
}

func WithGinHandle(fun func(c *gin.Context) string) Option {
	return func(conf *Conf) {
		conf.GinCustomHandle = fun
	}
}

func WithGatherer(g prometheus.Gatherer) Option {
	return func(conf *Conf) {
		conf.Gatherers = append(conf.Gatherers, g)
	}
}

func (conf Conf) BuildHandler() http.Handler {
	if len(conf.Gatherers) == 0 {
		return promhttp.Handler()
	}
	return promhttp.HandlerFor(conf.Gatherers, promhttp.HandlerOpts{})
}

var one sync.Once

func Patch() {

	one.Do(func() {
		InitUptime()
		StartHttpClientTrace()
	})

}
