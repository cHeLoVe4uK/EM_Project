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

// Клиент чата
type Client struct {
	models.User

	ChatRoom *Room

	clientCtx context.Context

	conn    *websocket.Conn
	recieve chan []*websocket.PreparedMessage
}

// Создает нового клиента из юзера
func NewClient(user models.User) *Client {

	return &Client{
		User:    user,
		recieve: make(chan []*websocket.PreparedMessage),
	}
}

// Запускает чат сессию клиента
func (c *Client) StartSession(ctx context.Context) error {
	defer func() {
		c.conn.Close()
		c.conn = nil
	}()

	log := logging.L(ctx)

	c.clientCtx = ctx

	log.Info("start new session")

	gr, grctx := errgroup.WithContext(ctx)

	gr.Go(c.read)
	gr.Go(c.write)

	if err := gr.Wait(); err != nil {
		log.Info("session closed")
		grctx.Done()

		return err
	}

	return nil
}

// Добавляет веб сокет соединение клиенту
func (c *Client) addConnection(conn *websocket.Conn) {
	c.conn = conn
}

// Отправляет сообщения клиенту из чата
func (c *Client) read() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	log := logging.L(c.clientCtx)

	for {
		select {
		case msgs, ok := <-c.recieve:

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Warn(
						"write close message",
						logging.Err(err),
					)
				}

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
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			log.Debug("send ping")

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Warn(
					"send ping",
					logging.Err(err),
				)

				return ErrClientNotAvailable
			}

			log.Debug("got pong")
		}
	}
}

// Обрабатывает входящие сообщения от клиента
func (c *Client) write() error {

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	log := logging.L(c.clientCtx)

	for {
		_, text, err := c.conn.ReadMessage()
		if err != nil {

			log.Error(
				"read message",
				logging.Any("error", err),
			)
			break
		}

		msg := &MessageDTO{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Warn(
				"decoding message",
				logging.Any("error", err),
			)
			continue
		}

		c.ChatRoom.Broadcast <- NewMessage(c, msg.Content)
	}

	return nil
}

// Отаправляет сообщение клиенту
func (c *Client) Send(msg *websocket.PreparedMessage) error {
	if c.conn == nil {
		return ErrClientNotAvailable
	}

	c.recieve <- []*websocket.PreparedMessage{msg}

	return nil
}

// Отправляет несколько сообщений клиенту
func (c *Client) SendBatch(msgs []*websocket.PreparedMessage) error {
	if c.conn == nil {
		return ErrClientNotAvailable
	}

	c.recieve <- msgs

	return nil
}

// Закрывает веб сокет соединение у клиента
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
