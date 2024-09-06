package file_cleaner

import (
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
