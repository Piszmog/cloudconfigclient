package cloudconfigclient

// NotFoundError is a special error that is used to propagate 404s.
type NotFoundError struct {
}

// Error return the error message.
func (r NotFoundError) Error() string {
	return "failed to find resource"
}

var notFoundErrorType *NotFoundError
var notFoundError = &NotFoundError{}
