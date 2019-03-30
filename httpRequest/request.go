package httpRequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mozillazg/request"
)

type ApiRequest struct {
	Url     string                 //请求的相对地址，如果BaseUri为空，就必须是绝对路径
	Params  map[string]interface{} //get,delete的Params参数
	BaseUri string                 //请求地址url的前缀
	Body    map[string]interface{}
	Header  map[string]interface{}
	Cookie  map[string]interface{}
	Method  string //请求的方法get,post,put,patch,delete,head等
	Proxy   string //请求设置的http_proxy代理

	//支持post,put,patch以json格式传递,[]int{1, 2, 3},map[string]string{"a":"b"}格式
	//也可以是json字符串
	Json interface{}
	FileField
}

type FileField struct {
	FieldName string //上传文件对应的表单file字段名
	File      string //上传文件名称,需要绝对路径
	FileName  string //上传后的文件名称
}

type Result struct {
	Err  error
	Body string
}

//标准的api返回格式
type ApiStdRes struct {
	Code    int
	Message string
	Data    interface{}
}

func (a *ApiRequest) Do() *Result {
	res := &Result{}
	a.Url = a.BaseUri + a.Url
	if a.Url == "" {
		res.Err = errors.New("request url is empty")
		return res
	}

	client := &http.Client{} //http客户端
	req := request.NewRequest(client)

	if a.Proxy != "" {
		req.Proxy = a.Proxy
	}

	//post,put,patch等支持json格式
	if a.Json != nil {
		req.Json = a.Json
	}

	//设置cookie
	if len(a.Cookie) > 0 {
		req.Cookies = a.ParseData(a.Cookie)
	}

	//设置header头
	if len(a.Header) > 0 {
		req.Headers = a.ParseData(a.Header)
	}

	a.Method = strings.ToLower(a.Method)
	switch a.Method {
	case "get":
		if len(a.Params) > 0 {
			req.Params = a.ParseData(a.Params)
		}

		resp, err := req.Get(a.Url)
		return a.GetData(resp, err)
	case "post":
		if len(a.Body) > 0 {
			req.Data = a.ParseData(a.Body) //请求的数据data
		}

		resp, err := req.Post(a.Url)
		return a.GetData(resp, err)
	case "put":
		if len(a.Body) == 0 {
			res.Err = errors.New("put data is empty")
			return res
		}

		req.Data = a.ParseData(a.Body)
		resp, err := req.Put(a.Url)
		res = a.GetData(resp, err)
	case "patch":
		if len(a.Body) == 0 {
			res.Err = errors.New("put data is empty")
			return res
		}

		req.Data = a.ParseData(a.Body)
		resp, err := req.Put(a.Url)
		return a.GetData(resp, err)
	case "delete":
		if len(a.Params) > 0 {
			req.Params = a.ParseData(a.Params)
		}

		resp, err := req.Delete(a.Url)
		return a.GetData(resp, err)
	case "head":
		resp, err := req.Head(a.Url)
		res = a.GetData(resp, err)
	case "file":
		if a.FileName == "" {
			res.Err = errors.New("file not exist")
			return res
		}

		fd, err := os.Open(a.FileName)
		if err != nil {
			res.Err = errors.New(fmt.Sprintf("open file:%s error:%s", a.FileName, err.Error()))
			return res
		}

		defer fd.Close()

		req.Files = []request.FileField{
			request.FileField{a.FieldName, a.FileName, fd},
		}

		resp, err := req.Post(a.Url)
		return a.GetData(resp, err)
	default:
	}

	res.Err = errors.New("request method not support")
	return res
}

func (a *ApiRequest) SetFileName(fileName string) {
	a.FileName = fileName
}

func (a *ApiRequest) ParseData(d map[string]interface{}) (data map[string]string) {
	if len(d) > 0 {
		for k, v := range a.Body {
			if val, ok := v.(string); ok {
				data[k] = val
			}

			data[k] = fmt.Sprintf("%v", v)
		}
	}

	return
}

//处理请求的结果
func (a *ApiRequest) GetData(resp *request.Response, err error) *Result {
	res := &Result{}
	if err != nil {
		res.Err = err
		return res
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 || !resp.OK() {
		res.Err = errors.New(resp.Reason())
		return res
	}

	body, err := resp.Text()
	if err != nil {
		res.Err = err
		res.Body = body
		return res
	}

	res.Body = body

	return res
}

//对返回的result.Body做json反序列化处理
func (result *Result) ParseJson() (res *ApiStdRes, e error) {
	if result.Body != "" {
		err := json.Unmarshal([]byte(result.Body), res)
		if err != nil {
			e = err
			return
		}
	}

	return
}
