# metrics性能监控
    使用方法：
    import "github.com/daheige/thinkgo/monitor"

    1、在init()方法中添加如下代码
    //注册监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

    2、在pprof中添加如下路由：
    //性能报告监控和健康检测
	//性能监控的端口port+1000,只能在内网访问
	go func() {
		//defer logger.Recover() //参考thinkgo/logger包

		pprof_address := fmt.Sprintf("0.0.0.0:%d", port+1000)
		log.Println("server pprof run on: ", pprof_address)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", routes.HealthCheck)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(pprof_address, httpMux); err != nil {
			log.Println(err)
		}
	}()

# 实战案例
    https://github.com/daheige/hg-mux
