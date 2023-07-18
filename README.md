# helpr
**helpr** is a simple utility package in Go.

It contains the following packages:  
* mgdb  
* rstatus
* str

---
## Installation

Install the dependency in your project.

    go get github.com/SDCDevOps/helpr

---
## mgdb

**mgdb** provides some simple methods to use MongoDB in Go. 

It is created to save you a few lines of codes. It's not necessary to use it.

### Usage

    import "github.com/SDCDevOps/helpr/mgdb"
    
    func myFunc() error {
      // Context with timeout.
      ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
      defer cancel()
      
      // Initialise.
      m, err := mgdb.New(ctx, "mongodb://myDbURI:1234", "myDbName")
      if err != nil {
        return err
      }
      defer m.Close(ctx) // Remember to close.
      
      db := m.Db
      
      // Use it as normal mongo calls.
      coll := db.Collection("humanBody")
      if _, err := coll.DeleteMany(ctx, bson.M{}); err != nil {
        return err
      }
      
      // Check document existence.
      isThere, err := mgdb.DocExists(coll, bson.M{"bodypart": "myhead"})
      if err != nil {
        return err
      }
      
      if !isThere {
        fmt.Print("Oh no! Where's my head?!")
      } else {
        fmt.Print("Thought I told you to cut off everything?!")
      }
    }

---
## rstatus

**rstatus** can be used as a function return status that can make the returned error more specific.

It's able to return a relevant http status code based on the error.

### Usage

    import "github.com/SDCDevOps/helpr/rstatus"

    func main() {
      s := stickYourHandInHere("left")
      if s.Err != nil { // Check if there's error.
        log.Printf("There's an error: %v", s.Message)
        
        if s.Type == rstatus.ReturnStatusFailDB {
          fixDatabase()
        } else if s.Type == rstatus.ReturnStatusFailExternalParty {
          screamAtTheWorld()
        } else if s.Type == rstatus.ReturnStatusFailInternal {
          goForCheckup()
        } else if s.Type == rstatus.ReturnStatusFailInvalidInput {
          checkInMentalHospital()
        }
        
      } else { // If no error, we can assume that function call is successful.
        announceYourHappinessToTheWorld(s.Message)
      }
      
      httpStatusCode := s.GetHTTPStatus(); // Get relevant HTTP status code.
      returnServiceCallUsingCode(httpStatusCode)
    }

    func stickYourHandInHere(hand string) rstatus.Status {
      if hand == "left" {
        stuff, err := getMysteryStuffFromDatabase(hand)
        if err != nil {
          return rstatus.New(rstatus.ReturnStatusFailDB, "Database error", err) // DB error
        }
        
        msg := fmt.Sprintf("Congrats! You got yourself a %v from our database!", stuff)
        return rstatus.New(rstatus.ReturnStatusSuccess, msg, nil) // No error
      } else if hand == "right" {
        stuff, err := grabStuffFromOutside(hand)
        if err != nil {
          return rstatus.New(rstatus.ReturnStatusFailExternalParty, "External party error", err) // External party error
        }
        
        msg := fmt.Sprintf("Congrats! You got yourself a %v from our involuntary sponsor!", stuff)
        return rstatus.New(rstatus.ReturnStatusSuccess, msg, nil) // No error
      } else if "fakehand" {
        msg := "Internal error: We're having indigestion..."
        err := errors.New(msg)
        return rstatus.New(rstatus.ReturnStatusFailInternal, msg, err) // Internal error.
      } else {
        msg := "What the hell did you put in here?!"
        err := errors.New(msg)
        return rstatus.New(rstatus.ReturnStatusFailInvalidInput, msg, err)
      }
    }


---
## str

**str** is a string utility.

### Usage

    import "github.com/SDCDevOps/helpr/str"

    func main() {
      s := "This is My String to Search"
      substr := "sea"

      if str.CaseInsensitiveSearch(s, substr) {
        log.Print("Result is true. This will be printed") // Will only print this.
      } else { 
        log.Print("Result is false. This will NEVER be printed")
      }
    }


---
## filemgr

**filemgr** is a file utility.

### Usage

    import "github.com/SDCDevOps/helpr/filemgr"

    func main() {
      content1 := "This is my CONTENT1"
      filename := "test.txt"

      // File will be created if it does not exist. 
      // If file already exist and overwriteIfExist=false (3rd param), nothing will happen (content will NOT be written into it).
      // If file already exist and overwriteIfExist=true, file will be truncated and content (2nd param) will be written into it.
      err := filemgr.CreateFileIfNotExist(filename, content1, true) 
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling CreateFileIfNotExist: %s", err.Error()))
      }

      content2 := "My CONTENT2"
      err = filemgr.AppendFileCreateIfNotExist(filename, content2) // content2 will be appended to file.
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling AppendFileCreateIfNotExist: %s", err.Error()))
      }
      
      notExist, err := filemgr.FileNotExist(filename) // notExist will be false since file is created earlier.
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling FileNotExist: %s", err.Error()))
      }
      
      err = filemgr.DeleteFileIfExist(filename) // File will be deleted.
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling DeleteFileIfExist: %s", err.Error()))
      }
    }


---
## filelog

**filelog** is a file logging utility.

### Usage

    import "github.com/SDCDevOps/helpr/filelog"

    func main() {
      filename := "mylog.log"
      content1 := "content1 CONTENT1"

      fl := filelog.New(filename) // Initiate object.

      // Calling LogAppend will append to log (with timestamp). If log file does not exist, it will create and log content1.
      // If log file already exist, it will append content1 to it.
      err := fl.LogAppend(content1)
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling LogAppend: %s", err.Error()))
      }

      // Calling LogNew will log as new (with timestamp). If log file does not exist, it will create and log content1.
      // If log file already exist, it will truncate it and then log content1 (thus only content1 will in the log file).
      err := fl.LogNew(content1)
      if err != nil {
        log.Panic(fmt.Sprintf("Error calling LogNew: %s", err.Error()))
      }

      // Calling LogAppendPanic will append to log (with timestamp) and "panic" (exit) program. If log file does not exist, it will create and log content1.
      // If log file already exist, it will append content1 to it.
      fl.LogAppendPanic(content1)
    }
