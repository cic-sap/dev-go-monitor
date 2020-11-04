package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var startTime = time.Now()

func InitUptime() {

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "process_uptime_seconds",
		Help: "process uptime seconds",
	}, func() float64 {
		return time.Since(startTime).Seconds()
	})

}
