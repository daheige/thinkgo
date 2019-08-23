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
	Url     string                 //请求的相对地址，如果BaseUri为空，就必须是完整的url地址
	Params  map[string]interface{} //get,delete的Params参数
	BaseUri string                 //请求地址url的前缀
	Data    map[string]interface{} //post请求的data表单数据
	Header  map[string]interface{} //header头信息
	Cookie  map[string]interface{} //cookie信息
	Method  string                 //请求的方法get,post,put,patch,delete,head等
	Proxy   string                 //请求设置的http_proxy代理

	//支持post,put,patch以json格式传递,[]int{1, 2, 3},map[string]string{"a":"b"}格式
	//json支持[],{}数据格式,主要是golang的基本数据类型，就可以
	Json interface{}
	FileField
}

type FileField struct {
	FieldName string //上传文件对应的表单file字段名
	File      string //上传文件名称,需要绝对路径
	FileName  string //上传后的文件名称
}

type Result struct {
	Err  error  //请求过程中，发生的error
	Body []byte //返回的body内容
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
		req.Params = a.ParseData(a.Params)
		resp, err := req.Get(a.Url)

		return a.GetData(resp, err)
	case "post":
		if len(req.Data) > 0 {
			req.Data = a.ParseData(a.Data) //请求的数据data
		}

		resp, err := req.Post(a.Url)

		return a.GetData(resp, err)
	case "put":
		req.Data = a.ParseData(a.Data)
		resp, err := req.Put(a.Url)

		res = a.GetData(resp, err)
	case "patch":
		req.Data = a.ParseData(a.Data)
		resp, err := req.Patch(a.Url)

		return a.GetData(resp, err)
	case "delete":
		req.Params = a.ParseData(a.Params)
		resp, err := req.Delete(a.Url)

		return a.GetData(resp, err)
	case "head":
		req.Params = a.ParseData(a.Params)
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

func (a *ApiRequest) ParseData(d map[string]interface{}) map[string]string {
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

//处理请求的结果
func (a *ApiRequest) GetData(resp *request.Response, err error) *Result {
	res := &Result{}
	if err != nil {
		res.Err = err
		return res
	}

	res.Body, res.Err = resp.Content()
	if resp.StatusCode != 200 || !resp.OK() {
		res.Err = errors.New(resp.Reason())
		return res
	}

	return res
}

//对返回的result.Body做json反序列化处理
func (result *Result) ParseJson() (*ApiStdRes, error) {
	res := &ApiStdRes{}
	err := result.Json(res)
	return res, err
}

func (result *Result) Text() string {
	return string(result.Body)
}

func (result *Result) Json(data interface{}) error {
	if len(result.Body) > 0 {
		err := json.Unmarshal(result.Body, data)
		if err != nil {
			return err
		}
	}

	return nil
}
