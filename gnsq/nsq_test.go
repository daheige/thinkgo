package gnsq

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

// 通过直接连接到nsqd进行消费，速度快，但不方便拓展，建议通过lookupd查找节点进行消费
var (
	// nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.0.1:4152" // nsqd连接的tcp地址
	lookupdAddrr = "127.0.0.1:4161" // lookupd http地址
)

// go test -v -test.run TestProduction
// 查看消息发送的结果 http://localhost:4171/topics/test
/*
send msg success
2019/07/20 16:14:07 production will exit
--- PASS: TestProduction (229.65s)
    nsq_test.go:17: test nsq production
PASS
ok      github.com/daheige/thinkgo/gNsq 229.649s
*/
func TestProduction(t *testing.T) {
	t.Log("test nsq production")
	conf := NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second // 心跳检查

	// 创建生产者
	tPro, err := InitProducer(tcpNsqdAddrr, conf)

	if err != nil {
		fmt.Println("new producer err:", err)
	}

	// 平滑退出
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// 测试发送100w消息
	nums := 100
	for i := 0; i < nums; i++ {
		// 监听是否有退出信号
		select {
		case sig := <-ch: // 接收到停止信号，就优雅的退出发送
			signal.Stop(ch) // 停止接收信号

			log.Println("exit signal: ", sig.String())
			// 优雅的停止发送
			// Stop initiates a graceful stop of the Producer (permanent)
			tPro.Stop()
			goto exit // 如果是退出函数可以写return
		default:
			log.Println("msg sending...")
		}

		// 主题
		topic := "test"
		// 主题内容
		tCommand := "hello:" + strconv.Itoa(i)
		// 发布消息

		// 这里可以调用tPro上面的方法Publish进行发送消息
		// err = tPro.Publish(topic, []byte(tCommand))

		// 可以调用这个Publish进行发送
		err = Publish(tPro, topic, []byte(tCommand))
		if err != nil {
			fmt.Println("publis msg error: ", err)
			continue
		}

		fmt.Println("current index: ", i)
		fmt.Println("send msg success")
	}

exit:
	log.Println("production will exit")

}

// 声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
// 实现nsq 底层的Handler
type nsqHandler struct {
	// 消息数
	msqCount int64
	// 标识ID
	nsqHandlerID string
}

// 实现 Handler接口上的HandleMessage方法
// message是接收到的消息
func (s *nsqHandler) HandleMessage(message *nsq.Message) error {
	// 每收到一条消息+1
	s.msqCount++
	// 打印输出信息和ID
	fmt.Println(s.msqCount, s.nsqHandlerID)
	// 打印消息的一些基本信息
	fmt.Printf("msg.Timestamp=%v, msg.nsqaddress=%s,msg.body=%s \n",
		time.Unix(0, message.Timestamp).Format("2006-01-02 03:04:05"),
		message.NSQDAddress, string(message.Body))
	return nil
}

// go test -v -test.run TestCust
/**
^C2019/07/20 16:26:35 exit signal:  interrupt
2019/07/20 16:26:35 INF    1 [test/channel1] stopping...
2019/07/20 16:26:35 shutting down
--- PASS: TestCust (259.80s)
    nsq_test.go:121: test cust nsq
2019/07/20 16:26:35 INF    1 [test/channel1] (127.0.0.1:4152) received CLOSE_WAIT from nsqd
2019/07/20 16:26:35 INF    1 [test/channel1] (127.0.0.1:4152) beginning close
2019/07/20 16:26:35 INF    1 [test/channel1] (127.0.0.1:4152) readLoop exiting
2019/07/20 16:26:35 INF    1 [test/channel1] (127.0.0.1:4152) breaking out of writeLoop
2019/07/20 16:26:35 INF    1 [test/channel1] (127.0.0.1:4152) writeLoop exiting
PASS
ok      github.com/daheige/thinkgo/gNsq 259.807s
*/
func TestCust(t *testing.T) {
	t.Log("test cust nsq")

	// 初始化配置
	conf := NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second // 心跳检查

	// 创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := InitConsumer("test", "channel1", conf)
	if err != nil {
		fmt.Println(err)
	}

	// 添加处理回调
	// com.AddHandler(&NsqHandler{nsqHandlerID: "One"}) //默认是单个goroutine处理消息

	// 通过并发的方式消费
	// 指定10个goroutine内部消费
	com.AddConcurrentHandlers(&nsqHandler{nsqHandlerID: "One"}, 10)
	// 连接对应的nsqd
	// err = com.ConnectToNSQD(tcpNsqdAddrr)

	// 通过lookupd查询到nsqd节点后，连接到对应的nsqd
	err = com.ConnectToNSQLookupd(lookupdAddrr)
	if err != nil {
		fmt.Println(err)
	}

	// 平滑退出
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())

	// 优雅的停止消费者
	com.Stop()

	log.Println("shutting down")
}

// 测试100条消息消费
// 通过lookupd查找nsqd节点后，进行连接消费
// go test -v -test.run TestCust2
/**
^C2019/07/20 16:35:13 exit signal:  interrupt
2019/07/20 16:35:13 shutting down
--- PASS: TestCust2 (1.20s)
    nsq_test.go:194: test cust nsq
PASS
ok      github.com/daheige/thinkgo/gNsq 1.207s
*/
func TestCust2(t *testing.T) {
	t.Log("test cust nsq")

	// 初始化配置
	conf := NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second // 心跳检查

	// 创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := InitConsumer("test", "channel1", conf)
	if err != nil {
		fmt.Println(err)
	}

	// 指定10个goroutine内部消费
	err = ConsumerConnectToNSQLookupd(com, lookupdAddrr, &nsqHandler{nsqHandlerID: "One"}, 10)
	if err != nil {
		log.Println("exec cust: ", err)
	}

	// 平滑退出
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())

	// 优雅的停止消费者
	com.Stop()

	log.Println("shutting down")
}

// 直接连接到nsqd上进行消费
// 测试100条消息消费
// go test -v -test.run TestCust3
/**
^C2019/07/20 16:35:38 exit signal:  interrupt
2019/07/20 16:35:38 shutting down
--- PASS: TestCust3 (0.58s)
    nsq_test.go:245: test cust nsq
PASS
ok      github.com/daheige/thinkgo/gNsq 0.588s
*/
func TestCust3(t *testing.T) {
	t.Log("test cust nsq")

	// 初始化配置
	conf := NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second // 心跳检查

	// 创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := InitConsumer("test", "channel1", conf)
	if err != nil {
		fmt.Println(err)
	}

	// 指定10个goroutine内部消费
	err = ConsumerConnectToNSQD(com, tcpNsqdAddrr, &nsqHandler{nsqHandlerID: "One"}, 10)
	if err != nil {
		log.Println("exec cust: ", err)
	}

	// 平滑退出
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())

	// 优雅的停止消费者
	com.Stop()

	log.Println("shutting down")
}
