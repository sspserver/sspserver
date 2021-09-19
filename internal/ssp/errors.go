package ssp

// SourceError contains only errors from some source drivers
type SourceError struct {
	Source  interface{}
	Message string
}

func (e *SourceError) Error() string {
	return e.Message
}
