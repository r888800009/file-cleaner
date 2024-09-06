package file_cleaner

import (
	"fmt"
	"os"
	"path/filepath"
)

// it is the argument for cmd line
type CmdLineArgs struct {
	DryRun           bool
	ReplaceAsSymlink bool
}

type ExecuteArgs struct {
	cmd    CmdLineArgs
	config Config
}

/*
Strategy defines the what to do with the file
*/
type Strategy interface {
	Load(name string, value map[string]interface{}) error
	Execute(parms ExecuteArgs) error
}

func ListFiles(dirEntry DirEntry) (map[int64]([]FileEntry), map[string]FileEntry) {
	recursively := dirEntry.recursively

	// list all target files and create a map of size to file, is can chceck quickly if a file exists without reading the file
	sizeIndex := make(map[int64]([]FileEntry))
	fileMap := make(map[string]FileEntry)
	filepath.Walk(dirEntry.path, func(path string, info os.FileInfo, err error) error {
		// check if file is directory
		if info.IsDir() && !recursively && path != dirEntry.path {
			return filepath.SkipDir
		}

		// check is not match
		if !dirEntry.Match(path) {
			return nil
		}

		entry := FileEntry{}
		entry.Load(path)

		// make index and map
		if sizeIndex[entry.size] == nil {
			sizeIndex[entry.size] = []FileEntry{entry}
		} else {
			sizeIndex[entry.size] = append(sizeIndex[entry.size], entry)
		}
		fileMap[path] = entry

		return nil
	})

	return sizeIndex, fileMap
}

func duplicateHandler(clean FileEntry, keep FileEntry, parms ExecuteArgs, strategy SourceToTargetDedupeStrategy) {
	fmt.Println("  Duplicate:", clean.path)
	fmt.Println("    Target:", keep.path)

	// move to trash
	absPath, _ := filepath.Abs(clean.path)
	trashPath := filepath.Join(strategy.trashPath, absPath)
	fmt.Println("    Moving to trash:", clean.path)
	fmt.Println("    Trash Path:", trashPath)
	if !parms.cmd.DryRun {
		os.MkdirAll(filepath.Dir(trashPath), os.ModePerm)
		os.Rename(clean.path, trashPath)
	} else {
		fmt.Println("    Dry Run: Not moving to trash")
	}

	if parms.cmd.ReplaceAsSymlink {
		fmt.Println("    Replacing with symlink:", clean.path, "->", keep.path)
		if !parms.cmd.DryRun {
			os.Symlink(keep.path, clean.path)
		} else {
			fmt.Println("    Dry Run: Not creating symlink")
		}
		// create symlink
	}

}

func IsPathNotIndepent(path1 string, path2 string) (bool, error) {
	path1 = filepath.Clean(path1)
	path2 = filepath.Clean(path2)
	path1, err := filepath.Abs(path1)
	if err != nil {
		return true, err
	}

	path2, err = filepath.Abs(path2)
	if err != nil {
		return true, err
	}

	// add separator to the end of the path
	path1 += string(filepath.Separator)
	path2 += string(filepath.Separator)

	if path1 == path2 {
		return true, nil
	}

	// set path1 to be the shorter one
	if len(path1) > len(path2) {
		path1, path2 = path2, path1
	}

	// check if path1 is subpath of path2
	if path2[:len(path1)] == path1 {
		return true, nil
	}

	return false, nil
}

func (strategy *SourceToTargetDedupeStrategy) Execute(parms ExecuteArgs) error {
	fmt.Println("Execute SourceToTargetDedupeStrategy")
	fmt.Println("Target:", strategy.target.path)

	// if target directory does not exist, throw an error
	if _, err := os.Stat(strategy.target.path); os.IsNotExist(err) {
		return err
	}

	sizeIndex, fileMap := ListFiles(strategy.target)

	// print all target files
	for path := range fileMap {
		fmt.Println("  Target:", path)
	}

	for _, source := range strategy.source {
		notIndepent, err := IsPathNotIndepent(source.path, strategy.target.path)
		if err != nil {
			return err
		}

		if notIndepent {
			panic("current not support source and target are the same")
		}

		fmt.Println("Source:", source.path)
		_, sourceFileMap := ListFiles(source)

		// print all source files
		for _, entry := range sourceFileMap {
			// check if file duplicates
			if targetEntries, ok := sizeIndex[entry.size]; ok {
				for _, targetEntry := range targetEntries {
					if entry.Equal(&targetEntry) && entry.path != targetEntry.path {
						duplicateHandler(entry, targetEntry, parms, *strategy)
					}
				}
			}
		}
	}
	return nil
}

func (config_struct *Config) Execute(cmdLineArgs CmdLineArgs) error {
	parms := ExecuteArgs{cmd: cmdLineArgs, config: *config_struct}

	for name, strategy := range config_struct.strategies {
		fmt.Println("Execute:", name)
		if err := strategy.Execute(parms); err != nil {
			return err
		}
	}
	return nil
}
