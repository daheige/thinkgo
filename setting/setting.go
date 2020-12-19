package setting

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Setting setting struct
type Setting struct {
	vp        *viper.Viper
	watchFile bool                   // 是否监听文件变化
	sections  map[string]interface{} // 存放key/val配置项
}

// NewSetting create a setting entry.
// 默认filename文件类型是yaml文件，其他类型ini,json,yml等格式也支持
func NewSetting(dir string, filename string, opts ...Option) (*Setting, error) {
	// 获取配置文件当前路径的绝对路径地址
	configDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln("config path error: ", err)
	}

	if filename == "" {
		log.Fatalln("config filename is empty")
	}

	// 文件拓展名
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext == "" {
		ext = "yaml" // 默认读取yaml文件
	}

	// 文件名不带后缀
	file := filepath.Base(filename)

	vp := viper.New()
	vp.SetConfigName(file) // 配置文件名，不包含文件后缀名称
	vp.AddConfigPath(configDir)

	vp.SetConfigType(ext)
	err = vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// 初始化Setting
	s := &Setting{vp: vp, sections: make(map[string]interface{}, 20)}
	for _, o := range opts {
		o(s)
	}

	if s.watchFile {
		s.WatchSettingChange()
	}

	return s, nil
}

// WatchSettingChange watch file change.
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}

// GetVp 返回viper.Viper指针对象
func (s *Setting) GetVp() *viper.Viper {
	return s.vp
}
