package file_cleaner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	file_cleaner "github.com/r888800009/file_cleaner/core"
	"github.com/stretchr/testify/assert"
)

func TestEntryEqual(t *testing.T) {
	assert := assert.New(t)
	var entry file_cleaner.FileEntry
	var entry2 file_cleaner.FileEntry

	entry.Load("/etc/hosts")
	entry2.Load("/etc/hosts")
	assert.True(entry.Equal(&entry2))

	entry.Load("/etc/hosts")
	entry2.Load("/etc/passwd")
	assert.False(entry.Equal(&entry2))
}

func TestIsPathNotIndepent(t *testing.T) {
	assert := assert.New(t)

	assert.True(file_cleaner.IsPathNotIndepent("/etc/hosts", "/etc"))
	assert.True(file_cleaner.IsPathNotIndepent("/etc/hosts", "/etc/hosts"))
	assert.True(file_cleaner.IsPathNotIndepent("/etc/hosts", "/etc/hosts/"))
	assert.True(file_cleaner.IsPathNotIndepent("/home/user////file", "/home/user"))
	assert.True(file_cleaner.IsPathNotIndepent("./file", "."))

	assert.False(file_cleaner.IsPathNotIndepent("/etc/hosts", "/etc/passwd"))
	assert.False(file_cleaner.IsPathNotIndepent("./file", "/home/user"))
	assert.False(file_cleaner.IsPathNotIndepent("/etc/host", "/etc/hostname"))

	// cd should not affect the result
	// note: the test may depend on operating system, for example, in macos would using symlink for /tmp and /home
	// check /bin is not symlink
	if link, _ := os.Readlink("/bin"); link != "" {
		assert.FailNow("The /bin is symlink, the test may not work")
	}

	os.Chdir("/bin")
	assert.True(file_cleaner.IsPathNotIndepent(".", "/bin"))
	assert.True(file_cleaner.IsPathNotIndepent("./", "/bin/vim"))
}

func TestPathNomalizePair(t *testing.T) {
	assert := assert.New(t)
	path1, path2, err := file_cleaner.PathNomalizePair("/etc/hosts", "/etc")
	path1, path2, swapped := file_cleaner.SetShorterPathFirst(path1, path2)
	assert.Nil(err)

	// first path should be shorter
	assert.Equal("/etc/", path1)
	assert.Equal("/etc/hosts/", path2)
	assert.True(swapped)
}

/*
this test consider if the path recursive=false
*/
func TestIsPathNotIndepentRecursiveFalse(t *testing.T) {
	assert := assert.New(t)

	// if we recursive /etc, but not recursive /etc/hosts, it not independent
	assert.True(file_cleaner.IsPathNotIndepentRecursive("/etc/hosts", false, "/etc", true))

	// swap the path, it should be the same result
	assert.True(file_cleaner.IsPathNotIndepentRecursive("/etc", true, "/etc/hosts", false))

	// if we recursive /etc/hosts, but not recursive /etc, it is independent
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/etc/hosts", true, "/etc", false))

	// if same path, it is not independent
	assert.True(file_cleaner.IsPathNotIndepentRecursive("/etc/hosts", false, "/etc/hosts", false))

	// if /etc/ not recursive, and /etc/dir not recursive, it is independent
	// because we search /etc/* not include /etc/dir/*
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/etc/dir", false, "/etc", false))

	// diffent path, it is independent
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/tmp", false, "/etc", false))

	// same name but different root, it is independent
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/tmp", false, "/user/tmp", false))
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/tmp", true, "/user/tmp", true))
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/tmp", false, "/user/tmp", true))
	assert.False(file_cleaner.IsPathNotIndepentRecursive("/tmp", true, "/user/tmp", false))
}

// test ListFiles should not return the directory itself
func TestListFiles(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)
	t.Logf("Current test directory: %s", filepath.Dir(filename))

	// change to the test directory
	os.Chdir(filepath.Dir(filename))

	// test recursive
	dirEntry := file_cleaner.CreateDirEntry("data/listfile/", true)
	_, fileMap := file_cleaner.ListFiles(dirEntry)
	assert.NotContains(fileMap, "data/listfile/")
	assert.NotContains(fileMap, "data/listfile/dir/")
	assert.Contains(fileMap, "data/listfile/dir/listfile")

	// test not recursive
	dirEntry = file_cleaner.CreateDirEntry("data/listfile/", false)
	_, fileMap = file_cleaner.ListFiles(dirEntry)
	assert.NotContains(fileMap, "data/listfile/")
}
