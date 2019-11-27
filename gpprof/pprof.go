package gpprof

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
)

// PProfHandler PProf性能监控
// 性能监控的端口port只能在内网访问
// 一般在程序启动init或main函数中执行
func PProfHandler(port int) {
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
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), httpMux); err != nil {
			log.Println("PProf listen error: ", err)
		}

	}()

}

// Check PProf心跳检测
func Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"alive": true}`))
}
