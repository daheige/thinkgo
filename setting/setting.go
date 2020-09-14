package setting

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Setting setting struct
type Setting struct {
	vp        *viper.Viper
	watchFile bool // 是否监听文件变化
}

// NewSetting create a setting entry.
func NewSetting(dir string, filename string, opts ...Option) (*Setting, error) {
	// 获取配置文件当前路径的绝对路径地址
	configDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln("config path error: ", err)
	}

	vp := viper.New()
	vp.SetConfigName(filename) // 配置文件名，不包含yaml后缀
	vp.AddConfigPath(configDir)

	vp.SetConfigType("yaml")
	err = vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}
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
