package errors

import (
	"context"
	osErr "errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rentiansheng/dashboard/middleware/errors/register"
)

type Error interface {
	Code() int32
	Message() string
	Error() string
	String() string
	Caller() []string
	SetError(error) Error
	RawErrorString() string
}

type errors struct {
	message string
	error   error
	code    int32
	stack   *stack
}

func (e errors) Code() int32 {
	return int32(e.code)
}

func (e errors) Error() string {
	return e.String()
}

func (e errors) Message() string {
	return e.message
}

func (e *errors) SetError(err error) Error {
	e.error = err
	return e
}

func (e *errors) RawErrorString() string {
	if e.error == nil {
		return ""
	}
	return e.error.Error()
}

func (e errors) String() string {
	message := e.message
	if e.error != nil {
		message += logSplitFlag + "raw_error:" + e.error.Error()
	}
	return message
}

func (e errors) Caller() []string {
	if e.stack == nil {
		return nil
	}
	return e.stack.StackTrace()
}

func (e errors) Unwrap() error {
	return e.error
}

// new 输出error的时候, 同时输出出错error,request id和错误码对应的错误信息
func new(ctx context.Context, err error, code int32, args ...interface{}) Error {
	format := register.Get(defaultLang, code)

	message := fmt.Sprintf(format, args...)
	if err == nil {
		err = fmt.Errorf(message)
	}
	return &errors{
		message: message,
		error:   err,
		code:    code,
		stack:   callers(),
	}
}

// New  输出error的时候, 同时输出出错error和错误码对应的错误信息
func New(err error, code int32, args ...interface{}) Error {
	format := register.Get(defaultLang, code)
	message := fmt.Sprintf(format, args...)
	if err == nil {
		err = fmt.Errorf(message)
	}
	return &errors{
		message: message,
		error:   err,
		code:    code,
		stack:   callers(),
	}
}

// NewError  输出error的时候, 同时输出出错error和错误码对应的错误信息
func NewError(code int32, message string) Error {
	return &errors{
		message: message,
		error:   fmt.Errorf(message),
		code:    code,
		stack:   callers(),
	}
}

func NewErrorOf(code int32, format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return &errors{
		message: message,
		error:   fmt.Errorf(message),
		code:    code,
		stack:   callers(),
	}
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) StackTrace() []string {
	if s == nil {
		return nil
	}

	frames := runtime.CallersFrames(*s)
	results := make([]string, 0, 10)
	for {
		frame, more := frames.Next()
		results = append(results, fmt.Sprintf("%s:%d", s.file(frame.File), frame.Line))
		// Check whether there are more frames to process after this one.
		if !more {
			break
		}

	}
	return results

}

func (s stack) file(file string) string {
	separator := string(filepath.Separator)
	filePathList := strings.Split(file, separator)
	pathCnt := len(filePathList)
	tmpLocation := ""
	if pathCnt >= 4 {
		tmpLocation = strings.Join(filePathList[pathCnt-4:], separator)
	} else if pathCnt == 3 {
		tmpLocation = strings.Join(filePathList[pathCnt-3:], separator)
	} else if pathCnt == 2 {
		tmpLocation = strings.Join(filePathList[pathCnt-2:], separator)
	} else if pathCnt == 1 {
		tmpLocation = filePathList[0]
	}

	return tmpLocation
}

func Is(err, target error) bool {
	return osErr.Is(err, target)
}

const logSplitFlag = "|"
const defaultLang = "default"
