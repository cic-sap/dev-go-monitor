package monitor

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Conf struct {
	Path            string
	GinCustomHandle func(c *gin.Context) string
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

var one sync.Once

func Patch() {

	one.Do(func() {
		InitUptime()
		StartHttpClientTrace()
	})

}
