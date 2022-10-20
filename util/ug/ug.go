package ug

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"utilware/tig"
)

const (
	Success = "success" // 成功
	Failure = "failure" // 失败
	Error   = "error"   // 错误
)

var (
	debugMode bool = false // 调试模式
)

// 设置/获取调试模式
// @param mode bool 调试模式
// @return bool 调试模式
func Debug(b ...bool) bool {
	if len(b) > 0 {
		debugMode = b[0]
	}

	return debugMode
}

// 主要数据结构
type Data struct {
	file    string      // 文件名
	line    int         // 行号信息
	code    int         // 错误码
	message string      // 提示信息
	err     error       // 错误信息
	data    interface{} // 返回数据(业务接口定义具体数据结构)
}

// 数据返回通用数据结构
type Response struct {
	Id      string      `json:"id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Raw struct {
	File    string      // 文件名
	Line    int         // 行号信息
	Code    int         // 错误码
	Message string      // 提示信息
	Err     error       // 错误信息
	Data    interface{} // 返回数据(业务接口定义具体数据结构)
}

// 新建返回数据结构
// @param code int 错误码
// @param message string 提示信息
// @param data interface{} 返回数据 (可选参数)
// @return *Data 返回数据结构
func New(code int, message string, data ...interface{}) *Data {
	r := &Data{
		code:    code,
		message: message,
	}

	if len(data) > 0 {
		r.data = data[0]
	}

	if r.message == "" {
		r.message = Success
	}

	_, r.file, r.line, _ = runtime.Caller(1)

	return r
}

// 追加错误
// @param err error 错误信息
// @return *Data 返回数据结构
func (r *Data) Add(e error) *Data {
	r.err = e
	return r
}

// 追加格式化错误
// @param format string 格式化字符串
// @param args ...interface{} 格式化参数
// @return *Data 返回数据结构
func (r *Data) Addf(format string, args ...interface{}) *Data {
	r.err = fmt.Errorf(format, args...)
	return r
}

// 获取错误
// @return error 错误信息
func (r *Data) Raw() *Raw {
	return &Raw{
		File:    r.file,
		Line:    r.line,
		Code:    r.code,
		Message: r.message,
		Err:     r.err,
		Data:    r.data,
	}
}

// recover 错误处理
// @return *Data 返回数据结构
func (r *Data) Recover() *Data {
	if rec := recover(); rec != nil {
		r.err = fmt.Errorf("%v", rec)
	}

	return r
}

// 错误消息唯一标识 (由文件名、行号、错误码、错误信息生成)
func (r *Data) Id() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%d:%s:%d:%s", r.code, r.file, r.line, r.message)))

	return hex.EncodeToString(h.Sum(nil))
}

// 实现 error 接口
// @return string 错误信息字符串
func (r *Data) Error() string {
	if r.err != nil && debugMode {
		return fmt.Sprintf("[%d:%s] [%s:%d] %s", r.code, r.message, r.file, r.line, r.err.Error())
	}

	if r.code != 0 {
		return fmt.Sprintf("[%d] %s", r.code, r.message)
	}

	return r.message
}

// 输出信息字符串 (不包括data数据)
// @return string 信息字符串
func (r *Data) Output() string {
	if r.err != nil {
		return fmt.Sprintf("[%d:%s] [%s:%d] %s", r.code, r.message, r.file, r.line, r.err.Error())
	}

	return fmt.Sprintf("[%d] [%s:%d] %s", r.code, r.file, r.line, r.message)
}

// 输出返回结果数据结构
// @return *Response 返回结果数据结构
func (r *Data) Response() *Response {
	return &Response{
		Id:      r.Id(),
		Code:    r.code,
		Message: r.Error(),
		Data:    r.data,
	}
}

// 标准返回结果数据结构封装
// @param c *tig.Context tig web 框架上下文
func (r *Data) Json(c *tig.Context) {
	c.JSON(200, r.Response())
}

// 返回JSON数据并退出当前HTTP执行函数
// @param c *tig.Context tig web 框架上下文
func (r *Data) JsonExit(c *tig.Context) {
	r.Json(c)
	c.End()
}

// 标准返回结果数据结构封装
// @param c *tig.Context tig web 框架上下文
func (r *Data) Text(c *tig.Context) {
	c.Sprintf(r.code, "%s - %s", r.Id(), r.Error())
}

// 返回Text数据并退出当前HTTP执行函数
// @param c *tig.Context tig web 框架上下文
func (r *Data) TextExit(c *tig.Context) {
	r.Text(c)
	c.End()
}
