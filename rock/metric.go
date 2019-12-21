package rock

import (
	"time"

	"github.com/byte-power/rockgo/util"
	"gopkg.in/alexcesaro/statsd.v2"
)

var managedMetricInstance *statsd.Client
var metricPrefix string

func initMetric(m util.AnyMap) {
	// TODO: init metric
}

func Metric() *statsd.Client {
	return managedMetricInstance
}

func MetricIncrease(key string) {
	if managedMetricInstance == nil {
		return
	}
	managedMetricInstance.Increment(metricPrefix + "." + key)
}

// duration: time.Now().Sub(oldTime)
func MetricTiming(key string, duration time.Duration) {
	if managedMetricInstance == nil {
		return
	}
	managedMetricInstance.Timing(metricPrefix+"."+key, int(duration/time.Millisecond))
}
