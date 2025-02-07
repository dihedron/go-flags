package flags

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

var (
	// ConfigSupportedExts is a list of supported config file extensions.
	ConfigSupportedExts = []string{".json", ".yaml", ".yml"}

	configPaths []string
	configFile  string
	configName  string

	configMap map[string]any
)

// AddConfigPath adds a path for go-falgs to search for the config file in.
// Can be called multiple times to define multiple search paths; paths will
// be searched in the order they are added.
func AddConfigPath(path string) {
	var err error
	path = os.ExpandEnv(path)
	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			path = ""
		}
	}
	path = filepath.Clean(path)

	if path != "" {
		slog.Info("adding path to search paths", "path", path)
		if !slices.Contains(configPaths, path) {
			configPaths = append(configPaths, path)
		}
	}
}

// SetConfigFile explicitly defines the path, name and extension of the config file.
// Viper will use this and not check any of the config paths.
func SetConfigFile(in string) {
	if in != "" {
		configFile = in
	}
}

// SetConfigName sets name for the config file.
// Does not include extension.
func SetConfigName(in string) {
	if in != "" {
		configName = in
		configFile = ""
	}
}

// ReadInConfig will discover and load the configuration file from disk
// and key/value stores, searching in one of the defined paths.
func ReadInConfig() error {
	slog.Info("attempting to read in config file")

	if configFile == "" {
	loop:
		for _, path := range configPaths {
			for _, ext := range ConfigSupportedExts {
				path := filepath.Join(path, configName+ext)
				slog.Debug("testing file", "path", path)
				if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
					slog.Debug("found file", "path", path)
					configFile = path
					break loop
				}
			}
		}
	}

	if !slices.Contains(ConfigSupportedExts, filepath.Ext(configFile)) {
		return errors.New("unsupported config file extension")
	}

	slog.Debug("reading file", "file", configFile)
	data, err := os.ReadFile(configFile)
	if err != nil {
		slog.Error("error reading config file", "error", err)
		return err
	}

	configMap = make(map[string]any)

	switch filepath.Ext(configFile) {
	case ".json":
		if err = json.Unmarshal(data, &configMap); err != nil {
			slog.Error("error unmarshalling JSON", "error", err)
			return err
		}
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(data, &configMap); err != nil {
			slog.Error("error unmarshalling YAML", "error", err)
			return err
		}
	}
	return nil
}
