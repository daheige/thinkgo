// package httpRequest
// go http client support get,post,delete,patch,put,head,file method
// author:daheige

package httpRequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mozillazg/request"
)

var DefaultReqTimeout = 5 //默认请求超时，单位s

// Service 请求句柄设置
type Service struct {
	BaseUri string //请求地址url的前缀
	Timeout int    //请求超时限制
	Proxy   string //请求设置的http_proxy代理
}

// ReqOpt 请求参数设置
type ReqOpt struct {
	Params map[string]interface{} //get,delete的Params参数
	Data   map[string]interface{} //post请求的data表单数据
	Header map[string]interface{} //header头信息
	Cookie map[string]interface{} //cookie信息
	Method string                 //请求的方法get,post,put,patch,delete,head等

	//支持post,put,patch以json格式传递,[]int{1, 2, 3},map[string]string{"a":"b"}格式
	//json支持[],{}数据格式,主要是golang的基本数据类型，就可以
	Json interface{}

	//上传文件参数
	FieldName string //上传文件对应的表单file字段名
	File      string //上传文件名称,需要绝对路径
	FileName  string //上传后的文件名称
}

// Reply 请求后的结果
type Reply struct {
	Err  error  //请求过程中，发生的error
	Body []byte //返回的body内容
}

// ApiStdRes 标准的api返回格式
type ApiStdRes struct {
	Code    int
	Message string
	Data    interface{}
}

// ParseData 解析ReqOpt Params和Data
func (ReqOpt) ParseData(d map[string]interface{}) map[string]string {
	dLen := len(d)
	if dLen == 0 {
		return nil
	}

	//对d参数进行处理
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

// Do 请求方法
// uri    string  请求的相对地址，如果BaseUri为空，就必须是完整的url地址
func (s *Service) Do(method string, reqUrl string, opt *ReqOpt) *Reply {
	if s.BaseUri != "" {
		reqUrl = strings.TrimRight(s.BaseUri, "/") + "/" + reqUrl
	}

	if s.Timeout == 0 {
		s.Timeout = DefaultReqTimeout
	}

	if reqUrl == "" {
		return &Reply{
			Err: errors.New("request url is empty"),
		}
	}

	//短连接的形式请求api
	//关于如何关闭http connection
	//https://www.cnblogs.com/cobbliu/p/4517598.html

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(s.Timeout) * time.Second,
			KeepAlive: time.Duration(s.Timeout) * time.Second,
		}).DialContext,
		DisableKeepAlives: true, //禁用长连接
	}

	//http客户端
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(s.Timeout) * time.Second,
	}

	//创建请求客户端
	req := request.NewRequest(client)

	if s.Proxy != "" {
		req.Proxy = s.Proxy
	}

	//post,put,patch等支持json格式
	if opt.Json != nil {
		req.Json = opt.Json
	}

	//设置cookie
	if len(opt.Cookie) > 0 {
		req.Cookies = opt.ParseData(opt.Cookie)
	}

	//设置header头
	if len(opt.Header) > 0 {
		req.Headers = opt.ParseData(opt.Header)
	}

	method = strings.ToLower(method)
	switch method {
	case "get":
		req.Params = opt.ParseData(opt.Params)
		resp, err := req.Get(reqUrl)

		return s.GetData(resp, err)
	case "post":
		if len(req.Data) > 0 {
			req.Data = opt.ParseData(opt.Data) //请求的数据data
		}

		resp, err := req.Post(reqUrl)

		return s.GetData(resp, err)
	case "put":
		req.Data = opt.ParseData(opt.Data)
		resp, err := req.Put(reqUrl)

		return s.GetData(resp, err)
	case "patch":
		req.Data = opt.ParseData(opt.Data)
		resp, err := req.Patch(reqUrl)

		return s.GetData(resp, err)
	case "delete":
		req.Params = opt.ParseData(opt.Params)
		resp, err := req.Delete(reqUrl)

		return s.GetData(resp, err)
	case "head":
		req.Params = opt.ParseData(opt.Params)
		resp, err := req.Head(reqUrl)

		return s.GetData(resp, err)
	case "file":
		if opt.FileName == "" {
			return &Reply{
				Err: errors.New("upload fileName is empty"),
			}
		}

		fd, err := os.Open(opt.File)
		if err != nil {
			return &Reply{
				Err: errors.New(fmt.Sprintf("open file:%s error:%s", opt.File, err.Error())),
			}
		}

		defer fd.Close()

		req.Files = []request.FileField{
			request.FileField{opt.FieldName, opt.FileName, fd},
		}

		resp, err := req.Post(reqUrl)
		return s.GetData(resp, err)
	default:
	}

	return &Reply{
		Err: errors.New("request method not support"),
	}
}

// GetData 处理请求的结果
func (s *Service) GetData(resp *request.Response, err error) *Reply {
	res := &Reply{}
	if err != nil {
		res.Err = err
		return res
	}

	//resp.Context() will auto close body
	//调用Content方法会自动关闭resp句柄
	res.Body, res.Err = resp.Content()
	if resp.StatusCode != 200 || !resp.OK() {
		res.Err = errors.New(resp.Reason())
		return res
	}

	return res
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
