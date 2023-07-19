package helpr

import (
	"os"
	"testing"

	"github.com/SDCDevOps/helpr/crypt"
	"github.com/SDCDevOps/helpr/filelog"
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

func TestFilelog(t *testing.T) {
	t.Log("\n\nTest filelog package...")

	filename := "mylog.log"
	fl := filelog.New(filename)

	t.Log("Calling LogAppend to log content1...")
	content1 := "content1 CONTENT1"
	err := fl.LogAppend(content1)
	if err != nil {
		t.Fatal("Error calling LogAppend")
	}

	t.Log("Reading content in log file after calling LogAppend...")
	byt, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading log file after calling LogAppend")
	}

	assert.Contains(t, string(byt), content1, "Log file should contain content1")

	t.Log("Calling LogNew to log content2...")
	content2 := "content2 CONTENT2"
	err = fl.LogNew(content2)
	if err != nil {
		t.Fatal("Error calling LogNew")
	}

	t.Log("Reading content in log file after calling LogNew...")
	byt, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading log file after calling LogNew")
	}

	str := string(byt)
	assert.Contains(t, str, content2, "Log file should contain content2")
	assert.NotContains(t, str, content1, "Log file should NO LONGER CONTAIN content1")

	t.Log("Calling LogAppend again to log content1...")
	err = fl.LogAppend(content1)
	if err != nil {
		t.Fatal("Error calling LogAppend")
	}

	t.Log("Reading content in log file after calling LogAppend again...")
	byt, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading log file after calling LogAppend again")
	}

	str = string(byt)
	assert.Contains(t, str, content1, "Log file should contain content1")
	assert.Contains(t, str, content2, "Log file should contain content2")

	// Clean up: removing log file.
	err = os.Remove(filename)
	if err != nil {
		t.Fatal("Error deleting log file")
	}
}

func TestCrypt(t *testing.T) {
	t.Log("\n\nTest str package...")

	keyInvalid := []byte("qwertyuiopasdfghjklzxcvbnmlkjh") // 30 bytes
	keyValid := []byte("qwertyui10asdfghjk20xcvbnmlk30ab") // 32 bytes
	fileContent := "This is a test file\nThis is another line"
	srcClearFile := "srcClearFile.txt"
	encryptedFile := "encryptedFile.bin"
	decryptedFile := "decryptedFile.txt"

	t.Log("Creating srcClearFile file...")
	f, err := os.OpenFile(srcClearFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal("Error creating srcClear file: " + err.Error())
	}
	defer f.Close()

	f.WriteString(fileContent)

	t.Log("Testing with AES key with invalid length...")
	err = crypt.EncryptAes(keyInvalid, srcClearFile, encryptedFile)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "error in AES key")
	}

	t.Log("Testing encryption...")
	err = crypt.EncryptAes(keyValid, srcClearFile, encryptedFile)
	assert.Nil(t, err)

	t.Log("Check that encrypted file created...")
	fileNotExist, err := filemgr.FileNotExist(encryptedFile)
	if err != nil {
		t.Fatal("Error checking encryptedFile existence")
	}

	assert.Equal(t, false, fileNotExist, "Encrypted file should be created")

	t.Log("Testing decryption...")
	err = crypt.DecryptAes(keyValid, encryptedFile, decryptedFile)
	assert.Nil(t, err)

	t.Log("Check that decrypted file created...")
	fileNotExist, err = filemgr.FileNotExist(decryptedFile)
	if err != nil {
		t.Fatal("Error checking decryptedFile existence")
	}

	assert.Equal(t, false, fileNotExist, "Decrypted file should be created")

	t.Log("Reading decrypted file content...")
	data, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatal("Error reading decryptedFile file")
	}

	assert.Equal(t, fileContent, string(data))

	// Clean up.
	err = filemgr.DeleteFileIfExist(srcClearFile)
	if err != nil {
		t.Fatal("Error deleting srcClearFile file")
	}

	err = filemgr.DeleteFileIfExist(encryptedFile)
	if err != nil {
		t.Fatal("Error deleting encryptedFile file")
	}

	err = filemgr.DeleteFileIfExist(decryptedFile)
	if err != nil {
		t.Fatal("Error deleting decryptedFile file")
	}

}
