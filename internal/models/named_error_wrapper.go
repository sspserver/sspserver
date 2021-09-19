package models

import "fmt"

// NamedErrorWrapper type
type NamedErrorWrapper struct {
	name string
	err  error
}

// NewNamedErrorWrapper returns new error object
func NewNamedErrorWrapper(name string, err error) NamedErrorWrapper {
	return NamedErrorWrapper{
		name: name,
		err:  err,
	}
}

func (err NamedErrorWrapper) Error() string {
	return fmt.Sprintf("[%s] %s", err.name, err.err.Error())
}

// GetName of the error
func (err NamedErrorWrapper) GetName() string {
	return err.name
}

// OriginalError of the wrapper
func (err NamedErrorWrapper) OriginalError() error {
	return err.err
}
