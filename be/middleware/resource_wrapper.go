package middleware

import (
	"bytes"

	"fmt"

	"io/ioutil"
	"strings"
	"sync"

	gorestful "github.com/emicklei/go-restful/v3"
	coreContext "github.com/rentiansheng/dashboard/middleware/context"
	appErrors "github.com/rentiansheng/dashboard/middleware/errors"
)

var (
	apiChecker sync.Once
)

const (
	responseHTTHeaderRequestID = "trace-Id"
	RequestId                  = responseHTTHeaderRequestID
)

type Handler func(ctx coreContext.Contexts) appErrors.Error
type HttpJsonResponse struct {
	Retcode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OkResponseExtra(message string, data interface{}, extraRespData map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	for key, val := range extraRespData {
		result[key] = val
	}
	result["retcode"] = 0
	result["message"] = message
	result["data"] = data

	return result

}

func OkResponse(message string, data interface{}) HttpJsonResponse {
	return HttpJsonResponse{
		Retcode: 0,
		Message: message,
		Data:    data,
	}
}

func FailResponse(retcode int, message string, data interface{}) HttpJsonResponse {
	return HttpJsonResponse{
		Retcode: retcode,
		Message: message,
		Data:    data,
	}
}

func Wrapper(h Handler) func(r *gorestful.Request, w *gorestful.Response) {
	return wrapperOptions(func(ctx coreContext.Contexts) appErrors.Error {
		return h(ctx)
	}, DefaultOption())
}

func wrapperOptions(h Handler, o Option) func(r *gorestful.Request, w *gorestful.Response) {
	return func(r *gorestful.Request, w *gorestful.Response) {
		ctx := coreContext.NewContext(r.Request.Context(), r, w)

		requestID := ctx.GetRequestID()
		w.AddHeader(responseHTTHeaderRequestID, requestID)

		// 从panic中恢复
		defer func() {
			if e := recover(); e != nil {
				ctx.Log().Panicf("panic. err: %#v", e)
				_ = w.WriteAsJson(e)
				return
			}
		}()

		// 记录请求body
		requestRecords(ctx)

		err := h(ctx)
		data := ctx.GetData()
		responseRecords(ctx, data, err)
		if err != nil {
			if eerr, ok := err.(appErrors.Error); ok {
				_ = w.WriteAsJson(FailResponse(int(eerr.Code()), eerr.Message(), data))
			} else {
				_ = w.WriteAsJson(FailResponse(-1, err.Error(), data))
			}
		} else {

			if fileName, fileContent, exists := ctx.GetResponseFile(); exists {
				w.Header().Add("Content-Disposition", "attachment; filename="+fileName)
				w.Header().Add("Content-Type", "application/octet-stream")
				w.Header().Add("Content-Length", fmt.Sprintf("%d", fileContent.Len()))
				_, _ = w.Write(fileContent.Bytes())

			} else if typ, body, exists := ctx.GetRawResponse(); exists {

				if typ != "" {
					w.Header().Add("Content-Type", typ)
				}
				w.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
				_, _ = w.Write(body)
			} else {
				extraRespData := ctx.GetExtraResponse()
				if len(extraRespData) > 0 {
					_ = w.WriteAsJson(OkResponseExtra("", data, extraRespData))
				} else {
					_ = w.WriteAsJson(OkResponse("", data))
				}
			}
		}
	}
}

func responseRecords(ctx coreContext.Contexts, data interface{}, err error) {

	log := ctx.Log()
	defer func() {
		if r := recover(); r != nil {
			log.PanicJSON("response record. panic info: %s", r)
		}
	}()
	if err != nil {
		if eerr, ok := err.(appErrors.Error); ok {
			log.ErrorJSON("response record, response error. err code: %d, err message: %s, raw msg: %s, caller: %s", eerr.Code(), eerr.Message(), eerr.RawErrorString(), eerr.Caller())
		} else if eerr, ok := err.(coreContext.BaseErrI); ok {
			log.ErrorJSON("response record, response error. err code: %d, err message: %s", eerr.Code(), eerr.Message())
		} else {
			log.ErrorJSON("response record, response error. err: %s", err)
		}
	} else {
		log.InfoJSON("response record. data: %s", data)
	}
	if ctx.Response() != nil {
		log.InfoJSON("response record http header. header: %s", ctx.Response().Header())
	}
}

// requestRecords 记录请求body < 1m且content_type=application/josn 的http 请求的body
func requestRecords(c coreContext.Context) {
	log := c.Log()

	parentReqId := c.Request().Header.Get(RequestId)

	if c.Request().ContentLength < 1024*1024*1 {

		// Ignore requests smaller than 1MB. This helps prevent delaying
		ct := c.Request().Header.Get("Content-Type")
		if strings.ToLower(ct) == "application/json" || strings.HasPrefix(strings.ToLower(ct), "application/json;") {
			bodyBytes, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				log.Errorf("io read request.Body fail, %+v", err)
			}
			// 去除换行符，避免日志换行
			strBody := strings.Replace(string(bodyBytes), "\n", "", -1)
			log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, body: %s",
				c.Request().Method, c.Request().RequestURI, parentReqId, strBody)
			c.Request().Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		} else {
			log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, body: not support Content-Type=%v",
				c.Request().Method, c.Request().RequestURI, parentReqId, ct)
		}
	} else {
		// body 超过1m 不记录
		log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, request body more than 1MB",
			c.Request().Method, c.Request().RequestURI, parentReqId)
	}
}
