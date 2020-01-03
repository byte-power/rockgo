package rock

import (
	"testing"

	"github.com/byte-power/rockgo/util"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	assert.Error(t, ImportConfigFilesFromDirectory(""))
	assert.Error(t, ImportConfigFilesFromDirectory("app.go"))

	pathSettings := "../_example/settings"
	assert.Nil(t, ImportConfigFilesFromDirectory(pathSettings+"/rockgo.yaml"))
	assert.Nil(t, ImportConfigFilesFromDirectory(pathSettings))
	
	assert.Equal(t, int64(1), util.AnyToInt64(ConfigIn("jd.a")))
	assert.Equal(t, "io", util.AnyToString(ConfigIn("yd.z")))
	assert.Equal(t, []interface{}{10}, ConfigIn("yd.a"))
	assert.Equal(t, int64(1), util.AnyToInt64(ConfigIn("yd.d.a")))
	assert.IsType(t, []interface{}{}, ConfigIn("ja"))
}
