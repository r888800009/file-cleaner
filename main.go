package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gofrs/flock"

	file_cleaner "github.com/r888800009/file_cleaner/core"
)

func parseArgs() (config *file_cleaner.Config, cmdArgs *file_cleaner.CmdLineArgs, err error) {
	// Define flags
	var configPath = flag.String("config", "", "Path to the configuration file")
	var dryRun = flag.Bool("dry-run", true, "Run the program in dry-run mode (no changes will be made)")
	var replaceAsSymlink = flag.Bool("replace-as-symlink", false, "Replace duplicate files with symlinks and move to trash")
	flag.Parse()

	if *configPath == "" {
		return nil, nil, errors.New("please provide a configuration file")
	}

	cmdArgs = new(file_cleaner.CmdLineArgs)
	cmdArgs.DryRun = *dryRun
	if *dryRun {
		fmt.Println("Running in dry-run mode")
	}

	cmdArgs.ReplaceAsSymlink = *replaceAsSymlink
	if *replaceAsSymlink {
		fmt.Println("Replacing duplicate files with symlinks and moving to trash")
	}

	config = new(file_cleaner.Config)
	err = config.Load(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration file:", err)
		return nil, nil, err
	}

	return config, cmdArgs, nil
}

func main() {
	config, cmdArgs, err := parseArgs()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		os.Exit(1)
	}

	lock_file := flock.New("/tmp/file_cleaner.lock")
	locked, err := lock_file.TryLock()
	if err != nil {
		fmt.Println("Error locking file:", err)
		os.Exit(1)
	}

	if locked {
		if err = config.Execute(*cmdArgs); err != nil {
			fmt.Println("Error executing configuration:", err)
			os.Exit(1)
		}

		lock_file.Unlock()
	} else {
		fmt.Println("Another instance of file_cleaner is already running")
		os.Exit(1)
	}
}
