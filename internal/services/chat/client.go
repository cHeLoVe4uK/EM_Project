package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	models.User

	ChatRoom *Room

	conn    *websocket.Conn
	recieve chan []*websocket.PreparedMessage
}

func NewClient(user models.User) *Client {
	return &Client{
		User:    user,
		recieve: make(chan []*websocket.PreparedMessage),
	}
}

func (c *Client) StartSession(ctx context.Context) error {
	defer func() {
		c.conn.Close()
		c.conn = nil
	}()

	slog.Info(
		"starting new session",
		slog.String("client_id", c.ID),
		slog.String("username", c.Username),
		slog.String("room_id", c.ChatRoom.ID),
	)

	gr, ctx := errgroup.WithContext(ctx)

	gr.Go(c.read)
	gr.Go(c.write)

	if err := gr.Wait(); err != nil {
		slog.Info(
			"session closed",
			slog.String("client_id", c.ID),
		)
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

	slog.With(
		slog.String("client_id", c.ID),
		slog.String("username", c.Username),
		slog.String("room_id", c.ChatRoom.ID),
	)

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
					return fmt.Errorf("write message: %w", err)
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			slog.Debug(
				"ping client",
			)

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return ErrClientNotAvailable
			}
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

	slog.With(
		slog.String("client_id", c.ID),
		slog.String("username", c.Username),
		slog.String("room_id", c.ChatRoom.ID),
	)

	for {

		_, text, err := c.conn.ReadMessage()
		if err != nil {

			slog.Error(
				"read message",
				slog.Any("error", err),
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
			slog.Warn(
				"decoding message",
				slog.Any("error", err),
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
		c.conn.Close()
	}
}
