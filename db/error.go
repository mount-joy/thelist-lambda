package db

import "errors"

// ErrorNotFound is the error returned when an item could not be found
var ErrorNotFound = errors.New("Not Found")

// ErrorBadRequest is the error returned when an error didn't satisfy the contract
var ErrorBadRequest = errors.New("Bad Request")

// ErrorIDExists is the error returned when an item could not created because it already exists
var ErrorIDExists = errors.New("ID Already Exists")
