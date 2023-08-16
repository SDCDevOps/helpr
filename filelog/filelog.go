package filelog

import (
	"os"
	"time"

	"github.com/SDCDevOps/helpr/filemgr"
)

type Logger struct {
	filepath string
}

// New - Initialise a new object.
// Note:
//   - If log file already exist, it will overwrite if calling LogNew(), or append if using calling LogAppend().
//   - If log file does not exist, it will create it when calling LogNew() or LogAppend().
//   - If using New(), maxSize > 0, file already exist and it's size is >= maxSize, it will rename existing file to
//     <filepath>-<current timestamp>, and then create a new file named <filepath>.
//   - If maxSize is 0 or less, it will not check if file exist and not check the size.
func New(filepath string, maxSize int64) (Logger, error) {
	var logr Logger
	logr.filepath = filepath

	if maxSize > 0 {
		isNotExist, err := filemgr.FileNotExist(filepath)
		if err != nil {
			return logr, err
		}

		if !isNotExist {
			stat, _ := os.Stat(filepath)
			if stat.Size() >= maxSize {
				os.Rename(filepath, filepath+"-"+time.Now().Format("20060102_150405")) // Rename log file with timestamp.
			}
		}
	}

	return logr, nil
}

// Empty file and write "content" (file will just contain "content" data).
func (lgr *Logger) LogNew(content string) error {
	file, err := os.OpenFile(lgr.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Empty the file.
	if err = os.Truncate(lgr.filepath, 0); err != nil {
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
	file, err := os.OpenFile(lgr.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

// Append content to file and panic/exit program.
func (lgr *Logger) LogAppendPanic(content string) {
	lgr.LogAppend(content)
	panic(content)
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
