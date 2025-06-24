package context

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

type Log interface {
	Errorf(format string, args ...interface{})
	Error(message string)
	ErrorJSON(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Info(message string)
	InfoJSON(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(message string)
	DebugJSON(format string, args ...interface{})
	Panic(message string)
	Panicf(format string, args ...interface{})
	PanicJSON(format string, args ...interface{})

	LogAndReturnErr(format string, args ...interface{}) error

	SubLog(suffix string) Log
}

type log struct {
	ctx context.Context
}

func NewSubLogCtx(ctx context.Context, suffix string) context.Context {
	requestID := GetLogID(ctx)
	subCtx := context.WithValue(ctx, CtxLogIDKey, requestID+":"+suffix)
	return subCtx
}

func NewLog(ctx context.Context) Log {
	if coreCtx, ok := ctx.(Context); ok {
		return coreCtx.Log()
	}
	return &log{
		ctx: adjustCtxLogID(ctx),
	}
}

func (l *log) SubLog(suffix string) Log {
	ctx := NewSubLogCtx(l.ctx, suffix)
	return &log{
		ctx: ctx,
	}
}

func GetLogID(ctx context.Context) string {
	if ctx == nil {
		ctx = context.TODO()
	}

	requestID, _ := ctx.Value(CtxLogIDKey).(string)

	return requestID
}

func NewLogID() string {
	return "svc:" + uuid.New().String()
}

// AdjustCtxLogID 适配http,grpc，没有请求 等情况
func AdjustCtxLogID(ctx context.Context) context.Context {
	return adjustCtxLogID(ctx)
}

// adjustCtxLogID 适配http,grpc，没有请求 等情况
func adjustCtxLogID(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	requestID := GetLogID(ctx)
	if requestID == "" {
		return context.WithValue(ctx, CtxLogIDKey, requestID)
	}

	return ctx

}

func (l *log) Errorf(format string, args ...interface{}) {
	logrus.Errorf(l.logPrefix()+format+"\n", args...)
}
func (l *log) Error(message string) {
	logrus.Error(l.logPrefix() + message + "\n")
}

// ErrorJSON 根据arg 的Error() string, String() string 来输出参数，  会将Struct，Interface,Array,Map, Slice复杂结构默认转换未json
func (l *log) ErrorJSON(format string, args ...interface{}) {
	args = argsJSON(args...)
	logrus.Errorf(l.logPrefix()+format+"\n", args...)
}

func (l *log) LogAndReturnErr(format string, args ...interface{}) error {
	args = argsJSON(args...)
	er := fmt.Errorf(format, args...)
	logrus.Errorf(l.logPrefix() + er.Error() + "\n")
	return er
}

func (l *log) Infof(format string, args ...interface{}) {
	logrus.Infof(l.logPrefix()+format+"\n", args...)
}
func (l *log) Info(message string) {
	logrus.Info(l.logPrefix() + message + "\n")
}

// InfoJSON 根据arg 的Error() string, String() string 来输出参数， 复杂结构默认转换未json
func (l *log) InfoJSON(format string, args ...interface{}) {
	args = argsJSON(args...)
	logrus.Infof(l.logPrefix()+format+"\n", args...)
}

func (l *log) Debugf(format string, args ...interface{}) {
	logrus.Debugf(l.logPrefix()+format+"\n", args...)
}
func (l *log) Debug(message string) {
	logrus.Debug(l.logPrefix() + message + "\n")
}

// DebugJSON 根据arg 的Error() string, String() string 来输出参数，  会将Struct，Interface,Array,Map, Slice复杂结构默认转换未json
func (l *log) DebugJSON(format string, args ...interface{}) {
	args = argsJSON(args...)
	logrus.Debugf(l.logPrefix() + format + "\n")
}

func (l *log) Panic(message string) {
	buf := make([]byte, 1*1024*1024)
	buf = buf[:runtime.Stack(buf, false)]
	logrus.Panic(l.logPrefix() + message + ", panic stack: " + string(buf) + "\n")
}

func (l *log) Panicf(format string, args ...interface{}) {
	buf := make([]byte, 1*1024*1024)
	buf = buf[:runtime.Stack(buf, false)]
	args = append(args, string(buf))
	logrus.Panicf(l.logPrefix()+format+", panic stack: %s\n", args...)
}

// PanicJSON 根据arg 的Error() string, String() string 来输出参数，  会将Struct，Interface,Array,Map, Slice复杂结构默认转换未json
func (l *log) PanicJSON(format string, args ...interface{}) {
	buf := make([]byte, 1*1024*1024)
	buf = buf[:runtime.Stack(buf, false)]
	args = append(args, string(buf))
	args = argsJSON(args...)
	logrus.Panicf(l.logPrefix()+format+", panic stack: %s\n", args...)
}

// argsJSON 根据arg 的Error() string, String() string 来输出参数， 会将Struct，Interface,Array,Map, Slice复杂结构默认转换未json
func argsJSON(args ...interface{}) []interface{} {
	params := []interface{}{}

	for _, arg := range args {
		if f, ok := arg.(errorFunc); ok {
			params = append(params, f.Error())
			continue
		}
		if f, ok := arg.(stringFunc); ok {
			params = append(params, f.String())
			continue
		}

		if arg == nil {
			params = append(params, []byte("null"))
			continue
		}

		kind := reflect.TypeOf(arg).Kind()
		if kind == reflect.Ptr {
			kind = reflect.TypeOf(arg).Elem().Kind()
		}
		if kind == reflect.Struct || kind == reflect.Interface ||
			kind == reflect.Array || kind == reflect.Map || kind == reflect.Slice {
			out, err := json.Marshal(arg)
			if err != nil {
				params = append(params, arg)
			} else {
				params = append(params, string(out))
			}
			continue
		}

		params = append(params, arg)
	}

	return params
}

func (l *log) logPrefix() string {
	return l.logSpanID() + codeFilePath()
}

func (l *log) logSpanID() string {
	return fmt.Sprintf("span_id|%v|", l.ctx.Value(spanIDKey))
}

func codeFilePath() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	separator := string(filepath.Separator)
	filePathList := strings.Split(file, separator)
	pathCnt := len(filePathList)
	tmpLocation := ""
	if pathCnt >= 3 {
		tmpLocation = strings.Join(filePathList[pathCnt-3:], separator)
	} else if pathCnt == 2 {
		tmpLocation = strings.Join(filePathList[pathCnt-2:], separator)
	} else if pathCnt == 1 {
		tmpLocation = filePathList[0]
	}
	return "location|" + fmt.Sprintf("%s:%d|", tmpLocation, line)
}

type errorFunc interface {
	Error() string
}

type stringFunc interface {
	String() string
}

func SetLogger(l *logrus.Logger) {
	if l == nil {
		return
	}
	logger = l
}
