package monitor

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// WebRequestTotal 初始化 web_request_total， counter类型指标， 表示接收http请求总次数
// 设置两个标签 请求方法和 路径 对请求总次数在两个
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_request_total",
		Help: "Number of hello requests in total",
	},
	[]string{"method", "endpoint"},
)

// WebRequestDuration web_request_duration_seconds，
// Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

// CpuTemp cpu情况
var CpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU",
})

var HdFailures = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors",
	},
	[]string{"device"},
)

// MonitorHandlerFunc 对于http原始的处理器函数，包装 handler function,不侵入业务逻辑
// 可以对单个接口做metrics监控
func MonitorHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h(w, r)

		// counter类型 metric 的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 metric的记录方式
		WebRequestDuration.With(prometheus.Labels{
			"method": r.Method, "endpoint": r.URL.Path,
		}).Observe(time.Since(start).Seconds())
	}
}

// MonitorHandler 性能监控处理器
// 可以作为中间件对接口进行打点监控
func MonitorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		// counter类型 metric 的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 metric 的记录方式
		WebRequestDuration.With(prometheus.Labels{
			"method": r.Method, "endpoint": r.URL.Path,
		}).Observe(time.Since(start).Seconds())
	})
}
