package errs

import "errors"

var (
	ErrSubjectNotFound = errors.New("subjectNotFound")
	ErrSubPubClosed    = errors.New("sub pub is closed")
	ErrNoConnection    = errors.New("no connection")
)
