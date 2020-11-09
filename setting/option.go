package setting

// Option option func for Setting.
type Option func(s *Setting)

// WithWatchFile 是否监听文件变化
func WithWatchFile(b bool) Option {
	return func(s *Setting) {
		s.watchFile = b
	}
}
