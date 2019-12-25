package rock

import (
	"time"
	"errors"

	"github.com/byte-power/rockgo/util"
	"gopkg.in/alexcesaro/statsd.v2"
)

var managedMetricInstance *statsd.Client
var metricPrefix string

func initMetric(m util.AnyMap) (err error) {
	var opts []statsd.Option
	if it, ok := m["host"]; ok {
		if host, ok := it.(string); ok {
			opts = append(opts, statsd.Address(host))
		} else {
			err = errors.New("metric.host should be string")
			return
		}
	}
	if it, ok := m["prefix"].(string); ok {
		metricPrefix = it
		opts = append(opts, statsd.Prefix(it))
	}
	if it := util.AnyToInt64(m["max_packet_size"]); it > 0 {
		opts = append(opts, statsd.MaxPacketSize(int(it)))
	}
	if it, ok := m["flush_period_seconds"]; ok {
		opts = append(opts, statsd.FlushPeriod(time.Second*time.Duration(util.AnyToInt64(it))))
	}
	if it, ok := m["network"].(string); ok {
		opts = append(opts, statsd.Network(it))
	}
	if it, ok := m["mute"]; ok {
		opts = append(opts, statsd.Mute(util.AnyToBool(it)))
	}
	if it, ok := m["sample_rate"]; ok {
		opts = append(opts, statsd.SampleRate(float32(util.AnyToFloat64(it))))
	}
	if it := util.AnyToAnyMap(m["tags"]); it != nil {
		tags := make([]string, 0, len(it)*2)
		for k, v := range it {
			tags = append(tags, k, util.AnyToString(v))
		}
		opts = append(opts, statsd.Tags(tags...))
	}
	managedMetricInstance, err = statsd.New(opts...)
	return
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

func MetricGauge(bucket string, value interface{}) {
	if managedMetricInstance == nil {
		return
	}
	managedMetricInstance.Gauge(bucket, value)
}

func MetricHistogram(bucket string, value interface{}) {
	if managedMetricInstance == nil {
		return
	}
	managedMetricInstance.Histogram(bucket, value)
}
