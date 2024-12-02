package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/gorilla/websocket"
	"github.com/meraiku/logging"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	models.User

	ChatRoom *Room

	conn    *websocket.Conn
	recieve chan []*websocket.PreparedMessage

	ctx context.Context
}

func NewClient(user models.User) *Client {
	ctx := context.Background()

	return &Client{
		User:    user,
		recieve: make(chan []*websocket.PreparedMessage),
		ctx:     ctx,
	}
}

func (c *Client) StartSession(ctx context.Context) error {
	defer func() {
		c.conn.Close()
		c.conn = nil
	}()

	log := logging.L(ctx)

	c.ctx = ctx

	log.Info("start new session")

	gr, ctx := errgroup.WithContext(ctx)

	gr.Go(c.read)
	gr.Go(c.write)

	if err := gr.Wait(); err != nil {
		log.Info("session closed")

		return err
	}

	return nil
}

func (c *Client) addConnection(conn *websocket.Conn) {
	c.conn = conn
}

func (c *Client) read() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	log := logging.L(c.ctx)

	for {
		select {
		case msgs, ok := <-c.recieve:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return ErrRoomClosed
			}

			for _, msg := range msgs {
				if err := c.conn.WritePreparedMessage(msg); err != nil {
					log.Warn(
						"write message",
						logging.Err(err),
					)
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			log.Debug("send ping")

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("send ping", logging.Err(err))
				return ErrClientNotAvailable
			}

			log.Debug("ping sent")
		}
	}
}

func (c *Client) write() error {

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	log := logging.L(c.ctx)

	for {
		_, text, err := c.conn.ReadMessage()
		if err != nil {

			log.Error(
				"read message",
				logging.Err(err),
			)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}

		msg := &MessageDTO{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Warn(
				"decoding message",
				logging.Err(err),
			)
			continue
		}

		c.ChatRoom.Broadcast <- NewMessage(c, msg.Content)
	}

	return nil
}

func (c *Client) Send(msg *websocket.PreparedMessage) error {
	if c.conn == nil {
		return ErrClientNotAvailable
	}

	c.recieve <- []*websocket.PreparedMessage{msg}

	return nil
}

func (c *Client) SendBatch(msgs []*websocket.PreparedMessage) error {
	if c.conn == nil {
		return ErrClientNotAvailable
	}

	c.recieve <- msgs

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		logging.L(c.ctx).Info("closing client connection")
		c.conn.Close()
	}
}
