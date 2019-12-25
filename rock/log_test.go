package rock

import (
	"github.com/byte-power/rockgo/util"
	"testing"
)

func TestFluentLog(t *testing.T) {
	config := make(util.AnyMap)
	config["level"] = "info"
	config["tag"] = "test-fluent"
	config["port"] = 24225

	_, err := parseFluentLogger(config)
	if err == nil {
		t.Error("Should error if host not in config.")
	}

	config["host"] = "127.0.0.1"
	config["async"] = true
	logger, _ := parseFluentLogger(config)

	if logger == nil {
		t.Error("Create FluentLogger failed.")
	}
}
