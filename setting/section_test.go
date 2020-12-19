package setting

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// app.yaml section config.
var (
	appServerConf = &appServerSettingS{}
)

// appServerSettingS server config.
type appServerSettingS struct {
	AppEnv              string        `json:"app_env"`
	AppDebug            bool          `json:"app_debug"`
	GRPCPort            int           `json:"grpc_port"`
	GRPCHttpGatewayPort int           `json:"grpc_http_gateway_port"`
	HttpPort            int           `json:"http_port"`
	ReadTimeout         time.Duration `json:"read_timeout"`
	WriteTimeout        time.Duration `json:"write_timeout"`
	LogDir              string        `json:"log_dir"`
	JobPProfPort        int           `json:"job_p_prof_port"`
}

// readConfig 读取配置文件
func readConfig(configDir string) error {
	// 测试拓展名获取
	filename := "abc.yaml"
	log.Println(strings.TrimPrefix(filepath.Ext(filename), ".")) // yaml

	log.Println(filepath.Dir("/abc/app.yaml"))
	s, err := NewSetting(configDir, "test")
	if err != nil {
		return err
	}

	err = s.ReadSection("AppServer", &appServerConf)
	if err != nil {
		return err
	}

	appServerConf.ReadTimeout *= time.Second
	appServerConf.WriteTimeout *= time.Second

	if appServerConf.AppDebug {
		log.Println("app server config: ", appServerConf)
	}

	return nil
}

// TestReadSection test config read
/**
AppServer:
  AppEnv: dev
  AppDebug: true
  GRPCPort: 50051
  GRPCHttpGatewayPort: 1336
  HttpPort: 1338
  ReadTimeout: 6
  WriteTimeout: 6
  LogDir: ./logs
  JobPProfPort: 30031
*/
func TestReadSection(t *testing.T) {
	readConfig("./")
	b, _ := json.Marshal(appServerConf)
	log.Println("section app config: ", string(b))
}

/**
=== RUN   TestReadSection
2020/09/14 21:27:41 app server config:  &{dev true 50051 1336 1338 6s 6s ./logs 30031}
2020/09/14 21:27:41 section app config:
{"app_env":"dev","app_debug":true,"grpc_port":50051,"grpc_http_gateway_port":1336,
"http_port":1338,"read_timeout":6000000000,"write_timeout":6000000000,
"log_dir":"./logs","job_p_prof_port":30031}
--- PASS: TestReadSection (0.00s)
PASS
*/
