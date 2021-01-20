package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
)

type Conf struct {
	Address  string
	Exchange string
	Key      string
}

type Publisher struct {
	conf    *Conf
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

func NewPublisher(conf *Conf) registry.PublisherInterface {
	return &Publisher{
		conf: conf,
		done: make(chan error),
	}
}

func (p *Publisher) Send(data []byte) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}

		//nolint:gosimple
		select {
		case <-time.After(d):
			if err := p.Connect(); err != nil {
				log.Printf("could not reconnect: %+v", err)
				continue
			}
			err := p.channel.Publish(
				p.conf.Exchange,
				p.conf.Key,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         data,
				},
			)
			if err != nil {
				fmt.Printf("failed to send data: %+v", err)
				continue
			}
			log.Printf(" [x] Sent %s", data)
			return nil
		}
	}
}

func (p *Publisher) Connect() error {
	conn, err := amqp.Dial(p.conf.Address)
	if err != nil {
		return err
	}
	p.conn = conn

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	go func() {
		log.Printf("closing: %s", <-p.conn.NotifyClose(make(chan *amqp.Error)))
		p.done <- errors.New("channel closed")
	}()

	return p.channel.ExchangeDeclare(
		p.conf.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (p *Publisher) Close() error {
	err := p.channel.Close()
	if err != nil {
		return err
	}

	return p.conn.Close()
}
