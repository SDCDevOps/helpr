package rstatus

import "net/http"

const (
	// ReturnStatusSuccess - Successful return.
	ReturnStatusSuccess = 1
	// ReturnStatusFailInternal - Failed return: Internal error.
	ReturnStatusFailInternal = -1
	// ReturnStatusFailInvalidInput - Failed return: Invalid input.
	ReturnStatusFailInvalidInput = -2
	// ReturnStatusFailExternalParty - Failed return: External party error.
	ReturnStatusFailExternalParty = -3
	// ReturnStatusFailDB - Failed return: DB error.
	ReturnStatusFailDB = -4
)

// Status - Return status.
type Status struct {
	Type    int
	Message string
	Err     error
}

// New - Instantiate new object.
func New(typ int, msg string, err error) (s Status) {
	s.Type = typ
	s.Message = msg
	s.Err = err

	return
}

// GetHTTPStatus - Get relevant http status.
func (stat *Status) GetHTTPStatus() (httpStatus int) {
	switch stat.Type {
	case ReturnStatusSuccess:
		httpStatus = http.StatusOK
	case ReturnStatusFailInvalidInput:
		httpStatus = http.StatusBadRequest
	default:
		httpStatus = http.StatusInternalServerError
	}

	return
}
