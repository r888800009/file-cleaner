# File Cleaner

a go application to clean files, scans a directory finds duplicates and deletes them.

## Installation
```bash
git clone https://github.com/r888800009/file-cleaner
cd file-cleaner
go build
```
## Usage
```bash
./file-cleaner -config path/to/config.json
```

## Configuration
`target_dir_strategy` would search the `source_dirs` files if it exists in the `target_dir` and deletes or symlinks them.
```json
{
    "version": "0.1",
    "target_dir_strategy": {
        "target_dir": "path/to/target/dir",
        "target_dir_recursive": true,
        "source_dirs": [
            {
                "source_dir": "path/to/source/dir",
                "recursive": true
            },
            {
                "source_dir": "path/to/source/dir",
                "recursive": true
            }
        ]
    }
}
```
`pdf-mover` would move matching files from `source_dir` to `target_dir` based on the `pdf_matcher` configuration.
- `conference_paper_detector` match keywords such as `usenix`, `ieee`, `acm`, etc. or other heuristics.
```json
{
    "version": "0.1",
    "pdf_mover": {
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
go test -v ./...
```