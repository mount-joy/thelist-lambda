package db

// Error gives an error including a type
type Error struct {
	error
	ErrorType string
}

// Error Types
const (
	// ErrorNotFound is the error returned when an item could not be found
	ErrorNotFound = "Not Found"
)

// NewError creates a new db.Error instance
func NewError(errorType string) Error {
	return Error{
		ErrorType: errorType,
	}
}
