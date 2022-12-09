package helpr

import (
	"testing"

	"github.com/skker/helpr/filemgr"
	"github.com/skker/helpr/str"
	"github.com/stretchr/testify/assert"
)

func TestStr(t *testing.T) {
	t.Log("\n\nTest str package...")

	s := "This is My String to Search"

	substr := "sea"
	result := str.CaseInsensitiveSearch(s, substr)
	assert.True(t, result, "String ("+s+") should contain "+substr)

	substr = "land"
	result = str.CaseInsensitiveSearch(s, substr)
	assert.False(t, result, "String ("+s+") should not contain "+substr)
}

func TestFilemgr(t *testing.T) {
	t.Log("\n\nTest str package...")

	filename := "test.txt"
	content := "test test only..."
	err := filemgr.CreateFileIfNotExist(filename, content)
	if err != nil {
		t.Fatal("Error calling CreateFileIfNotExist")
	}

	notExist, err := filemgr.FileNotExist(filename)
	if err != nil {
		t.Fatal("Error calling FileNotExist")
	}

	assert.Equal(t, false, notExist, "FileNotExist should return false after creating file")

	err = filemgr.DeleteFileIfExist(filename)
	if err != nil {
		t.Fatal("Error calling DeleteFileIfExist")
	}

	notExist, err = filemgr.FileNotExist(filename)
	if err != nil {
		t.Fatal("Error calling FileNotExist")
	}

	assert.Equal(t, true, notExist, "FileNotExist should return true after calling DeleteFileIfExist")
}
