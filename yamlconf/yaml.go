package yamlconf

import (
	"log"
	"path/filepath"

	"github.com/daheige/thinkgo/setting"
)

// ConfigEngine 配置文件结构体
type ConfigEngine struct {
	// 配置文件目录
	dir string

	// 比如app.yaml,app.env,app.yml,app.json
	filename string

	// 监听文件变化
	watchFile bool

	s *setting.Setting
}

// NewConf new an entry
func NewConf(opts ...Option) *ConfigEngine {
	conf := &ConfigEngine{}
	conf.Apply(opts)

	return conf
}

// Apply 应用函数配置项
func (conf *ConfigEngine) Apply(opts []Option) {
	for _, o := range opts {
		o(conf)
	}
}

// GetData get data
func (c *ConfigEngine) GetData() map[string]interface{} {
	return c.s.GetSections()
}

// LoadConf 加载yaml,yml内容,兼容之前的版本
func (c *ConfigEngine) LoadConf(path string) error {
	c.dir = filepath.Dir(path)
	c.filename = filepath.Base(path)
	return c.LoadData()
}

// LoadData 加载yaml,yml内容
func (c *ConfigEngine) LoadData() error {
	var err error
	c.s, err = setting.NewSetting(c.dir, c.filename, setting.WithWatchFile(c.watchFile))
	if err != nil {
		return err
	}

	return nil
}

// Get 从配置文件中获取值,v必须是一个指针类型
func (c *ConfigEngine) Get(name string, v interface{}) error {
	return c.s.ReadSection(name, v)
}

// GetString 从配置文件中获取string类型的值
func (c *ConfigEngine) GetString(name string, defaultValue string) string {
	var str string
	err := c.Get(name, &str)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	if str == "" {
		return defaultValue
	}

	return str
}

// GetInt 从配置文件中获取int类型的值,defaultValue是默认值的
func (c *ConfigEngine) GetInt(name string, defaultValue int) int {
	var i int
	err := c.Get(name, &i)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	if i == 0 {
		return defaultValue
	}

	return i
}

// GetInt64 get int64
func (c *ConfigEngine) GetInt64(name string, defaultValue int64) int64 {
	var i int64
	err := c.Get(name, &i)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	if i == 0 {
		return defaultValue
	}

	return i
}

// GetBool 从配置文件中获取bool类型的值
func (c *ConfigEngine) GetBool(name string, defaultValue bool) bool {
	var b bool
	err := c.Get(name, &b)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	return b
}

// GetFloat64 从配置文件中获取Float64类型的值
func (c *ConfigEngine) GetFloat64(name string, defaultValue float64) float64 {
	var f float64
	err := c.Get(name, &f)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	if f == 0 {
		return defaultValue
	}

	return f
}

// GetFloat32 获取浮点数float32
func (c *ConfigEngine) GetFloat32(name string, defaultValue float32) float32 {
	var f float32
	err := c.Get(name, &f)
	if err != nil {
		log.Printf("get key of %s error: %s", name, err.Error())
		return defaultValue
	}

	if f == 0 {
		return defaultValue
	}

	return f
}

// GetStruct 从配置文件中获取Struct类型的值
// 这里的struct是你自己定义的根据配置文件
// s必须是传递结构体的指针
// name是yaml定义的结构体名称
func (c *ConfigEngine) GetStruct(name string, s interface{}) interface{} {
	return c.s.ReadSection(name, s)
}
