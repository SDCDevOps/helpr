package helpr

import (
	"os"
	"testing"

	"github.com/SDCDevOps/helpr/filemgr"
	"github.com/SDCDevOps/helpr/str"
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
	t.Log("\n\nTest filemgr package...")

	filename := "test.txt"
	content1 := "CONTENT1 CONTENT1"

	t.Log("Creating file with content1, overwriting if exist...")
	err := filemgr.CreateFileIfNotExist(filename, content1, true)
	if err != nil {
		t.Fatal("Error calling CreateFileIfNotExist")
	}

	t.Log("Check that file with content1 exist...")
	notExist, err := filemgr.FileNotExist(filename)
	if err != nil {
		t.Fatal("Error calling FileNotExist")
	}

	assert.Equal(t, false, notExist, "FileNotExist should return false after creating file")

	t.Log("Creating file with content2, NOT TO OVERWRITE if exist...")
	content2 := "This is content2 CONTENT2"
	err = filemgr.CreateFileIfNotExist(filename, content2, false)
	if err != nil {
		t.Fatal("Error calling CreateFileIfNotExist")
	}

	t.Log("Read file content after calling CreateFileIfNotExist with overwriteIfExist=FALSE...")
	byteData, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading file content")
	}

	assert.Equal(t, content1, string(byteData), "File content should still remain as content1")

	t.Log("Creating file with content2, to OVERWRITE if exist...")
	err = filemgr.CreateFileIfNotExist(filename, content2, true)
	if err != nil {
		t.Fatal("Error calling CreateFileIfNotExist")
	}

	t.Log("Read file content after calling CreateFileIfNotExist with overwriteIfExist=TRUE...")
	byteData, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading file content")
	}

	assert.Equal(t, content2, string(byteData), "File content should change to content2")

	err = filemgr.AppendFileCreateIfNotExist(filename, content1)
	if err != nil {
		t.Fatal("Error calling AppendFileCreateIfNotExist")
	}

	t.Log("Read file content after calling AppendFileCreateIfNotExist...")
	byteData, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading file content")
	}
	assert.Equal(t, content2+content1, string(byteData), "File content should be content1 appended to content2")

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
