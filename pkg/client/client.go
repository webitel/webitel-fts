package client

import (
	"github.com/webitel/webitel-fts/internal/model"
)

const exchange = "fts-stock"

type Publisher interface {
	Send(exchange string, rk string, body []byte) error
}

type Client struct {
	publisher Publisher
}

func New(p Publisher) *Client {
	return &Client{
		publisher: p,
	}
}

func (c *Client) Create(domainId int64, objectName string, id any, row any) error {
	msg, err := model.NewMessageJSON(domainId, objectName, id, row)
	if err != nil {
		return err
	}

	return c.publisher.Send(exchange, model.MessageCreate, msg)
}

func (c *Client) Update(domainId int64, objectName string, id any, row any) error {
	msg, err := model.NewMessageJSON(domainId, objectName, id, row)
	if err != nil {
		return err
	}

	return c.publisher.Send(exchange, model.MessageUpdate, msg)
}

func (c *Client) Delete(domainId int64, objectName string, id any) error {
	msg, err := model.NewMessageJSON(domainId, objectName, id, nil)
	if err != nil {
		return err
	}

	return c.publisher.Send(exchange, model.MessageDelete, msg)
}
