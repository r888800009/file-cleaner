# File Cleaner

a go application to clean files, scans a directory finds duplicates and deletes them.

## Installation
```bash
git clone https://github.com/r888800009/file-cleaner
cd file-cleaner
go build
```
## Note: currently file_cleaner is under development, it might not work as expected. or have some bugs cause data loss. you should backup your data before using it. we are not responsible for any data loss.

## Usage
```bash
./file_cleaner -config path/to/config.json
```
for safety, we make the cli tool dry-run by default, you can use `-dry-run=false` to move file to trash.
```bash
./file_cleaner -config path/to/config.json -dry-run=false
```
if you want remove empty trash directory, you can use `find` command to remove them.
```bash
find ./trash -type d -empty -delete
```

## Features
- [x] `source_to_target_dedupe`
- [ ] `pdf_mover`

## Configuration
`source_to_target_dedupe` would search the `source_dirs` files if it exists in the `target_dir` and deletes or symlinks them.
duplicate files are moved to the `trash_dir` and keep the newest file. `ignore` is supported go regex.
```json
{
    "version": "0.1",
    "name1": {
        "strategy": "source_to_target_dedupe",
        "target_dir": {
            "path": "~/organized_dir",
            "recursive": true,
            "ignore": "(.*/.git.*|.*/__init__.py)"
        },
        "trash_dir": "~/trash",
        "source_dirs": [
            {
                "path": "~/Downloads",
                "recursive": true,
                "ignore": "regex"
            },
            {
                "path": "path/to/source/dir",
                "recursive": true
            }
        ]
    }
}
```
`target_dir` and `source_dir` are `DirEntry` struct, that allows you to specify the `path` and `recursive` flag.
```json
{
    "path": "path/to/source/dir",
    "recursive": true
}
```
for `ignore` you can use regex to ignore files or directories.
```json
{
    "path": "~/organized_dir",
    "recursive": true,
    "ignore": "(.*/.git.*|.*/__init__.py)"
}
```
for `match` you can use regex to match files or directories.
```json
{
    "path": "~/organized_dir",
    "recursive": true,
    "match": "(.*/.pdf)"
}
```
if you set match and ignore same time, it means the file must match the `match` and not match the `ignore`. 
note it might be slow if you have a lot of files. because we check two regex for each file.
```json
{
    "path": "~/organized_dir",
    "recursive": true,
    "match": "(.*/.pdf)",
    "ignore": "(.*/.git.*|.*/__init__.py)"
}
```

trash_dir directory structure, each directory is a timestamp (ISO 8601) of the deletion. you can recover the files if needed.
```
- trash
    - YYYY-MM-DD-HH-MM-SS.sss
        - original/file/path
        - original/file2/path
    - YYYY-MM-DD-HH-MM-SS.sss
        - original/file/path
        - original/file2/path
```

`pdf-mover` would move matching files from `source_dir` to `target_dir` based on the `pdf_matcher` configuration.
- `conference_paper_detector` match keywords such as `usenix`, `ieee`, `acm`, etc. or other heuristics.
```json
{
    "version": "0.1",
    "name2": {
        "strategy": "pdf_mover",
        "pdf_matcher": "conference_paper_detector",
        "target_dir": {
            "path": "path/to/target/dir",
            "recursive": true
        },
        "source_dir": {
            "path": "path/to/source/dir",
            "recursive": true
        }
    }
}
```
it supports multiple entries
```json
{
    "version": "0.1",
    "name3": {
        "strategy": "source_to_target_dedupe",
        ...
    },
    "name4": {
        "strategy": "pdf_mover",
        ...
    },
    ...
}
```

## Development
generate the godoc, if you don't have godoc installed
```bash
go get golang.org/x/tools/cmd/godoc
```
run the godoc server
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
godoc -http=:6060
```
[open the godoc](http://localhost:6060/pkg/github.com/r888800009/file_cleaner/)

run the tests
```bash
go test -v ./... -cover -coverpkg=./...
```

using as a library currently not supported, but it might possible call the `file_cleaner` as a go package.
there no guarantee of the stability of the API.
```go
package main

import (
    "github.com/r888800009/file_cleaner"
)
```

download specific version of the file-cleaner
```bash
go get github.com/r888800009/file_cleaner@v0.1.0
```

publish a new version
```bash
git tag v0.1.0
git push origin v0.1.0
```

## Related Projects
- [adrianlopezroche/fdupes: FDUPES is a program for identifying or deleting duplicate files residing within specified directories.](https://github.com/adrianlopezroche/fdupes)
    - `file_cleaner` support differentiated the `source` and `target` directories, found `source` files in the `target` directory and delete them. it supports regex for ignore/match files.
    - `fdupes` is cli tool, you might need to write a script to parse the output and delete the files.
- [markfasheh/duperemove: Tools for deduping file systems](https://github.com/markfasheh/duperemove)
    - duperemove is filesystem layer deduplication. `file_cleaner` is a higher level deduplication tool, it don't care about the filesystem layer.