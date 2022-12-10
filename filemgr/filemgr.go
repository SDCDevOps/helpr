package filemgr

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Return true if file does not exist.
func FileNotExist(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}

		return false, err // There's other error.
	}

	return false, nil // File exist.
}

// Create file if file does not exist.
// If overwriteIfExist = false and file exist, nothing will happen (content will not be written to existing file)
func CreateFileIfNotExist(filename string, content string, overwriteIfExist bool) error {
	notExist, err := FileNotExist(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Error checking if file exist: %s", err.Error()))
	}

	if notExist || overwriteIfExist {
		file, err := os.Create(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("Error creating file: %s", err.Error()))
		}
		defer file.Close()

		if content != "" {
			_, err := file.WriteString(content)
			if err != nil {
				return errors.New(fmt.Sprintf("Error writing to file: %s", err.Error()))
			}
		}
	}

	return nil
}

// Delete file if exist.
func DeleteFileIfExist(filename string) error {
	notExist, err := FileNotExist(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Error checking whether file exist: %s", err.Error()))
	}

	if notExist {
		log.Println("File does not exist:", filename)
	} else {
		err := os.Remove(filename) // Delete file.
		if err != nil {
			return errors.New(fmt.Sprintf("Error deleting file: %s", err.Error()))
		}
		log.Println("File exists and deleted:", filename)
	}

	return nil
}
