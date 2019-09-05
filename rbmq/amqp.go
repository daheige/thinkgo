package rbmq

import (
	"flag"
	"fmt"

	"github.com/streadway/amqp"
)

const (
	Transient  uint8 = 1 //暂存的消息
	Persistent uint8 = 2 //消息持久化

	//消息优先级
	PriorityMax uint8 = 9
	PriorityMin uint8 = 0
)

var (
	flgPrefetch = flag.Int("prefetch", 64, "prefetch message from mq")
)

// ConnectMq new a rbmq connection,return conn,ch,error
func ConnectMq(url string) (conn *amqp.Connection, channel *amqp.Channel, err error) {
	conn, err = amqp.Dial(url)
	if err == nil {
		channel, err = conn.Channel()
	}

	return
}

// NewMqExchange New a Exchange
// _type: 四种exchange类型：direct, topic, headers, fanout
//  fanout类型的exchange很简单，顾名思义：它将所接收到的消息广播给所有绑定的队列
func NewMqExchange(channel *amqp.Channel, name, _type string, durable bool) error {
	return channel.ExchangeDeclare(
		name,    // name
		_type,   // type
		durable, // durable
		false,   // auto-delete
		false,   // internal
		false,   // nowait
		nil,     // args
	)
}

// Publish publish a msg
func Publish(channel *amqp.Channel, exchange, rkey string, msg []byte) error {
	if channel == nil {
		return fmt.Errorf("channel is nil")
	}

	return channel.Publish(exchange, rkey, false, false,
		amqp.Publishing{ContentType: "application/octet-stream", Body: msg})
}

// PriorityPublish 指定优先级priority和是否消息持久化DeliveryMode
func PriorityPublish(channel *amqp.Channel, exchange, rkey string, savedisk, priority uint8, msg []byte) error {
	if channel == nil {
		return fmt.Errorf("channel is nil")
	}

	return channel.Publish(exchange, rkey, false, false,
		amqp.Publishing{ContentType: "application/octet-stream", Body: msg,
			DeliveryMode: savedisk, Priority: priority})
}

// NewMqQueue declare a queue
// 声明一个队列，可以指定exchange,routing key
// durable 消息是否持久化
func NewMqQueue(channel *amqp.Channel, exchange, queue, rkey string, durable, exclusive bool) error {
	if _, err := channel.QueueDeclare(
		queue,     // name of the queue
		durable,   // durable
		exclusive, // delete when usused
		exclusive, // exclusive
		false,     // noWait
		nil,       // arguments
	); err != nil {
		return err
	}

	if err := channel.QueueBind(
		queue,    // name of the queue
		rkey,     // bindingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return err
	}

	return nil
}

func newMqConsumer(url, exchange, queue, rkey, ctag string, ack, durable, exclusive bool) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	conn, channel, err = ConnectMq(url)
	if err != nil {
		return
	}

	err = NewMqQueue(channel, exchange, queue, rkey, durable, exclusive)
	if err != nil {
		return
	}

	deliveries, err = channel.Consume(
		queue,     // name
		ctag,      // consumerTag,
		!ack,      // noAck
		exclusive, // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		conn.Close()
		return
	}

	if ack {
		channel.Qos(*flgPrefetch, 0, true)
	}
	return
}

// NewMqConsumer make a msg consumer创建一个消费者
func NewMqConsumer(url, exchange, queue, rkey, ctag string, ack, durable bool) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	return newMqConsumer(url, exchange, queue, rkey, ctag, ack, durable, false)
}

// NewExclusiveMqConsumer new a exclusive mq consumer
func NewExclusiveMqConsumer(url, exchange, queue, rkey, ctag string, ack, durable bool) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	return newMqConsumer(url, exchange, queue, rkey, ctag, ack, durable, true)
}
