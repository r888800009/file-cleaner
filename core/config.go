package file_cleaner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DirEntry struct {
	path         string
	recursively  bool
	ignore_regex *regexp.Regexp
	match_regex  *regexp.Regexp
}

type StrategyConfig struct {
	name string
	// Strategy allow `source_to_target_dedupe` and `pdf_mover`
	strategy string
}

type SourceToTargetDedupeStrategy struct {
	super     StrategyConfig
	target    DirEntry
	source    []DirEntry
	trashPath string
}

type Config struct {
	version    string
	strategies map[string]Strategy
}

// print dir entry
func (dir *DirEntry) Print() {
	fmt.Println("Path:", dir.path, "Recursively:", dir.recursively)
}

// Load a strategy entry
func (config *SourceToTargetDedupeStrategy) Load(name string, value map[string]interface{}) error {
	config.super.name = name
	config.super.strategy = value["strategy"].(string)
	fmt.Println("Strategy:", config.super.strategy)

	config.target.Load(value["target_dir"].(map[string]interface{}))
	config.target.Print()

	config.trashPath = value["trash_dir"].(string)
	config.trashPath = expandDir(config.trashPath)
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02-15-04-05.000")
	config.trashPath = filepath.Join(config.trashPath, formattedTime)
	fmt.Println("Trash Path:", config.trashPath)

	// Load source directories
	sourceDirs := value["source_dirs"].([]interface{})
	for _, sourceDir := range sourceDirs {
		dir := DirEntry{}
		dir.Load(sourceDir.(map[string]interface{}))
		dir.Print()
		config.source = append(config.source, dir)
	}
	return nil
}

/*
this function expands the directory path
supports ~ and ~/ expansion
*/
func expandDir(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		return dir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}
	return path
}

// Load DirEntry
func (dirEntry *DirEntry) Load(value map[string]interface{}) {
	dirEntry.path = value["path"].(string)
	dirEntry.path = expandDir(dirEntry.path)

	dirEntry.recursively = value["recursive"].(bool)

	// load ignore regex if it exists
	if ignore, ok := value["ignore"]; ok {
		fmt.Println("Ignore regex:", ignore)
		dirEntry.ignore_regex = regexp.MustCompile(ignore.(string))
	} else {
		fmt.Println("No ignore regex")
		dirEntry.ignore_regex = nil
	}

	// load match regex if it exists
	if match, ok := value["match"]; ok {
		fmt.Println("Match regex:", match)
		dirEntry.match_regex = regexp.MustCompile(match.(string))
	} else {
		fmt.Println("No match regex")
		dirEntry.match_regex = nil
	}
}

func (dirEntry *DirEntry) Match(path string) bool {
	result := true
	if dirEntry.match_regex != nil {
		result = result && dirEntry.match_regex.MatchString(path)
	}
	if dirEntry.ignore_regex != nil {
		result = result && !dirEntry.ignore_regex.MatchString(path)
	}
	return result
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
		fmt.Println("Loading pdf_mover strategy")
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
	fmt.Println("Config Version:", config_struct.version)
	if config_struct.version != "0.1" {
		return errors.New("unsupported config version")
	}
	delete(config, "version")

	config_struct.strategies = make(map[string]Strategy)
	for key, jsonValue := range config {
		fmt.Println("Found strategy entry:", key)
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
