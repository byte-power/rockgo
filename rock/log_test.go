package rock

import (
	"testing"

	"github.com/byte-power/rockgo/log"
	"github.com/byte-power/rockgo/util"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert.NotNil(t, Logger(""))
}

func TestFluentLog(t *testing.T) {
	tag := "test-fluent"
	config := make(util.AnyMap)
	config["level"] = "info"
	config["port"] = 24225

	_, err := parseFluentLogger(config, tag)
	assert.Error(t, err)

	config["host"] = "127.0.0.1"
	config["async"] = true
	logger, err := parseFluentLogger(config, tag)
	assert.NotNil(t, logger)
	assert.Nil(t, err)
}

func TestParseLogComponents(t *testing.T) {
	assert.Equal(t, log.LevelDebug, parseLevel(nil))
	assert.Equal(t, log.LevelDebug, parseLevel("Debug"))
	assert.Equal(t, log.LevelInfo, parseLevel("inFo"))
	assert.Equal(t, log.LevelWarn, parseLevel("Warn"))
	assert.Equal(t, log.LevelError, parseLevel("erroR"))
	assert.Equal(t, log.LevelFatal, parseLevel("fatAl"))
	assert.Equal(t, log.MessageFormatJSON, parseMessageFormat("json"))
	assert.Equal(t, log.MessageFormatText, parseMessageFormat("text"))
	assert.Equal(t, log.TimeFormatISO8601, log.MakeTimeFormat("iso8601"))
	assert.Equal(t, log.TimeFormatSeconds, log.MakeTimeFormat("seconds"))
	assert.Equal(t, log.TimeFormatMillis, log.MakeTimeFormat("millis"))
	assert.Equal(t, log.TimeFormatNanos, log.MakeTimeFormat("nanos"))
}
