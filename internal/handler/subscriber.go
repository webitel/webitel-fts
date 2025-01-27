package handler

import (
	"context"
	"encoding/json"
	"github.com/webitel/webitel-fts/infra/pubsub"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type SubscriberService interface {
	Create(ctx context.Context, msg model.Message) error
	Update(ctx context.Context, msg model.Message) error
	Delete(ctx context.Context, msg model.Message) error
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
		const exchange = "fts-stock"
		const queueName = "fts-stock"

		if err = channel.DeclareExchange(pubsub.Exchange{
			Name:    exchange,
			Type:    pubsub.ExchangeTypeDirect,
			Durable: true,
		}); err != nil {
			return err
		}
		if err = channel.DeclareDurableQueue(queueName, nil); err != nil {
			return err
		}

		if err = channel.BindQueue(queueName, model.MessageCreate, exchange, nil); err != nil {
			return err
		}
		if err = channel.BindQueue(queueName, model.MessageUpdate, exchange, nil); err != nil {
			return err
		}
		if err = channel.BindQueue(queueName, model.MessageDelete, exchange, nil); err != nil {
			return err
		}

		ch, err := channel.ConsumeQueue(queueName, false)
		if err != nil {
			return err
		}

		go func() {
			var err error
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						return
					}

					if msg.ContentType != "text/json" {
						h.log.Warn("don't support context type "+msg.ContentType,
							wlog.String("content-type", msg.ContentType),
						)
						continue
					}

					var m model.Message
					err = json.Unmarshal(msg.Body, &m)
					if err != nil {
						h.log.Error(err.Error())
						continue
					}

					rlog := h.log.With(
						wlog.Any("id", m.Id),
						wlog.String("object_name", m.ObjectName),
						wlog.Int64("domain_id", m.DomainId),
					)

					switch msg.RoutingKey {
					case model.MessageCreate:
						err = h.NewRecord(m)
					case model.MessageUpdate:
						err = h.UpdateRecord(m)
					case model.MessageDelete:
						err = h.DeleteRecord(m)
					default:
						rlog.Error("no handle routing key " + msg.RoutingKey)
					}

					if err == nil {
						msg.Ack(true)
						rlog.Debug("method " + msg.RoutingKey + " success")
					} else {
						rlog.Error(err.Error(), wlog.Err(err))
					}

				}
			}
		}()

		return nil
	})

	return h
}

func (s *Subscriber) NewRecord(msg model.Message) error {
	return s.svc.Create(context.TODO(), msg)
}

func (s *Subscriber) UpdateRecord(msg model.Message) error {
	return s.svc.Update(context.TODO(), msg)
}

func (s *Subscriber) DeleteRecord(msg model.Message) error {
	return s.svc.Delete(context.TODO(), msg)
}
