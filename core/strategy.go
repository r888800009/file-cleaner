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
	includeDirs := dirEntry.include_dirs

	// list all target files and create a map of size to file, is can chceck quickly if a file exists without reading the file
	sizeIndex := make(map[int64]([]FileEntry))
	fileMap := make(map[string]FileEntry)
	filepath.Walk(dirEntry.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if file is directory
		if info.IsDir() && !recursively && path != dirEntry.path {
			return filepath.SkipDir
		}

		// check is not match
		if !dirEntry.Match(path) {
			return nil
		}

		// skip if it is directory
		if info.IsDir() && !includeDirs {
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

func pathNomalize(path string) (string, error) {
	path = filepath.Clean(path)
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path += string(filepath.Separator)
	return path, nil
}

/*
PathNomalizePair normalize the path and return the normalized path
*/
func PathNomalizePair(path1 string, path2 string) (string, string, error) {
	path1, err := pathNomalize(path1)
	if err != nil {
		return "", "", err
	}

	path2, err = pathNomalize(path2)
	if err != nil {
		return "", "", err
	}

	// return: set path1 to be the shorter one
	return path1, path2, nil
}

// check if path1 is subpath of path2
// it path1 should be shorter than path2
func checkIsSubPath(path1 string, path2 string) bool {
	return path2[:len(path1)] == path1
}

func SetShorterPathFirst(path1 string, path2 string) (string, string, bool) {
	swapped := false
	// set path1 to be the shorter one
	if len(path1) > len(path2) {
		path1, path2 = path2, path1
		swapped = true
	}
	return path1, path2, swapped
}

func IsPathNotIndepent(path1 string, path2 string) (bool, error) {
	path1, path2, err := PathNomalizePair(path1, path2)
	if err != nil {
		return true, err
	}

	path1, path2, _ = SetShorterPathFirst(path1, path2)
	return checkIsSubPath(path1, path2), nil
}

func IsPathNotIndepentRecursive(path1 string, recPath1 bool, path2 string, recPath2 bool) (bool, error) {
	if recPath1 && recPath2 {
		return IsPathNotIndepent(path1, path2)
	}

	path1, path2, err := PathNomalizePair(path1, path2)
	if err != nil {
		return true, err
	}

	// make sure path1 is shorter
	path1, path2, swapped := SetShorterPathFirst(path1, path2)
	if swapped {
		recPath1, recPath2 = recPath2, recPath1
	}

	// if we recursive /etc, but not recursive /etc/hosts, it not independent
	if recPath1 && !recPath2 {
		return checkIsSubPath(path1, path2), nil
	}

	// if we not recursive /etc, but recursive /etc/hosts, it is independent
	if !recPath1 && recPath2 {
		return false, nil
	}

	// if we not recursive path1 and path2, it is independent only if path1 not equal to path2
	return path1 == path2, nil
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
		notIndepent, err := IsPathNotIndepentRecursive(source.path, source.recursively, strategy.target.path, strategy.target.recursively)
		if err != nil {
			return err
		}

		if notIndepent {
			panic(fmt.Sprintf("current not support source and target are the same %s %s", source.path, strategy.target.path))
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
