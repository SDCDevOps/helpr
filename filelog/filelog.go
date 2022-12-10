package filelog

import (
	"os"
	"time"
)

type Logger struct {
	filename string
}

// New - Initialise a new object.
// Note:
// - If log file already exist, it will overwrite if calling LogNew(), or append if using calling LogAppend().
// - If log file does not exist, it will create it when calling LogNew() or LogAppend().
func New(filename string) (logr Logger) {
	logr.filename = filename
	return
}

// Empty file and write "content" (file will just contain "content" data).
func (lgr *Logger) LogNew(content string) error {
	file, err := os.OpenFile(lgr.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Empty the file.
	if err = os.Truncate(lgr.filename, 0); err != nil {
		return err
	}

	err = log(file, content)
	if err != nil {
		return err
	}

	return nil
}

// Append content to file.
func (lgr *Logger) LogAppend(content string) error {
	file, err := os.OpenFile(lgr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = log(file, content)
	if err != nil {
		return err
	}

	return nil
}

// Log content with timestamp.
func log(file *os.File, content string) error {
	currTime := time.Now()
	ct := "\n[" + currTime.Format("2006-01-02 15:04:05") + "] " + content

	if _, err := file.WriteString(ct); err != nil {
		return err
	}

	return nil
}
