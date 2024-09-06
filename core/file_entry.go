package file_cleaner

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type Md5Sum []byte

type FileEntry struct {
	name  string
	path  string
	isDir bool
	size  int64
	md5   Md5Sum
}

/*
Print prints the file entry information to the console.
*/
func (entry *FileEntry) Print() {
	fmt.Println("File name:", entry.name)
	fmt.Println("File path:", entry.path)
	fmt.Println("Is directory:", entry.isDir)
	fmt.Println("File size:", entry.size)

	// Lazy load MD5
	md5, err := entry.MD5()
	if err != nil {
		fmt.Println("Error calculating MD5:", err)
	} else {
		fmt.Printf("MD5: %x\n", md5)
	}
}

/*
Create a new FileEntry object and load the file information from the given path.
Note: The MD5 hash is not calculated until the FileEntry.MD5() method is called.
*/
func (entry *FileEntry) Load(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	entry.name = fileInfo.Name()
	entry.path = path
	entry.isDir = fileInfo.IsDir()
	entry.size = fileInfo.Size()

	// lazy load md5
	entry.md5 = nil

	return nil
}

/*
MD5 calculates the MD5 hash of the file and returns it as a byte slice.
The MD5 hash is cached after the first call to this method.
*/
func (entry *FileEntry) MD5() (Md5Sum, error) {
	if entry.md5 != nil {
		return entry.md5, nil
	}

	file, err := os.Open(entry.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	entry.md5 = hash.Sum(nil)
	return entry.md5, nil
}

// compareBufferSize is the size of the buffer used to compare file contents.
const compareBufferSize = 1024

/*
Compare compares the two file is the same or not.
it would compare the size and md5 first.
then compare the content of the file.
*/
func (entry *FileEntry) Compare(other *FileEntry) bool {
	if entry.size != other.size {
		return false
	}

	// check md5
	md5, err := entry.MD5()
	if err != nil {
		return false
	}

	otherMd5, err := other.MD5()
	if err != nil {
		return false
	}

	if fmt.Sprintf("%x", md5) != fmt.Sprintf("%x", otherMd5) {
		return false
	}

	// open file and compare the content
	file1, err := os.Open(entry.path)
	file2, err2 := os.Open(other.path)
	if err != nil || err2 != nil {
		return false
	}
	defer file1.Close()
	defer file2.Close()

	for {
		block1 := make([]byte, compareBufferSize)
		size1, err1 := file1.Read(block1)
		block2 := make([]byte, compareBufferSize)
		size2, err2 := file2.Read(block2)

		// check if it is the end of the file
		if err1 == io.EOF && err2 == io.EOF {
			return true
		}

		if err1 != nil || err2 != nil {
			return false
		}

		if size1 != size2 {
			return false
		}

		// check block content is the same
		if !bytes.Equal(block1, block2) {
			return false
		}
	}
}

/*
Equal is an alias for FileEntry.Compare()
*/
func (entry *FileEntry) Equal(other *FileEntry) bool {
	return entry.Compare(other)
}
