package rock

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/byte-power/rockgo/util"
	"gopkg.in/yaml.v2"
)

// LoadConfigFromFile support json and yaml format with extension (json/yaml/yml)
func LoadConfigFromFile(path string, into interface{}) (err error) {
	ext := strings.ToLower(filepath.Ext(path))
	var content []byte
	switch ext {
	case ".json":
		content, err = readNonEmptyFile(path)
		if err == nil {
			err = json.Unmarshal(content, into)
		}
	case ".yaml", ".yml":
		content, err = readNonEmptyFile(path)
		if err == nil {
			err = yaml.Unmarshal(content, into)
		}
	default:
		err = errors.New(ErrNameUnsupportedFormat)
	}
	return
}

func readNonEmptyFile(path string) ([]byte, error) {
	content, err := util.ReadFileFromPath(path)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("empty content in " + path)
	}
	return content, nil
}

var sharedConfig sync.Map

// ImportConfigFromPaths load multiple config files.
// And then store in a shared map with key by filename (stripped ext), e.g. abc for abc.json.
// - Parameters:
//   - paths: Path list to be loaded, empty string would be ignored.
// - Return: each got error
func ImportConfigFromPaths(paths ...string) error {
	for _, path := range paths {
		if path == "" {
			continue
		}
		basename := filepath.Base(path)
		if basename == "." || basename == ".." {
			continue
		}
		pos := strings.Index(basename, ".")
		if pos <= 0 {
			defaultLogger.Warnf("skipped config file '%s', it should named with ext", basename)
			continue
		}
		var content interface{}
		if err := LoadConfigFromFile(path, &content); err != nil {
			return err
		}
		sharedConfig.Store(basename[:pos], content)
	}
	return nil
}

// Load each file in the <dir> without recursive by ImportConfigFromPaths().
// Contents of file would be stored in shared config with filename (without extension).
// e.g. app.yaml would use "app" as key.
// - Return: each got error
func ImportConfigFilesFromDirectory(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		err = ImportConfigFromPaths(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// Find value from shared configs with <keyPath> delimited by ".".
func ConfigIn(keyPath string) interface{} {
	if strings.Contains(keyPath, ".") {
		comps := strings.Split(keyPath, ".")
		keys := make([]interface{}, len(comps))
		for i, v := range comps {
			keys[i] = v
		}
		return util.FindInSyncMapWithKeys(&sharedConfig, keys)
	}
	v, _ := sharedConfig.Load(keyPath)
	return v
}

// Get shared configs
func Config() *sync.Map {
	return &sharedConfig
}
