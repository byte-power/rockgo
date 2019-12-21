package rock

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/byte-power/rockgo/util"
	"gopkg.in/yaml.v2"
)

// LoadConfigFromFile support json and yaml format with extension (json/yaml/yml)
func LoadConfigFromFile(path string, into interface{}) error {
	content, err := util.ReadFileFromPath(path)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return errors.New("empty content in " + path)
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return json.Unmarshal(content, into)
	case ".yaml", ".yml":
		return yaml.Unmarshal(content, into)
	default:
		return errors.New("unsupported format")
	}
}
