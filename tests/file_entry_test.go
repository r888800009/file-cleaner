package file_cleaner

import (
	"os"
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
