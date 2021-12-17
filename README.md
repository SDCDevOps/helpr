# helpr
**helpr** is a _not so useful_ helper package in Go.

It contains the following packages:  
* mgdb  
* rstatus

---
## Installation

Install the dependency in your project.

    go get github.com/skker/helpr

---
## mgdb

**mgdb** provides some simple methods to use MongoDB in Go. 

It is created to save you a few lines of codes. Other than that, it's quite useless.

### Usage

    import "github.com/skker/helpr/mgdb"
    
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

    import "github.com/skker/helpr/rstatus"

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


