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
	recieve chan []byte
}

func NewClient(user models.User) *Client {
	return &Client{
		User:    user,
		recieve: make(chan []byte),
	}
}

func (c *Client) StartSession(ctx context.Context, conn *websocket.Conn, room *Room) error {
	defer func() {
		c.conn.Close()
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

	for {
		select {
		case msg, ok := <-c.recieve:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return ErrRoomClosed
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return fmt.Errorf("next writer: %v", err)
			}

			w.Write(msg)

			n := len(c.recieve)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				return fmt.Errorf("close writer: %v", err)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
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

	for {
		_, text, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}

		msg := &MessageDTO{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
		}

		c.ChatRoom.Broadcast <- NewMessage(c, msg.Content)
	}

	return nil
}
