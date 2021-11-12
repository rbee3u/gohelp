package status

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rbee3u/gohelp/epkg"
)

func Error(code int, message string) error {
	return &baseError{code: code, message: message}
}

func Errorf(code int, format string, a ...interface{}) error {
	return &baseError{code: code, message: fmt.Sprintf(format, a...)}
}

func ErrorWL(code int, message string) error {
	return epkg.FileWithSkip(Error(code, message), 1)
}

func ErrorfWL(code int, format string, a ...interface{}) error {
	return epkg.FileWithSkip(Errorf(code, format, a...), 1)
}

type baseError struct {
	code    int
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

func Wrap(err error, code int, message string) error {
	if err == nil {
		return nil
	}

	return &wrapError{err: err, code: code, message: message}
}

func Wrapf(err error, code int, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapError{err: err, code: code, message: fmt.Sprintf(format, a...)}
}

func WrapWL(err error, code int, message string) error {
	if err == nil {
		return nil
	}

	return epkg.FileWithSkip(Wrap(err, code, message), 1)
}

func WrapfWL(err error, code int, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return epkg.FileWithSkip(Wrapf(err, code, format, a...), 1)
}

type wrapError struct {
	err     error
	code    int
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

func GetCode(err error) int {
	var be *baseError
	if errors.As(err, &be) {
		return be.code
	}

	var we *wrapError
	if errors.As(err, &we) {
		return we.code
	}

	return http.StatusInternalServerError
}
