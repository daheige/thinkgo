package monitor

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var IsWebRequest bool //是否是web请求的监控

// 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of hello requests in total",
	},
	[]string{"method", "endpoint"}, //设置两个标签 请求方法和 路径 对请求总次数在两个
)

// web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

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
//可以对单个接口做metrics监控
func MonitorHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h(w, r)

		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	}
}

// MonitorHandler 性能监控处理器
// 可以作为中间件对接口进行打点监控
func MonitorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 meric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	})
}

// Check 心跳检测
func Check(w http.ResponseWriter, r *http.Request) {
	//模拟业务查询耗时0~1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	w.Write([]byte(`{"alive": true}`))
}

// PrometheusHandler 监控,包含了PProf性能监控
// 性能监控的端口port只能在内网访问
// 一般在程序启动init或main函数中执行
func PrometheusHandler(port int) {
	//注册监控指标
	if IsWebRequest {
		prometheus.MustRegister(WebRequestTotal)
		prometheus.MustRegister(WebRequestDuration)
	}

	prometheus.MustRegister(CpuTemp)
	prometheus.MustRegister(HdFailures)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PProf exec recover: ", err)
			}
		}()

		log.Println("server PProf run on: ", port)

		//创建一个http ServeMux实例
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", Check)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), httpMux); err != nil {
			log.Println("PProf listen error: ", err)
		}
	}()

}
