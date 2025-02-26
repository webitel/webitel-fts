package handler

import (
	"context"
	"encoding/json"
	"github.com/webitel/webitel-fts/infra/pubsub"
	"github.com/webitel/webitel-fts/pkg/client"
	"github.com/webitel/wlog"
)

const XDLExpire = 86400000 * 7 // 7 days

type SubscriberService interface {
	Create(ctx context.Context, msg client.Message) error
	Update(ctx context.Context, msg client.Message) error
	Delete(ctx context.Context, msg client.Message) error
}

type Subscriber struct {
	svc SubscriberService
	log *wlog.Logger
}

func NewSubscriber(p *pubsub.Manager, log *wlog.Logger, svc SubscriberService) *Subscriber {
	h := &Subscriber{
		svc: svc,
		log: log,
	}
	p.AddOnConnect(func(channel *pubsub.Channel) error {
		var err error
		var delivery pubsub.Delivery
		const rejectExchange = "fts-reject"
		const exchange = "fts-stock"
		const queueName = "fts-stock"

		if err = channel.DeclareExchange(pubsub.Exchange{
			Name:    exchange,
			Type:    pubsub.ExchangeTypeDirect,
			Durable: true,
		}); err != nil {
			return err
		}

		if err = channel.DeclareExchange(pubsub.Exchange{
			Name:    rejectExchange,
			Type:    pubsub.ExchangeTypeTopic,
			Durable: true,
		}); err != nil {
			return err
		}

		if err = channel.DeclareDurableQueue(queueName, pubsub.Headers{
			"x-dead-letter-exchange": rejectExchange,
		}); err != nil {
			return err
		}

		if err = channel.DeclareDurableQueue(rejectExchange, nil); err != nil {
			return err
		}

		if err = channel.BindQueue(rejectExchange, "#", rejectExchange, pubsub.Headers{
			"x-expires": XDLExpire,
		}); err != nil {
			return err
		}

		if err = channel.BindQueue(queueName, client.MessageCreate, exchange, nil); err != nil {
			return err
		}
		if err = channel.BindQueue(queueName, client.MessageUpdate, exchange, nil); err != nil {
			return err
		}
		if err = channel.BindQueue(queueName, client.MessageDelete, exchange, nil); err != nil {
			return err
		}

		delivery, err = channel.ConsumeQueue(queueName, false)
		if err != nil {
			return err
		}

		go func() {
			var err error
			for {
				select {
				case msg, ok := <-delivery:
					if !ok {
						return
					}

					/* TODO
					if msg.ContentType != "text/json" {
						h.log.Warn("don't support context type "+msg.ContentType,
							wlog.String("content-type", msg.ContentType),
						)
						continue
					}
					*/

					var m client.Message
					err = json.Unmarshal(msg.Body, &m)
					if err != nil {
						h.log.Error(err.Error())
						msg.Reject(false)
						continue
					}

					rlog := h.log.With(
						wlog.Any("id", m.Id),
						wlog.String("object_name", m.ObjectName),
						wlog.Int64("domain_id", m.DomainId),
						wlog.String("action", msg.RoutingKey),
						wlog.Any(m.ObjectName, m.Body),
					)

					switch msg.RoutingKey {
					case client.MessageCreate:
						err = h.NewRecord(m)
					case client.MessageUpdate:
						err = h.UpdateRecord(m)
					case client.MessageDelete:
						err = h.DeleteRecord(m)
					default:
						rlog.Error("no handle routing key " + msg.RoutingKey)
					}

					if err == nil {
						msg.Ack(true)
						rlog.Debug("method " + msg.RoutingKey + " success")
					} else {
						msg.Reject(false)
						rlog.Error(err.Error(), wlog.Err(err))
					}

				}
			}
		}()

		return nil
	})

	return h
}

func (s *Subscriber) NewRecord(msg client.Message) error {
	return s.svc.Create(context.TODO(), msg)
}

func (s *Subscriber) UpdateRecord(msg client.Message) error {
	return s.svc.Update(context.TODO(), msg)
}

func (s *Subscriber) DeleteRecord(msg client.Message) error {
	return s.svc.Delete(context.TODO(), msg)
}
