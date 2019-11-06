package yamlconf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

//配置文件结构体
type ConfigEngine struct {
	data map[string]interface{}
}

func NewConf() *ConfigEngine {
	return &ConfigEngine{}
}

func (c *ConfigEngine) GetData() map[string]interface{} {
	return c.data
}

// 加载yaml,yml内容到c.Data
func (c *ConfigEngine) LoadConf(path string) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// yaml解析的时候c.data如果没有被初始化，会自动为你做初始化
		err = yaml.Unmarshal(data, &c.data)
		if err != nil {
			return errors.New("can not parse " + path + " config")
		}

		return nil
	default:
		panic("error path ext")
	}

	return nil
}

// 从配置文件中获取值
func (c *ConfigEngine) Get(name string) interface{} {
	path := strings.Split(name, ".")
	data := c.data
	for key, value := range path {
		v, ok := data[value]
		if !ok {
			break
		}

		if (key + 1) == len(path) {
			return v
		}

		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			data = v.(map[string]interface{})
		}
	}

	return nil
}

// 从配置文件中获取string类型的值
func (c *ConfigEngine) GetString(name string, defaultValue string) string {
	return setString(c.Get(name), defaultValue)
}

func setString(value interface{}, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	if v, ok := value.(string); ok {
		return v
	}

	return defaultValue
}

// 从配置文件中获取int类型的值,defaultValue是默认值的
func (c *ConfigEngine) GetInt(name string, defaultValue int) int {
	return setInt(c.Get(name), defaultValue)
}

func setInt(v interface{}, defaultValue int) int {
	if v == nil {
		return defaultValue
	}

	//类型断言
	switch value := v.(type) {
	case string:
		i, _ := strconv.Atoi(value)
		return i
	case int:
		return value
	case bool:
		if value {
			return 1
		}

		return 0
	case float64:
		return int(value)
	default:
		return defaultValue
	}
}

func (c *ConfigEngine) GetInt64(name string, defaultValue int64) int64 {
	return setInt64(c.Get(name), defaultValue)
}

func setInt64(v interface{}, defaultValue int64) int64 {
	if v == nil {
		return defaultValue
	}

	//类型断言
	switch value := v.(type) {
	case int64:
		return value
	case string:
		i, _ := strconv.ParseInt(value, 10, 64)
		return i
	case int:
		return int64(value)
	case bool:
		if value {
			return 1
		}

		return 0
	case float64:
		return int64(value)
	default:
		return defaultValue
	}
}

// 从配置文件中获取bool类型的值
func (c *ConfigEngine) GetBool(name string, defaultValue bool) bool {
	return setBool(c.Get(name), defaultValue)
}

func setBool(v interface{}, defaultValue bool) bool {
	if v == nil {
		return defaultValue
	}

	switch value := v.(type) {
	case string:
		str, _ := strconv.ParseBool(value)
		return str
	case int:
		if value != 0 {
			return true
		}
		return false
	case bool:
		return value
	case float64:
		if value != 0.0 {
			return true
		}
		return false
	default:
		return defaultValue
	}
}

// 从配置文件中获取Float64类型的值
func (c *ConfigEngine) GetFloat64(name string, defaultValue float64) float64 {
	return setFloat64(c.Get(name), defaultValue)
}

func setFloat64(v interface{}, defaultValue float64) float64 {
	if v == nil {
		return defaultValue
	}

	switch value := v.(type) {
	case string:
		str, _ := strconv.ParseFloat(value, 64)
		return str
	case int:
		return float64(value)
	case bool:
		if value {
			return float64(1)
		}
		return float64(0)
	case float64:
		return value
	default:
		return defaultValue
	}
}

// 从配置文件中获取Struct类型的值
// 这里的struct是你自己定义的根据配置文件
// s必须是传递结构体的指针
// name是yaml定义的结构体名称
func (c *ConfigEngine) GetStruct(name string, s interface{}) interface{} {
	d := c.Get(name)
	// log.Printf("%T", d)

	switch d.(type) {
	case string:
		SetField(s, name, d)
	case map[interface{}]interface{}:
		//log.Println("s", s)
		mapToStruct(d.(map[interface{}]interface{}), s)
	}

	return s
}

// 将map转换为struct
// obj必须是一个结构体的指针
// 用data(map类型)填充结构obj
func mapToStruct(data map[interface{}]interface{}, obj interface{}) error {
	for k, v := range data {
		if v == nil { //当配置项的值是空，直接跳过
			continue
		}

		//打印k,v
		//log.Println("k = ", k)
		//log.Printf("k type = %T\n", k)
		//log.Println("v = ", v)

		if val, ok := k.(string); ok {
			//log.Println("key:", val)

			err := SetField(obj, val, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//设置obj结构体的k/v值
func SetField(obj interface{}, name string, value interface{}) error {
	// reflect.Indirect 返回value对应的值
	structValue := reflect.Indirect(reflect.ValueOf(obj))
	structFieldValue := structValue.FieldByName(name) //结构体的字段名

	// isValid 显示的测试一个空指针
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	// CanSet判断值是否可以被更改
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	// 获取要更改值的类型
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType.Kind() == reflect.Struct && val.Kind() == reflect.Map {
		vint := val.Interface()
		//类型断言,根据不同的类型设置k,v
		switch v := vint.(type) {
		case map[interface{}]interface{}:
			for key, value := range v {
				SetField(structFieldValue.Addr().Interface(), key.(string), value)
			}
		case map[string]interface{}:
			for key, value := range v {
				SetField(structFieldValue.Addr().Interface(), key, value)
			}
		}
	} else if structFieldType.Kind() == reflect.Slice && val.Kind() == reflect.Slice {
		//log.Println("k:", name)
		//log.Println("v", value)
		//arr := getSlice(value.([]interface{}))
		//log.Println("arr: ", arr)
		structFieldValue.Set(getSlice(value.([]interface{})))
	} else {
		if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}

	return nil
}

// getSlice 根据v []interface{}进行类型断言，返回指定类型的切片
func getSlice(v []interface{}) reflect.Value {
	vType := reflect.ValueOf(v[0]).Kind()
	vLen := len(v)

	switch vType {
	case reflect.String: //字符串类型
		arr := make([]string, 0, vLen)
		for _, _v := range v {
			arr = append(arr, setString(_v, ""))
		}

		return reflect.ValueOf(arr)
	case reflect.Int: //整数类型，这里不区分int32,int64统一用int类型
		arr := make([]int, 0, vLen)
		for _, _v := range v {
			arr = append(arr, setInt(_v, 0))
		}

		return reflect.ValueOf(arr)
	case reflect.Float64: //浮点类型统一用float64类型
		arr := make([]float64, 0, vLen)
		for _, _v := range v {
			arr = append(arr, setFloat64(_v, 0))
		}

		return reflect.ValueOf(arr)
	default:
	}

	//默认[]interface{}
	return reflect.ValueOf(v)
}
