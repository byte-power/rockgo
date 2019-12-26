package rock

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"path/filepath"
	"strings"

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

var sharedConfig = util.AnyMap{}

// 载入多个配置文件，遇到任何错误将立即返回
//
// Load multiple config files
// and then store in a shared map with key by filename (stripped ext), e.g. abc for abc.json.
//
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
		sharedConfig[basename[:pos]] = content
	}
	return nil
}
// 载入指定目录下的所有文件（不递归）
// 文件将按照文件名作为map中的key，如app.yaml的内容将保存至data["app"]
//
// Load each file in the [dir] without recursive by ImportConfigFromPaths()
//
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
		err = ImportConfigFromPaths(dir + "/" + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// 按以.分隔的keyPath查找配置
//
// Find value in map key path from shared configs
func ConfigIn(keyPath string) interface{} {
	return util.FindInAnyMapWithKeys(sharedConfig, strings.Split(keyPath, "."))
}

// 载入的所有配置
//
// Get shared configs
func Config() map[string]interface{} {
	return sharedConfig
}
