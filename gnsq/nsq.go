/**
golang nsq消费队列封装，提供如下功能点
	1、初始化生产者
	2、初始化消费者
	3、提供不同方式的消费者消费模式
	4、当调用InitProducer,InitConsumer后可以直接调用nsq上底层方法
	也可以使用本包提供的方法，其实也是调用nsq底层方法
	5、关于优雅退出生产者和消费者，请看nsq_test.go
	6、通过直接连接到nsqd进行消费，速度快，但不方便拓展，建议通过lookupd查找节点进行消费
*/
package gnsq

import (
	"errors"
	"log"

	"github.com/nsqio/go-nsq"
)

var defaultGroutines = 1 //消费者连接到Nsqd内部groutine个数

// NewConfig 初始化nsq config
func NewConfig() *nsq.Config {
	return nsq.NewConfig()
}

// InitProducer 初始化生产者
// address是nsqd连接的tcp地址
func InitProducer(address string, conf *nsq.Config) (*nsq.Producer, error) {
	log.Println("nsqd tcp address: ", address)
	producer, err := nsq.NewProducer(address, conf)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// Publish 发布消息
// 当消息发送完毕后，需要producer.Stop() 让生产者优雅退出
func Publish(producer *nsq.Producer, topic string, msgBytes []byte) error {
	if len(msgBytes) == 0 { //不能发布空串，否则会导致error
		return errors.New("msg is empty")
	}

	if producer != nil {
		return producer.Publish(topic, msgBytes) // 发布消息
	}

	return errors.New("producer is nil")
}

// InitConsumer 初始化消费者
// 新建一个消费者
func InitConsumer(topic string, channel string, conf *nsq.Config) (*nsq.Consumer, error) {
	return nsq.NewConsumer(topic, channel, conf)
}

// ConsumerConnectToNSQLookupd 通过lookupd找到nsqd中的节点，进行消费
// nums是nsqd消费者内部指定goroutine个数
func ConsumerConnectToNSQLookupd(c *nsq.Consumer, address string, h nsq.Handler, nums int) error {
	if nums <= 0 {
		nums = defaultGroutines
	}

	c.SetLogger(nil, 0)              //屏蔽系统日志
	c.AddConcurrentHandlers(h, nums) //添加消费者接口

	//建立NSQLookupd连接
	if err := c.ConnectToNSQLookupd(address); err != nil {
		log.Println("nsq connection error: ", err)
		return err
	}

	return nil
}

// ConsumerConnectToNSQLookupds 通过lookupd找到nsqd中的节点，进行消费
// nums是nsqd消费者内部指定goroutine个数
// addressList 表示有多个lookupd地址
// hander消费者回调句柄是一个接口
func ConsumerConnectToNSQLookupds(c *nsq.Consumer, addressList []string, h nsq.Handler, nums int) error {
	if nums <= 0 {
		nums = defaultGroutines
	}

	c.SetLogger(nil, 0)              //屏蔽系统日志
	c.AddConcurrentHandlers(h, nums) //添加消费者接口

	//建立NSQLookupd连接
	if err := c.ConnectToNSQLookupds(addressList); err != nil {
		log.Println("nsq connection error: ", err)
		return err
	}

	return nil
}

// ConsumerConnectToNSQDs 消费者直接连接到单个nsqd进行消费
// hander消费者回调句柄是一个接口
func ConsumerConnectToNSQD(c *nsq.Consumer, address string, h nsq.Handler, nums int) error {
	if nums <= 0 {
		nums = defaultGroutines
	}

	c.SetLogger(nil, 0)              //屏蔽系统日志
	c.AddConcurrentHandlers(h, nums) //添加消费者接口

	//建立NSQd连接
	if err := c.ConnectToNSQD(address); err != nil {
		log.Println("nsq connection error: ", err)
		return err
	}

	return nil
}

// ConsumerConnectToNSQDs 消费者直接连接到多个nsqd进行消费
// hander消费者回调句柄是一个接口
func ConsumerConnectToNSQDs(c *nsq.Consumer, addressList []string, h nsq.Handler, nums int) error {
	if nums <= 0 {
		nums = defaultGroutines
	}

	c.SetLogger(nil, 0)              //屏蔽系统日志
	c.AddConcurrentHandlers(h, nums) //添加消费者接口

	//建立NSQd连接
	if err := c.ConnectToNSQDs(addressList); err != nil {
		log.Println("nsq connection error: ", err)
		return err
	}

	return nil
}
