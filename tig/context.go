package tig

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"utilware/dep/cast"
)

type HandlerFunc func(*Context)

type HandlerWriter struct {
	exported bool

	statusCode int
	header     http.Header
	bodyBuffer *bytes.Buffer
}

type Context struct {
	t *Tig

	// 虚拟的writer
	writer *HandlerWriter
	// 外部调用Write接口
	Writer http.ResponseWriter
	// 原始的writer
	rawWriter http.ResponseWriter
	Req       *http.Request

	middlewareCache *sync.Map

	Path   string
	Method string
	Params map[string]string

	index int

	handlerList []HandlerFunc
}

func newWriterContext(w http.ResponseWriter) *HandlerWriter {
	return &HandlerWriter{
		exported: false,

		statusCode: 200,
		header:     w.Header().Clone(),
		bodyBuffer: bytes.NewBuffer(nil),
	}
}

func (w *HandlerWriter) Header() http.Header {
	return w.header
}

func (w *HandlerWriter) Write(b []byte) (int, error) {
	return w.bodyBuffer.Write(b)
}

func (w *HandlerWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (t *Tig) newContext(w http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		t: t,

		rawWriter: w,
		Req:       req,

		middlewareCache: &sync.Map{},

		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}

	c.writer = newWriterContext(w)
	c.Writer = c.writer
	return c
}

func (c *Context) Param(key string) string {
	if v, ok := c.Params[key]; ok {
		return v
	}

	return ""
}

// URL Query
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// POST Form Value
func (c *Context) PostForm(key string) string {
	c.Req.ParseForm()
	return c.Req.FormValue(key)
}

// 获取当前请求的Cookit
func (c *Context) Cookie(key string) string {
	if v, _ := c.Req.Cookie(key); v != nil {
		return v.Value
	}

	return ""
}

func (c *Context) ReqBody() (body []byte) {
	if c.Req.Body != nil {
		body, _ = io.ReadAll(c.Req.Body)
	}
	c.Req.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}

// 设置状态码
// 也可以获取当前设置的状态码
// 加了魔法, 可以重复设置状态码
func (c *Context) Status(code ...int) int {
	if len(code) > 0 {
		c.Writer.WriteHeader(code[0])
	}

	return c.writer.statusCode
}

// 设置header的简单封装
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置cookie的简单封装
func (c *Context) SetCookie(key, value, path string, httpOnly bool, expires int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		HttpOnly: httpOnly,
		MaxAge:   expires,
		Expires:  time.Now().Add(time.Duration(expires) * time.Second),
	})
}

// 占位符填充输出
// 内部封装了fmt.Sprintf
// 默认content-type为text/html
func (c *Context) Sprintf(code int, format string, values ...any) {
	c.Status(code)

	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 输出字符串
// 默认content-type为text/html
func (c *Context) String(code int, format string) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	c.Writer.Write([]byte(format))
}

func (c *Context) JSON(code int, obj any) {
	c.SetHeader(HeaderContentType, "application/json; charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) XML(code int, obj any) {
	c.SetHeader(HeaderContentType, "application/xml; charset=utf-8")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 输出字节数据
// 默认content-type为application/octet-stream
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "application/octet-stream")
	}

	c.Writer.Write(data)
}

func (c *Context) File(code int, ffs fs.FS, filename string) {
	c.Status(code)
	if c.Writer.Header().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, "text/html;  charset=utf-8")
	}

	b, e := fs.ReadFile(ffs, filename)
	if e != nil {
		http.Error(c.Writer, e.Error(), 500)
		return
	}

	c.Writer.Write(b)
}

func (c *Context) WriteTo(w io.Writer) (int64, error) {
	return c.writer.bodyBuffer.WriteTo(w)
}

// Hijacker 接口封装
// 用于将writer转换为bufio.ReadWriter
// 一旦调用此函数, 其他输出函数将无效
// 并且此函数关闭时将调用End方法结束请求
func (c *Context) Hijacker(f func(bufrw *bufio.ReadWriter) error) error {
	hj, ok := c.rawWriter.(http.Hijacker)
	if !ok {
		return errors.New("http.ResponseWriter is not http.Hijacker")
	}

	c.writer.exported = true

	conn, bufrw, e := hj.Hijack()
	if e != nil {
		return e
	}

	defer conn.Close()

	if e := f(bufrw); e != nil {
		return e
	}

	c.End()
	return nil
}

func (c *Context) GetRawWriter() (http.ResponseWriter, error) {
	if c.writer.exported {
		return nil, errors.New("writer is exported")
	}

	c.writer.exported = true
	return c.rawWriter, nil
}

func (c *Context) WriterExported() bool {
	return c.writer.exported
}

// 循环执行下一个HandlerFunc
func (c *Context) Next() {
	for c.index++; c.index >= 0 && c.index < len(c.handlerList); c.index++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic(r)
				}

				if c.Index() < -1 {
					panic(nil)
				}
			}()

			c.handlerList[c.index](c)
		}()
	}
}

// handler调用深度
func (c *Context) Index() int {
	return c.index
}

// 清除缓冲区
func (c *Context) Clear() {
	c.writer = newWriterContext(c.rawWriter)
	c.Writer = c.writer
}

// 结束请求, 并跳过后续的HandlerFunc
// 将正常输出缓冲区内容
func (c *Context) End() {
	c.index = len(c.handlerList)
	panic(nil)
}

// 重置请求, 并跳过后续的HandlerFunc
// 将无任何输出, 并且浏览器显示连接已重置
// 但是仍然有响应头 "HTTP 1.1 400 Bad Request\r\nConnection: close"
func (c *Context) Close() {
	c.index = -100
	panic(nil)
}

func (c *Context) SetStore(key string, val interface{}) {
	c.middlewareCache.Store(key, val)
}

func (c *Context) GetStore(key string) (interface{}, bool) {
	return c.middlewareCache.Load(key)
}

func (c *Context) DeleteStore(key string) {
	c.middlewareCache.Delete(key)
}

func (c *Context) MustStruct(value any, tagForm string) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
		}
	}()

	// 解析From
	c.Req.ParseForm()
	// 数据暂存
	values := make(map[string]any)

	// 提取From数据
	if len(c.Req.Form) > 0 {
		for k, v := range c.Req.Form {
			values[k] = strings.Join(v, " ")
		}
	}

	// 提取Body数据
	reqBody := bytes.TrimSpace(c.ReqBody())
	if len(reqBody) > 0 {
		if reqBody[0] == '{' && reqBody[len(reqBody)-1] == '}' {
			e = json.Unmarshal(reqBody, &values)
			if e != nil {
				return e
			}
		}
	}

	t := reflect.TypeOf(value).Elem()
	val := reflect.ValueOf(value).Elem()
	for k := 0; k < t.NumField(); k++ {
		// 获取标签名称
		tagName := t.Field(k).Tag.Get(tagForm)

		// 搞事情开始
		if v, ok := values[tagName]; ok && tagName != "" {
			switch val.Field(k).Kind() {
			case reflect.String:
				v = cast.ToString(v)
			case reflect.Int64:
				v = cast.ToInt64(v)
			case reflect.Int32:
				v = cast.ToInt32(v)
			case reflect.Int16:
				v = cast.ToInt16(v)
			case reflect.Int8:
				v = cast.ToInt8(v)
			case reflect.Int:
				v = cast.ToInt(v)
			case reflect.Float64:
				v = cast.ToFloat64(v)
			case reflect.Float32:
				v = cast.ToFloat32(v)
			case reflect.Bool:
				v = cast.ToBool(v)
			case reflect.Slice:
				switch val.Field(k).Type().Elem().Kind() {
				case reflect.String:
					v = cast.ToStringSlice(v)
				case reflect.Int:
					v = cast.ToIntSlice(v)
				case reflect.Bool:
					v = cast.ToBoolSlice(v)
				default:
					continue
				}
			default:
				continue
			}

			val.Field(k).Set(reflect.ValueOf(v))
		}
	}

	return nil
}
