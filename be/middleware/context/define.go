package context

import (
	"bytes"
	osCtx "context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	gorestful "github.com/emicklei/go-restful/v3"

	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/mapper"
)

const (
	spanIDKey = "span-id"
	// CtxLogIDKey context log id key
	CtxLogIDKey = "trace-id"
)

type Contexts interface {
	Response() http.ResponseWriter

	// JSONDecode 使用json格式解析request body
	// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
	JSONDecode(target interface{}) errors.Error
	// Decode 根据http中的信息适配需要解析方式。已经支持自动解析query string, post(application/x-www-form-urlencoded,json)
	// 默认是先解析 query string ,然后根据 http content type 解析 合并数据。默认是json 方式
	// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
	Decode(target interface{}) errors.Error
	FromFile(name string) (multipart.File, *multipart.FileHeader, error)
	Query(name string) []string
	PathParameter(name string) string
	PathParameters() map[string]string
	SetExtraResponse(key string, val interface{})
	// SetPageResponse  这里的数据回与data 同级返回给用户page字段信息，只有用contexts 返回有效
	SetPageResponse(val interface{})
	SetResponseFile(fileName string, content *bytes.Buffer)
	GetResponseFile() (fileName string, content *bytes.Buffer, exists bool)
	GetExtraResponse() map[string]interface{}

	SetRawResponse(typ string, body []byte)
	GetRawResponse() (typ string, body []byte, exists bool)

	SelectedRoutePath() string

	SetData(data interface{})
	GetData() interface{}

	HTTPBody() ([]byte, error)

	Context
}

type Context interface {
	osCtx.Context
	Cookie() []*http.Cookie
	Header() http.Header
	Request() *http.Request

	GetRequestID() string

	Log() Log
	Error() errImpl

	// IsDone 判断context 是否已经结束
	IsDone() bool
	// WithTimeout 设置超时间时间
	WithTimeout(timeout time.Duration)

	WithValue(key, value interface{})
	SubContext(suffix string) Context
	WithTimeoutCtx(timeout time.Duration) Context

	Mapper(action string, src, dst interface{}) errors.Error
	AllMapper(action string, src, dst interface{}) errors.Error

	Cancel() osCtx.CancelFunc

	WithSpan() Context
	WithSpanPrefix(prefix string) Context
	WithSpanID(id string) Context
}

type contexts struct {
	requestID string
	spanID    string
	ctx       osCtx.Context
	req       *gorestful.Request
	resp      *gorestful.Response
	cancelFn  osCtx.CancelFunc
	meErr     errImpl

	isResponseFile bool
	fileName       string
	fileContent    *bytes.Buffer
	extraData      map[string]interface{}

	rawRes *rawResponse

	// response $.data
	data interface{}
}

type rawResponse struct {
	body []byte
	typ  string
}

// NewContext 不能向 Contexts 对象转换，回panic
func NewContext(ctx osCtx.Context, r *gorestful.Request, w *gorestful.Response) Contexts {
	if ctx == nil {
		ctx = osCtx.Background()
	}
	if rctx, ok := ctx.(*contexts); ok {
		newCtx := rctx.clone()
		return newCtx
	} else {
		ctx = AdjustCtxLogID(ctx)
	}

	me := &contexts{
		requestID: GetLogID(ctx),
		ctx:       ctx,
		meErr:     &err{},
		req:       r,
		resp:      w,
		extraData: make(map[string]interface{}, 0),
	}
	return me
}

func (c contexts) Cookie() []*http.Cookie {
	return c.req.Request.Cookies()
}

func (c *contexts) Header() http.Header {
	if c.req == nil || c.req.Request == nil {
		return http.Header{}
	}
	return c.req.Request.Header
}

func (c *contexts) AddResponseHeader(key, val string) {
	c.resp.AddHeader(key, val)
}

func (c *contexts) Request() *http.Request {
	return c.req.Request
}

func (c *contexts) SelectedRoutePath() string {
	if c.req == nil {
		return ""
	}
	return c.req.SelectedRoutePath()
}

func (c *contexts) Response() http.ResponseWriter {
	return c.resp.ResponseWriter
}

func (c *contexts) GetRequestID() string {
	return c.requestID
}

// Errorf 格式化错误, err 不会出现在错误信息中， args 是错误码对应format的参数
func (c *contexts) Errorf(code int32, err error, args ...interface{}) errors.Error {
	return errors.New(err, code, args...)
}

func (c *contexts) Error() errImpl {
	return c.meErr
}

func (c *contexts) Log() Log {
	return &log{ctx: c}
}

// JSONDecode 使用json格式解析request body
// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
func (c *contexts) JSONDecode(target interface{}) errors.Error {

	if err := decodeJSON(c.Request().Body, target); err != nil {
		return c.meErr.LegacyWrapCode(code.JSONDecodeErrCode, err)
	}

	if valid, ok := target.(ValidateI); ok {
		if err := valid.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Decode 根据http中的信息适配需要解析方式。已经支持自动解析query string, post(application/x-www-form-urlencoded,json)
// 默认是先解析 query string ,然后根据 http content type 解析 合并数据。默认是json 方式
// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
func (c *contexts) Decode(target interface{}) errors.Error {
	pathParameters := c.PathParameters()
	urlParams := make(map[string][]string, 0)
	for k, v := range pathParameters {
		urlParams[k] = []string{v}
	}

	if err := autoDecode(c.Request(), urlParams, target); err != nil {
		return c.meErr.LegacyWrapCode(code.JSONDecodeErrCode, err)
	}

	if valid, ok := target.(ValidateI); ok {

		if err := valid.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *contexts) FromFile(name string) (multipart.File, *multipart.FileHeader, error) {
	if c.req != nil && c.req.Request != nil {
		return c.req.Request.FormFile(name)
	}
	return nil, nil, c.meErr.Errorf(code.FileNotFoundErrCode, "file not found")
}

func (c *contexts) Query(name string) []string {
	if c.req == nil || c.req.Request == nil {
		return []string{}
	}
	return c.req.QueryParameters(name)
}

func (c *contexts) PathParameter(name string) string {
	if c.req != nil {
		return c.req.PathParameter(name)
	}
	return ""
}

func (c *contexts) PathParameters() map[string]string {
	if c.req != nil {
		return c.req.PathParameters()
	}
	return map[string]string{}
}

// WithTimeout 设置超时间时间
func (c *contexts) WithTimeout(timeout time.Duration) {
	c.ctx, c.cancelFn = osCtx.WithTimeout(c.ctx, timeout)
	return
}

// IsDone 判断context 是否已经结束
func (c *contexts) IsDone() bool {
	select {
	case <-c.Done():
		return true
	default:
		return false
	}
}

// SetExtraResponse 这里的数据回与data 同级返回给用户，只有用contexts 返回有效
func (c *contexts) SetExtraResponse(key string, val interface{}) {
	c.extraData[key] = val
}

// SetPageResponse  这里的数据回与data 同级返回给用户page字段信息，只有用contexts 返回有效
func (c *contexts) SetPageResponse(val interface{}) {
	c.extraData["page"] = val
}

func (c *contexts) GetExtraResponse() map[string]interface{} {
	return c.extraData
}

func (c *contexts) SetResponseFile(fileName string, content *bytes.Buffer) {
	c.isResponseFile = true
	c.fileContent = content
	c.fileName = fileName
}

func (c *contexts) GetResponseFile() (fileName string, content *bytes.Buffer, exists bool) {
	return c.fileName, c.fileContent, c.isResponseFile
}

func (c *contexts) WithValue(key, value interface{}) {
	c.ctx = osCtx.WithValue(c.ctx, key, value)
}

func (c *contexts) WithTimeoutCtx(timeout time.Duration) Context {
	newCtx := c.clone()
	newCtx.ctx, newCtx.cancelFn = osCtx.WithTimeout(newCtx.ctx, timeout)
	return newCtx
}

func (c *contexts) clone() *contexts {
	ctx := osCtx.WithValue(c.ctx, "_", c.Value("_"))
	return &contexts{
		ctx:            ctx,
		requestID:      c.requestID,
		req:            c.req,
		resp:           c.resp,
		meErr:          c.meErr,
		isResponseFile: c.isResponseFile,
		fileName:       c.fileName,
		fileContent:    c.fileContent,
		extraData:      c.extraData,
	}
}

func (c *contexts) SubContext(suffix string) Context {
	ctx := NewSubLogCtx(c.ctx, suffix)
	newCtx := c.clone()
	newCtx.ctx = ctx
	newCtx.requestID = GetLogID(ctx)
	return newCtx
}

// Mapper automatic data
func (c *contexts) Mapper(action string, src, dst interface{}) errors.Error {
	if src == nil {
		return nil
	}
	if err := mapper.Mapper(c, src, dst); err != nil {
		c.Log().ErrorJSON("action: %s, src: %s, err: %s", action, src, err.Error())
		return c.Error().Errorf(code.MapperActionErrCode, action, err.Error())
	}
	return nil
}

// AllMapper automatic data with struct private field
func (c *contexts) AllMapper(action string, src, dst interface{}) errors.Error {
	if src == nil {
		return nil
	}
	if err := mapper.AllMapper(c, src, dst); err != nil {
		c.Log().ErrorJSON("action: %s, src: %s, err: %s", action, src, err.Error())
		return c.Error().Errorf(code.MapperActionErrCode, action, err.Error())
	}

	return nil
}

func (c *contexts) Cancel() osCtx.CancelFunc {
	if c.cancelFn == nil {
		c.ctx, c.cancelFn = osCtx.WithCancel(c.ctx)
	}
	return c.cancelFn
}

func (c *contexts) WithSpan() Context {
	newCtx := c.clone()
	newCtx.spanID = NewLogID()
	newCtx.ctx = osCtx.WithValue(newCtx.ctx, spanIDKey, newCtx.spanID)
	return newCtx
}

func (c *contexts) WithSpanPrefix(prefix string) Context {
	newCtx := c.clone()
	newCtx.spanID = prefix + "-" + NewLogID()
	newCtx.ctx = osCtx.WithValue(newCtx.ctx, spanIDKey, newCtx.spanID)
	return newCtx
}

func (c *contexts) WithSpanID(id string) Context {
	newCtx := c.clone()
	newCtx.spanID = id
	newCtx.ctx = osCtx.WithValue(newCtx.ctx, spanIDKey, id)
	return newCtx
}

func (c *contexts) SetData(data interface{}) {
	c.data = data
}

func (c contexts) GetData() interface{} {
	return c.data
}

func (c *contexts) HTTPBody() ([]byte, error) {
	if c.req == nil || c.req.Request == nil {
		return nil, nil
	}

	if c.req.Request.Body == nil {
		return nil, nil
	}

	bodyBytes, err := ioutil.ReadAll(c.req.Request.Body)
	if err != nil {
		c.Log().Errorf("io read request.Body fail, %+v", err)
		return nil, c.Error().LegacyWrapCode(code.JSONDecodeErrCode, err)
	}
	c.req.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	return bodyBytes, nil
}

func (c *contexts) SetRawResponse(typ string, body []byte) {
	c.rawRes = &rawResponse{
		body: body,
		typ:  typ,
	}
}

func (c *contexts) GetRawResponse() (typ string, body []byte, exists bool) {
	if c.rawRes == nil {
		return "", nil, false
	}
	return c.rawRes.typ, c.rawRes.body, true
}

func TODO() Context {

	ctx := osCtx.Background()
	ctx = AdjustCtxLogID(ctx)

	me := &contexts{
		requestID: GetLogID(ctx),
		ctx:       ctx,
		meErr:     &err{},
	}

	return me
}
