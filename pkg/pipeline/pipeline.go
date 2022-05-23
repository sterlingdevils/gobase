package pipeline

type Pipelineable[T any] interface {
	PipelineChan() chan T
	Close()
}

// Checkerror will check an error value and return nil if error is not nil
func Checkerror[T any](val *T, e error) *T {
	if e != nil {
		return nil
	}
	return val
}
