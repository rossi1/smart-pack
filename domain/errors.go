package domain

import "errors"

const (
	BadRequestStatus          = 400
	notFoundStatus            = 404
	conflictStatus            = 409
	UnprocessableEntity       = 422
	InternalServerErrorStatus = 500
)

type CustomError struct {
	l string
	e string
	s int
}

func NewCustomError(l, e string, s int) *CustomError {
	return &CustomError{l: l, e: e, s: s}
}

func (x CustomError) Error() string {
	return x.e
}

func (x *CustomError) Status() int {
	return x.s
}

func (x *CustomError) Label() string {
	return x.l
}

func IsHTTPCustomError(err error) (*CustomError, bool) {
	var cerr *CustomError
	ok := errors.As(err, &cerr)
	return cerr, ok
}

type ErrorReporter interface {
	ReportError(err error)
}
