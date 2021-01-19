package rock

import (
	"errors"
	"time"

	"github.com/byte-power/rockgo/util"
	"gopkg.in/alexcesaro/statsd.v2"
)

var managedMetricInstance *statsd.Client
var managedMetricPrefix string

func initMetric(prefix string, cfg util.AnyMap) (err error) {
	var opts []statsd.Option
	if it, ok := cfg["host"]; ok {
		if host, ok := it.(string); ok {
			opts = append(opts, statsd.Address(host))
		} else {
			err = errors.New("metric.host should be string")
			return
		}
	}
	managedMetricPrefix = prefix
	opts = append(opts, statsd.Prefix(prefix))
	if it := util.AnyToInt64(cfg["max_packet_size"]); it > 0 {
		opts = append(opts, statsd.MaxPacketSize(int(it)))
	}
	if it, ok := cfg["flush_period_seconds"]; ok {
		opts = append(opts, statsd.FlushPeriod(time.Second*time.Duration(util.AnyToInt64(it))))
	}
	if it, ok := cfg["network"].(string); ok {
		opts = append(opts, statsd.Network(it))
	}
	if it, ok := cfg["mute"]; ok {
		opts = append(opts, statsd.Mute(util.AnyToBool(it)))
	}
	if it, ok := cfg["sample_rate"]; ok {
		opts = append(opts, statsd.SampleRate(float32(util.AnyToFloat64(it))))
	}
	if it := util.AnyToStrMap(cfg["tags"]); it != nil {
		tags := make([]string, 0, len(it)*2)
		for k, v := range it {
			tags = append(tags, k, util.AnyToString(v))
		}
		opts = append(opts, statsd.Tags(tags...))
	}
	managedMetricInstance, err = statsd.New(opts...)
	return
}

// Metric pass statsd client to make custom record.
//
// - Return: may be nil if not calling rock.Application#InitWithConfig() or not configure correctly.
func Metric() *statsd.Client {
	return managedMetricInstance
}

// MetricCount would change count on <num> for key.
func MetricCount(key string, num interface{}) {
	if managedMetricInstance != nil {
		// println("MetricCount ", managedMetricPrefix, key, num)
		managedMetricInstance.Count(key, num)
	}
}

// MetricIncrease would increase count on 1 for key with statsd count.
func MetricIncrease(key string) {
	MetricCount(key, 1)
}

// MetricDecrease would decrease count on 1 for key with statsd count.
func MetricDecrease(key string) {
	MetricCount(key, -1)
}

// MetricTimeDuration would record time duration for key with statsd timing.
//
// - Parameters:
//   - duration: e.g. time.Now().Sub(oldTime) or time.Second * 4
func MetricTimeDuration(key string, duration time.Duration) {
	if managedMetricInstance != nil {
		sec := float64(duration) / float64(time.Millisecond)
		// println("MetricTiming", managedMetricPrefix, key, sec)
		managedMetricInstance.Timing(key, sec)
	}
}

func MetricTiming(key string, value interface{})  {
	if managedMetricInstance != nil {
		managedMetricInstance.Timing(key, value)
	}	
}

func MetricGauge(bucket string, value interface{}) {
	if managedMetricInstance != nil {
		managedMetricInstance.Gauge(bucket, value)
	}
}

func MetricHistogram(bucket string, value interface{}) {
	if managedMetricInstance != nil {
		managedMetricInstance.Histogram(bucket, value)
	}
}
