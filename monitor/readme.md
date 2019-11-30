# metrics性能监控
    使用方法：
    import "github.com/daheige/thinkgo/monitor"

    1、在init()方法中添加如下代码
    //注册监控指标
    
    //web程序的性能监控，如果是job/rpc服务就不需要这两行
    prometheus.MustRegister(WebRequestTotal)
    prometheus.MustRegister(WebRequestDuration)
    
    	
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

    2、在pprof中添加如下路由：
    //性能报告监控和健康检测
	//性能监控的端口只能在内网访问
	var PProfPort = 2338
	go func() {
		//defer logger.Recover() //参考thinkgo/logger包

		PProfAddress := fmt.Sprintf("0.0.0.0:%d", PProfPort)
		log.Println("server pprof run on: ", PProfAddress)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", routes.HealthCheck)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(PProfAddress, httpMux); err != nil {
			log.Println(err)
		}
	}()
	
	也可以直接调用 PrometheusHandler 监控,包含了PProf性能监控

# 实战案例

    https://github.com/daheige/hg-mux
