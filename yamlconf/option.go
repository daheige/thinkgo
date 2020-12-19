package yamlconf

// Option 配置项采用功能函数模式
type Option func(conf *ConfigEngine)

// WithDir 设置配置文件目录
func WithDir(dir string) Option {
	return func(conf *ConfigEngine) {
		conf.dir = dir
	}
}

// WithFilename 设置配置文件名，比如app.yaml,app.env,app.yml,app.json
func WithFilename(f string) Option {
	return func(conf *ConfigEngine) {
		conf.filename = f
	}
}

// WithWatchFile 是否监听文件变化
func WithWatchFile(b bool) Option {
	return func(conf *ConfigEngine) {
		conf.watchFile = b
	}
}
