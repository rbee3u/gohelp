package epkg

import (
	"fmt"
	"io"
	"runtime"
)

func Error(message string) error {
	return &baseError{message: message}
}

func Errorf(format string, a ...interface{}) error {
	return &baseError{message: fmt.Sprintf(format, a...)}
}

func ErrorWL(message string) error {
	return FileWithSkip(Error(message), 1)
}

func ErrorfWL(format string, a ...interface{}) error {
	return FileWithSkip(Errorf(format, a...), 1)
}

type baseError struct {
	message string
}

func (e *baseError) Error() string {
	return e.message
}

func (e *baseError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		_, _ = io.WriteString(s, e.message)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.message)
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &wrapError{err: err, message: message}
}

func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapError{err: err, message: fmt.Sprintf(format, a...)}
}

func WrapWL(err error, message string) error {
	if err == nil {
		return nil
	}

	return FileWithSkip(Wrap(err, message), 1)
}

func WrapfWL(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return FileWithSkip(Wrapf(err, format, a...), 1)
}

type wrapError struct {
	err     error
	message string
}

func (e *wrapError) Error() string {
	return e.message + ": " + e.err.Error()
}

func (e *wrapError) Unwrap() error {
	return e.err
}

func (e *wrapError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%s: %+v", e.message, e.err)

			return
		}

		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func File(err error) error {
	return FileWithSkip(err, 1)
}

func FileWithSkip(err error, skip int) error {
	pc, pre, num, _ := runtime.Caller(skip + 1)
	if pc != 0 {
		fn := runtime.FuncForPC(pc).Name()
		if len(fn) != 0 {
			pre = fn
		}
	}

	return &fileError{err: err, pre: pre, num: num}
}

type fileError struct {
	err error
	pre string
	num int
}

func (e *fileError) Error() string {
	return e.err.Error()
}

func (e *fileError) Unwrap() error {
	return e.err
}

func (e *fileError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "(%s:%d)%+v", e.pre, e.num, e.err)

			return
		}

		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}
