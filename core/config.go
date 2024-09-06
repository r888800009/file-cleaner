package file_cleaner

import (
	"encoding/json"
	"errors"
	"os"
)

/*
Strategy defines the what to do with the file
*/
type Strategy interface {
	Load(name string, value map[string]interface{}) error
}

type DirEntry struct {
	path        string
	recursively bool
}

type StrategyConfig struct {
	name string
	// Strategy allow `source_to_target_dedupe` and `pdf_mover`
	strategy string
}

type SourceToTargetDedupeStrategy struct {
	super  StrategyConfig
	target DirEntry
	source []DirEntry
}

type Config struct {
	version    string
	strategies map[string]Strategy
}

// print dir entry
func (dir *DirEntry) Print() {
	println("Path:", dir.path, "Recursively:", dir.recursively)
}

// Load a strategy entry
func (config *SourceToTargetDedupeStrategy) Load(name string, value map[string]interface{}) error {
	config.super.name = name
	config.super.strategy = value["strategy"].(string)
	println("Strategy:", config.super.strategy)

	config.target.path = value["target_dir"].(string)
	config.target.recursively = value["target_dir_recursive"].(bool)
	config.target.Print()

	// Load source directories
	sourceDirs := value["source_dirs"].([]interface{})
	for _, sourceDir := range sourceDirs {
		dir := DirEntry{}
		dir.path = sourceDir.(map[string]interface{})["path"].(string)
		dir.recursively = sourceDir.(map[string]interface{})["recursive"].(bool)
		dir.Print()
		config.source = append(config.source, dir)
	}
	return nil
}

// Parse Strategy
func parseStrategy(key string, value map[string]interface{}) (Strategy, error) {
	strageKey, ok := value["strategy"].(string)
	if !ok {
		return nil, errors.New("strategy key not found")
	}

	switch strageKey {
	case "source_to_target_dedupe":
		strategy := new(SourceToTargetDedupeStrategy)
		strategy.Load(key, value)
		return strategy, nil
	case "pdf_mover":
		println("Loading pdf_mover strategy")
		return nil, errors.New("pdf_mover strategy not implemented")
		// strategy := new(PdfMoverStrategy)
		// strategy.Load(value.(map[string]interface{}))
		// config.strategies[key] = strategy
	default:
		return nil, errors.New("unknown strategy")
	}
}

// Load a configuration file
func (config_struct *Config) Load(path string) error {
	config, err := loadJson(path)
	if err != nil {
		return err
	}

	config_struct.version = config["version"].(string)
	println("Config Version:", config_struct.version)
	if config_struct.version != "0.1" {
		return errors.New("unsupported config version")
	}
	delete(config, "version")

	config_struct.strategies = make(map[string]Strategy)
	for key, jsonValue := range config {
		println("Found strategy entry:", key)
		value, ok := jsonValue.(map[string]interface{})
		if !ok {
			return errors.New("strategy parse error")
		}

		strategy, err := parseStrategy(key, value)
		if err != nil {
			return err
		}
		config_struct.strategies[key] = strategy
	}
	return nil
}

// Load a JSON file into a map
func loadJson(path string) (map[string]interface{}, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
