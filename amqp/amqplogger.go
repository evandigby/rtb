package amqp

import (
	"encoding/json"
	"github.com/evandigby/rtb"
	"github.com/streadway/amqp"
)

type AmqpBidLogger struct {
	addr      string
	appDomain string
}

func (l *AmqpBidLogger) logName() string {
	return l.appDomain + ":" + "log"
}

func (l *AmqpBidLogger) queueName() string {
	return l.appDomain + ":" + "logQueue"
}

func (l *AmqpBidLogger) LogItem(logItem *rtb.BidLogItem) {
	conn, err := amqp.Dial(l.addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	name := l.logName()

	err = ch.ExchangeDeclare(
		name,     // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	js, err := json.Marshal(logItem)

	if err != nil {
		return
	}

	err = ch.Publish(
		name,  // exchange
		"",    // routing key
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "applications/json",
			Body:        js,
		})

	if err != nil {
		panic(err)
	}

}

func (l *AmqpBidLogger) LogChannel() chan *rtb.BidLogItem {
	conn, err := amqp.Dial(l.addr)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	name := l.logName()
	queueName := l.queueName()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	err = ch.QueueBind(
		queueName, // queue name
		"",        // routing key
		name,      // exchange
		false,
		nil)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	logItems := make(chan *rtb.BidLogItem)

	go func() {
		for d := range msgs {
			var logItem rtb.BidLogItem

			err = json.Unmarshal(d.Body, &logItem)

			logItems <- &logItem
		}
	}()

	return logItems
}

func NewAmqpBidLogger(addr string, appDomain string) rtb.BidLogger {
	l := new(AmqpBidLogger)
	l.addr = addr
	l.appDomain = appDomain
	return l
}
