package gresty

//  go http client support get,post,delete,patch,put,head,file method
//  go-resty/resty: https:// github.com/go-resty/resty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// 默认请求超时
var defaultTimeout = 3 * time.Second

var defaultMaxRetries = 3 // 默认最大重试次数

//  Service 请求句柄设置
type Service struct {
	BaseUri         string        // 请求地址url的前缀
	Timeout         time.Duration // 请求超时限制
	Proxy           string        // 请求设置的http_proxy代理
	EnableKeepAlive bool          // 是否允许长连接方式请求接口，默认短连接方式
}

//  ReqOpt 请求参数设置
type ReqOpt struct {
	// 重试机制设置
	RetryCount       int           // 重试次数
	RetryWaitTime    time.Duration // 重试间隔,默认100ms
	RetryMaxWaitTime time.Duration // 重试最大等待间隔,默认2s

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

//  Reply 请求后的结果
type Reply struct {
	Err  error  // 请求过程中，发生的error
	Body []byte // 返回的body内容
}

//  ApiStdRes 标准的api返回格式
type ApiStdRes struct {
	Code    int
	Message string
	Data    interface{}
}

//  ParseData 解析ReqOpt Params和Data
func (ReqOpt) ParseData(d map[string]interface{}) map[string]string {
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

//  Do 请求方法
//  method string  请求的方法get,post,put,patch,delete,head等
//  uri    string  请求的相对地址，如果BaseUri为空，就必须是完整的url地址
//  opt 	  *ReqOpt 请求参数ReqOpt
func (s *Service) Do(method string, reqUrl string, opt *ReqOpt) *Reply {
	if method == "" {
		return &Reply{
			Err: errors.New("request method is empty"),
		}
	}

	if reqUrl == "" {
		return &Reply{
			Err: errors.New("request url is empty"),
		}
	}

	if opt == nil {
		opt = &ReqOpt{}
	}

	if s.BaseUri != "" {
		reqUrl = strings.TrimRight(s.BaseUri, "/") + "/" + reqUrl
	}

	if s.Timeout == 0 {
		s.Timeout = defaultTimeout
	}

	// 短连接的形式请求api
	// 关于如何关闭http connection
	// https:// www.cnblogs.com/cobbliu/p/4517598.html

	// 创建请求客户端
	client := resty.New()
	client = client.SetTimeout(s.Timeout) // timeout设置

	if !s.EnableKeepAlive {
		client = client.SetHeader("Connection", "close") // 显示指定短连接
	}

	if s.Proxy != "" {
		client = client.SetProxy(s.Proxy)
	}

	// 重试次数，重试间隔，最大重试超时时间
	if opt.RetryCount > 0 {
		if opt.RetryCount > defaultMaxRetries {
			opt.RetryCount = defaultMaxRetries // 最大重试次数
		}

		client = client.SetRetryCount(opt.RetryCount)

		if opt.RetryWaitTime != 0 {
			client = client.SetRetryWaitTime(opt.RetryWaitTime)
		}

		if opt.RetryMaxWaitTime != 0 {
			client = client.SetRetryMaxWaitTime(opt.RetryMaxWaitTime)
		}
	}

	//  设置cookie
	if cLen := len(opt.Cookies); cLen > 0 {
		cookies := make([]*http.Cookie, cLen)
		for k := range opt.Cookies {
			cookies = append(cookies, &http.Cookie{
				Name:     k,
				Value:    fmt.Sprintf("%v", opt.Cookies[k]),
				Path:     opt.CookiePath,
				Domain:   opt.CookieDomain,
				MaxAge:   opt.CookieMaxAge,
				HttpOnly: opt.CookieHttpOnly,
			})
		}

		client = client.SetCookies(cookies)
	}

	// 设置header头
	if len(opt.Headers) > 0 {
		client = client.SetHeaders(opt.ParseData(opt.Headers))
	}

	var resp *resty.Response
	var err error

	method = strings.ToLower(method)
	switch method {
	case "get", "delete", "head":
		client = client.SetQueryParams(opt.ParseData(opt.Params))
		if method == "get" {
			resp, err = client.R().Get(reqUrl)
			return s.GetResult(resp, err)
		}

		if method == "delete" {
			resp, err = client.R().Delete(reqUrl)
			return s.GetResult(resp, err)
		}

		if method == "head" {
			resp, err = client.R().Head(reqUrl)
			return s.GetResult(resp, err)
		}

	case "post", "put", "patch":
		req := client.R()
		if len(opt.Data) > 0 {
			//  SetFormData method sets Form parameters and their values in the current request.
			//  It's applicable only HTTP method `POST` and `PUT` and requests content type would be
			//  set as `application/x-www-form-urlencoded`.

			req = req.SetFormData(opt.ParseData(opt.Data))
		}

		// setBody: for struct and map data type defaults to 'application/json'
		//  SetBody method sets the request body for the request. It supports various realtime needs as easy.
		//  We can say its quite handy or powerful. Supported request body data types is `string`,
		//  `[]byte`, `struct`, `map`, `slice` and `io.Reader`. Body value can be pointer or non-pointer.
		//  Automatic marshalling for JSON and XML content type, if it is `struct`, `map`, or `slice`.
		if opt.Json != nil {
			req = req.SetBody(opt.Json)
		}

		if method == "post" {
			resp, err = req.Post(reqUrl)
			return s.GetResult(resp, err)
		}

		if method == "put" {
			resp, err = req.Put(reqUrl)
			return s.GetResult(resp, err)
		}

		if method == "patch" {
			resp, err = req.Patch(reqUrl)
			return s.GetResult(resp, err)
		}
	case "file":
		b, err := ioutil.ReadFile(opt.FileName)
		if err != nil {
			return &Reply{
				Err: errors.New("read file error: " + err.Error()),
			}
		}

		// 文件上传
		resp, err := client.R().
			SetFileReader(opt.FileParamName, opt.FileName, bytes.NewReader(b)).
			Post(reqUrl)
		return s.GetResult(resp, err)
	default:
	}

	return &Reply{
		Err: errors.New("request method not support"),
	}
}

//  NewClient创建一个resty客户端，支持post,get,delete,head,put,patch,file文件上传等
//  可以快速使用go-resty/resty上面的方法
//  参考文档： https:// github.com/go-resty/resty
func NewClient() *resty.Client {
	return resty.New()
}

//  GetData 处理请求的结果
func (s *Service) GetResult(resp *resty.Response, err error) *Reply {
	res := &Reply{}
	if err != nil {
		res.Err = err
		return res
	}

	// 请求返回的body
	res.Body = resp.Body()
	if !resp.IsSuccess() || resp.StatusCode() != 200 {
		res.Err = errors.New("request error: " + fmt.Sprintf("%v", resp.Error()) + "http StatusCode: " + strconv.Itoa(resp.StatusCode()) + "status: " + resp.Status())
		return res
	}

	return res
}

//  Text 返回Reply.Body文本格式
func (r *Reply) Text() string {
	return string(r.Body)
}

//  Json 将响应的结果Reply解析到data
//  对返回的Reply.Body做json反序列化处理
func (r *Reply) Json(data interface{}) error {
	if len(r.Body) > 0 {
		err := json.Unmarshal(r.Body, data)
		if err != nil {
			return err
		}
	}

	return nil
}
