package gresty

// go http client support get,post,delete,patch,put,head,file method.
// go-resty/resty: https:// github.com/go-resty/resty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// 默认请求超时
	defaultTimeout = 3 * time.Second

	// 默认最大重试次数
	defaultMaxRetries = 3

	// resp is nil
	respEmpty = errors.New("resp is empty")
)

// Service 请求句柄设置
type Service struct {
	BaseUri         string        // 请求地址url的前缀
	Timeout         time.Duration // 请求超时限制
	Proxy           string        // 请求设置的http_proxy代理
	EnableKeepAlive bool          // 是否允许长连接方式请求接口，默认短连接方式
}

// RequestOption 请求参数设置
type RequestOption struct {
	Method string // 请求的方法
	Url    string // 请求url

	RetryCount       int                        // 重试次数
	RetryWaitTime    time.Duration              // 重试间隔,默认100ms
	RetryMaxWaitTime time.Duration              // 重试最大等待间隔,默认2s
	RetryConditions  []resty.RetryConditionFunc // 重试条件，是一个函数切片

	Params  map[string]interface{} // get,delete的Params参数
	Data    map[string]interface{} // post请求form data表单数据
	Headers map[string]interface{} // header头信息

	// cookie参数设置
	Cookies        map[string]interface{} // cookie信息
	CookiePath     string                 // 可选参数
	CookieDomain   string                 // cookie domain可选
	CookieMaxAge   int                    // cookie MaxAge
	CookieHttpOnly bool                   // cookie httpOnly

	// 支持post,put,patch以json格式传递,[]int{1, 2, 3},map[string]string{"a":"b"}格式
	// json支持[],{}数据格式,主要是golang的基本数据类型，就可以
	// 直接调用SetBody方法，自动添加header头"Content-Type":"application/json"
	Json interface{}

	// 支持文件上传的参数
	FileName      string // 文件名称
	FileParamName string // 文件上传的表单file参数名称
}

// Reply 请求后的结果
// statusCode,body,error.
type Reply struct {
	StatusCode int    // http request 返回status code
	Err        error  // 请求过程中，发生的error
	Body       []byte // 返回的body内容
}

// Text 返回Reply.Body文本格式
func (r *Reply) Text() string {
	return string(r.Body)
}

// Json 将响应的结果Reply解析到data
// 对返回的Reply.Body做json反序列化处理
func (r *Reply) Json(data interface{}) error {
	if len(r.Body) > 0 {
		err := json.Unmarshal(r.Body, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// ApiStdRes 标准的api返回格式
type ApiStdRes struct {
	Code    int
	Message string
	Data    interface{}
}

// New 创建一个service实例
func New(opts ...Option) *Service {
	s := &Service{
		Timeout: defaultTimeout,
	}

	s.apply(opts)

	if s.Timeout == 0 {
		s.Timeout = defaultTimeout
	}

	return s
}

// NewRestyClient 创建一个resty client
// 创建一个resty客户端，支持post,get,delete,head,put,patch,file文件上传等
// 可以快速使用go-resty/resty上面的方法
// 参考文档： https:// github.com/go-resty/resty
func (s *Service) NewRestyClient() *resty.Client {
	client := resty.New()

	return client
}

// Do 请求方法
// method string  请求的方法get,post,put,patch,delete,head等
// uri    string  请求的相对地址，如果BaseUri为空，就必须是完整的url地址
// opt 	  *RequestOption 请求参数ReqOpt
// 短连接的形式请求api
// 关于如何关闭http connection
// https:// www.cnblogs.com/cobbliu/p/4517598.html
func (s *Service) Do(method string, reqUrl string, opt *RequestOption) *Reply {
	if method == "" || reqUrl == "" {
		return &Reply{
			Err: errors.New("request Method or request url is empty"),
		}
	}

	client := s.NewRestyClient()
	if opt == nil {
		opt = &RequestOption{}
	}

	opt.Method = method
	opt.Url = reqUrl
	return s.Request(opt, client)
}

// Request 请求方法
// resty.setBody: for struct and map data type defaults to 'application/json'
// SetBody method sets the request body for the request. It supports various realtime needs as easy.
// We can say its quite handy or powerful. Supported request body data types is `string`,
// `[]byte`, `struct`, `map`, `slice` and `io.Reader`. Body value can be pointer or non-pointer.
// Automatic marshalling for JSON and XML content type, if it is `struct`, `map`, or `slice`.
//
// client.R().SetFormData method sets Form parameters and their values in the current request.
// It's applicable only HTTP method `POST` and `PUT` and requests content type would be
// set as `application/x-www-form-urlencoded`.
func (s *Service) Request(reqOpt *RequestOption, client *resty.Client) *Reply {
	if client == nil {
		client = s.NewRestyClient()
	}

	if s.BaseUri != "" {
		reqOpt.Url = strings.TrimRight(s.BaseUri, "/") + "/" + reqOpt.Url
	}

	if s.Timeout == 0 {
		s.Timeout = defaultTimeout
	}

	client = client.SetTimeout(s.Timeout) // timeout设置
	if !s.EnableKeepAlive {
		client = client.SetHeader("Connection", "close") // 显示指定短连接
	}

	if s.Proxy != "" {
		client = client.SetProxy(s.Proxy)
	}

	// 重试次数，重试间隔，最大重试超时时间
	if reqOpt.RetryCount > 0 {
		if reqOpt.RetryCount >= defaultMaxRetries {
			reqOpt.RetryCount = defaultMaxRetries // 最大重试次数
		}

		if len(reqOpt.RetryConditions) > 0 {
			client.RetryConditions = reqOpt.RetryConditions
		}

		client = client.SetRetryCount(reqOpt.RetryCount)
		if reqOpt.RetryWaitTime != 0 {
			client = client.SetRetryWaitTime(reqOpt.RetryWaitTime)
		}

		if reqOpt.RetryMaxWaitTime != 0 {
			client = client.SetRetryMaxWaitTime(reqOpt.RetryMaxWaitTime)
		}
	}

	if cLen := len(reqOpt.Cookies); cLen > 0 {
		cookies := make([]*http.Cookie, cLen)
		for k := range reqOpt.Cookies {
			cookies = append(cookies, &http.Cookie{
				Name:     k,
				Value:    fmt.Sprintf("%v", reqOpt.Cookies[k]),
				Path:     reqOpt.CookiePath,
				Domain:   reqOpt.CookieDomain,
				MaxAge:   reqOpt.CookieMaxAge,
				HttpOnly: reqOpt.CookieHttpOnly,
			})
		}

		client = client.SetCookies(cookies)
	}

	// 设置header头
	if len(reqOpt.Headers) > 0 {
		client = client.SetHeaders(s.ParseData(reqOpt.Headers))
	}

	var resp *resty.Response
	var err error
	method := strings.ToLower(reqOpt.Method)
	switch method {
	case "get", "delete", "head":
		client = client.SetQueryParams(s.ParseData(reqOpt.Params))
		if method == "get" {
			resp, err = client.R().Get(reqOpt.Url)
			return s.GetResult(resp, err)
		}

		if method == "delete" {
			resp, err = client.R().Delete(reqOpt.Url)
			return s.GetResult(resp, err)
		}

		if method == "head" {
			resp, err = client.R().Head(reqOpt.Url)
			return s.GetResult(resp, err)
		}
	case "post", "put", "patch":
		request := client.R()
		if len(reqOpt.Data) > 0 {
			request = request.SetFormData(s.ParseData(reqOpt.Data))
		}

		if reqOpt.Json != nil {
			request = request.SetBody(reqOpt.Json)
		}

		if method == "post" {
			resp, err = request.Post(reqOpt.Url)
			return s.GetResult(resp, err)
		}

		if method == "put" {
			resp, err = request.Put(reqOpt.Url)
			return s.GetResult(resp, err)
		}

		if method == "patch" {
			resp, err = request.Patch(reqOpt.Url)
			return s.GetResult(resp, err)
		}
	case "file":
		b, err := ioutil.ReadFile(reqOpt.FileName)
		if err != nil {
			return &Reply{
				Err: errors.New("read file error: " + err.Error()),
			}
		}

		// 文件上传
		resp, err := client.R().
			SetFileReader(reqOpt.FileParamName, reqOpt.FileName, bytes.NewReader(b)).
			Post(reqOpt.Url)
		return s.GetResult(resp, err)
	default:
	}

	return &Reply{
		Err:        errors.New("request method not support"),
		StatusCode: http.StatusServiceUnavailable,
	}
}

// ParseData 解析ReqOpt Params和Data
func (s *Service) ParseData(d map[string]interface{}) map[string]string {
	dLen := len(d)
	if dLen == 0 {
		return nil
	}

	// 对d参数进行处理
	data := make(map[string]string, dLen)
	for k, v := range d {
		if val, ok := v.(string); ok {
			data[k] = val
		} else {
			data[k] = fmt.Sprintf("%v", v)
		}
	}

	return data
}

// GetResult 处理请求的结果statusCode,body,error.
// 首先判断是否出错，然后判断http resp是否请求成功或有错误产生
func (s *Service) GetResult(resp *resty.Response, err error) *Reply {
	res := &Reply{}
	if err != nil {
		if resp != nil {
			res.StatusCode = resp.StatusCode()
			res.Body = resp.Body()
		}

		res.Err = err
		return res
	}

	if resp == nil {
		res.StatusCode = http.StatusServiceUnavailable
		res.Err = respEmpty
		return res
	}

	res.Body = resp.Body()
	res.StatusCode = resp.StatusCode()
	if !resp.IsSuccess() || resp.IsError() {
		res.Err = fmt.Errorf("resp error: %v", resp.Error())
		return res
	}

	return res
}
