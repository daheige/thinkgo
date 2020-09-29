package gresty

import (
	"time"
)

// Option service option.
type Option func(s *Service)

func (s *Service) apply(opts []Option) {
	for _, o := range opts {
		o(s)
	}
}

// WithBaseUri 设置请求地址url的前缀
func WithBaseUri(uri string) Option {
	return func(s *Service) {
		s.BaseUri = uri
	}
}

// WithTimeout 请求超时时间
func WithTimeout(d time.Duration) Option {
	return func(s *Service) {
		s.Timeout = d
	}
}

// WithProxy 设置请求proxy
func WithProxy(proxy string) Option {
	return func(s *Service) {
		s.Proxy = proxy
	}
}

// WithEnableKeepAlive 是否长连接模式，默认短连接
func WithEnableKeepAlive(b bool) Option {
	return func(s *Service) {
		s.EnableKeepAlive = b
	}
}
