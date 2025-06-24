package context

import (
	"google.golang.org/grpc/status"

	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
)

type errImpl interface {
	Error(code int32, err error) errors.Error
	Errorf(code int32, args ...interface{}) errors.Error
	NewError(code int32, string string) errors.Error
	LegacyWrap(err error) errors.Error
	LegacyWrapCode(code int32, err error) errors.Error
	FromStatus(s *status.Status) errors.Error
}

type err struct {
}

// Errorf 格式化错误, err 不会出现在错误信息中， args 是错误码对应format的参数
func (c *err) Errorf(code int32, args ...interface{}) errors.Error {
	err := error(nil)
	for _, arg := range args {
		if e, ok := arg.(error); ok {
			err = e
		}

	}
	return errors.New(err, code, args...)
}

// Errorf 格式化错误，err 不会出现在错误信息中，
func (c *err) Error(code int32, err error) errors.Error {
	return errors.New(err, code)
}

// NewError 错误透传，
func (c *err) NewError(code int32, message string) errors.Error {
	return errors.NewError(code, message)
}

// LegacyWrap historical legacy wrap
// 如果底层已经是 errors.Error或者ecode.Code 直接用
func (c *err) LegacyWrap(err error) errors.Error {
	if meErr, ok := err.(errors.Error); ok {
		return meErr
	}

	if gcode, ok := err.(eCodeStatusI); ok {
		s := gcode.GetStatus()
		return errors.NewError(int32(s.Code()), s.Message())
	}
	if baseErr, ok := err.(BaseErrI); ok {
		return errors.NewError(baseErr.Code(), baseErr.Message())
	}

	return errors.New(err, code.RawErrWrap, err)
}

type eCodeStatusI interface {
	GetStatus() *status.Status
}

// FromStatus From Status
func (c *err) FromStatus(s *status.Status) errors.Error {
	return errors.NewError(int32(s.Code()), s.Message())
}

// LegacyWrapCode historical legacy wrap
// 如果底层已经是 errors.Error或者ecode.Code 直接用
func (c *err) LegacyWrapCode(code int32, err error) errors.Error {
	if meErr, ok := err.(errors.Error); ok {
		return meErr
	}

	return errors.New(err, code, err)
}

type BaseErrI interface {
	Code() int32
	Message() string
}
