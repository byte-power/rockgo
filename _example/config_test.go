package main

import (
	"testing"
	"github.com/byte-power/rockgo/rock"
	"github.com/byte-power/rockgo/util"
)

func TestConfigImporting(t *testing.T) {
	if err := rock.ImportConfigFilesFromDirectory("settings"); err != nil {
		t.Error(err)
	}
	if v := rock.ConfigIn("jd.a"); util.AnyToInt64(v) != 1 {
		t.Errorf("jd.a should be 1: %v", v)
	}
	if v := rock.ConfigIn("yd.z"); util.AnyToString(v) != "io" {
		t.Errorf("ya.z should be 'io': %v", v)
	}
	if v, ok := rock.ConfigIn("yd.a").([]interface{}); !ok || len(v) != 1 || v[0] != 10 {
		t.Errorf("yd.a should be [10]: %v", v)
	}
	if v := rock.ConfigIn("yd.d.a"); util.AnyToInt64(v) != 1 {
		t.Errorf("yd.d.a should be 1: %v", v)
	}
}
