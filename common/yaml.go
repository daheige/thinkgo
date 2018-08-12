package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

//配置文件结构体
type ConfigEngine struct {
	Data map[interface{}]interface{}
}

func New() *ConfigEngine {
	return &ConfigEngine{}
}

// 将ymal文件中的内容进行加载
func (c *ConfigEngine) LoadConf(path string) error {
	ext := c.guessFileType(path)
	if ext == "" {
		return errors.New("cant not load" + path + " config")
	}

	return c.loadFromYaml(path)
}

//判断配置文件名是否为yaml格式
func (c *ConfigEngine) guessFileType(path string) string {
	s := strings.Split(path, ".")
	ext := s[len(s)-1]
	switch ext {
	case "yaml", "yml":
		return "yaml"
	}
	return ""
}

// 将配置yaml文件中的进行加载
func (c *ConfigEngine) loadFromYaml(path string) error {
	yamlS, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	// yaml解析的时候c.data如果没有被初始化，会自动为你做初始化
	err := yaml.Unmarshal(yamlS, &c.Data)
	if err != nil {
		return errors.New("can not parse " + path + " config")
	}

	return nil
}

// 从配置文件中获取值
func (c *ConfigEngine) Get(name string) interface{} {
	path := strings.Split(name, ".")
	data := c.Data
	for key, value := range path {
		v, ok := data[value]
		if !ok {
			break
		}

		if (key + 1) == len(path) {
			return v
		}

		if reflect.TypeOf(v).String() == "map[interface {}]interface {}" {
			data = v.(map[interface{}]interface{})
		}
	}

	return nil
}

// 从配置文件中获取string类型的值
func (c *ConfigEngine) GetString(name string, defaultValue string) string {
	value := c.Get(name)
	if value == nil {
		return defaultValue
	}

	switch value := value.(type) {
	case string:
		return value
	case bool, float64, int:
		return fmt.Sprintf("%s", value)
	default:
		return defaultValue
	}

}

// 从配置文件中获取int类型的值,defaultValue是默认值的
func (c *ConfigEngine) GetInt(name string, defaultValue int) int {
	value := c.Get(name)
	if value == nil {
		return defaultValue
	}

	//类型断言
	switch value := value.(type) {
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
	value := c.Get(name)
	if value == nil {
		return defaultValue
	}

	//类型断言
	switch value := value.(type) {
	case int64:
		return value
	case string:
		i, _ := strconv.ParseInt("1234", 10, 64)
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
	value := c.Get(name)
	if value == nil {
		return defaultValue
	}

	switch value := value.(type) {
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
	value := c.Get(name)
	if value == nil {
		return defaultValue
	}

	switch value := value.(type) {
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

// 从配置文件中获取Struct类型的值,这里的struct是你自己定义的根据配置文件
// s必须是传递结构体的指针
// name是yaml定义的结构体名称
func (c *ConfigEngine) GetStruct(name string, s interface{}) interface{} {
	d := c.Get(name)
	switch d.(type) {
	case string:
		c.setField(s, name, d)
	case map[interface{}]interface{}:
		c.mapToStruct(d.(map[interface{}]interface{}), s)
	}

	return s
}

// map转换为struct
// s必须是一个结构体的指针
func (c *ConfigEngine) mapToStruct(m map[interface{}]interface{}, s interface{}) {
	for key, value := range m {
		switch key.(type) {
		case string:
			c.setField(s, key.(string), value)
		}
	}

}

// 利用反射设置obj的属性值,这部分代码是重点，需要多看看
func (c *ConfigEngine) setField(obj interface{}, name string, value interface{}) error {
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

		switch vint.(type) {
		case map[interface{}]interface{}:
			for key, value := range vint.(map[interface{}]interface{}) {
				c.setField(structFieldValue.Addr().Interface(), key.(string), value)
			}
		case map[string]interface{}:
			for key, value := range vint.(map[string]interface{}) {
				c.setField(structFieldValue.Addr().Interface(), key, value)
			}
		}

	} else {
		if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}

	return nil
}
