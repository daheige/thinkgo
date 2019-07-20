# golang nsq
    golang nsq消费队列封装，提供如下功能点
	1、初始化生产者
	2、初始化消费者
	3、提供不同方式的消费者消费模式
	4、当调用InitProducer,InitConsumer后可以直接调用nsq上底层方法
	也可以使用本包提供的方法，其实也是调用nsq底层方法
	5、关于优雅退出生产者和消费者，请看nsq_test.go
	6、通过直接连接到nsqd进行消费，速度快，但不方便拓展，建议通过lookupd查找节点进行消费