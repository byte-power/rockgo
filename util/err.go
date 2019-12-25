package util

import "strings"

// Error may wrap another error
type Error interface {
	error
	Unwrap() error
}

func NewError(message string, wrap error) Error {
	return &errorImpl{message: message, wrapped: wrap}
}

type errorImpl struct {
	message string
	wrapped error
}

func (e *errorImpl) Error() string {
	var m strings.Builder
	m.WriteString(e.message)
	if e.wrapped != nil {
		m.WriteByte('<')
		m.WriteString(e.wrapped.Error())
		m.WriteByte('>')
	}
	return m.String()
}

func (e *errorImpl) Unwrap() error {
	return e.wrapped
}
