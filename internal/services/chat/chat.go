package chat

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Room struct {
	models.Chat

	ActiveUsers map[*Client]struct{}
	Users       map[string]bool

	Broadcast chan *MessageDTO
	History   [100]MessageDTO
	Manager   *RoomManger

	msgService MessageService
	mu         sync.RWMutex
}

type RoomManger struct {
	Add    chan *Client
	Logout chan *Client
	Kick   chan *Client
}

func (s *Service) newRoom(chat models.Chat) *Room {
	return &Room{
		Chat:        chat,
		ActiveUsers: make(map[*Client]struct{}),
		Users:       make(map[string]bool),
		Broadcast:   make(chan *MessageDTO, 100),
		History:     [100]MessageDTO{},
		msgService:  s.msgService,
		Manager: &RoomManger{
			Add:    make(chan *Client),
			Logout: make(chan *Client),
			Kick:   make(chan *Client),
		},
	}
}

func (r *Room) Run(ctx context.Context) {

	slog.With(
		slog.String("room_id", r.ID),
	)

	for {
		select {
		case <-ctx.Done():

			slog.Info(
				"closing room",
			)
			for client := range r.ActiveUsers {

				slog.Debug(
					"close connection",
					slog.String("client_id", client.ID),
				)

				// TODO: Find better way to close connection

				r.Manager.Kick <- client

				// Is it good??
				client.conn.Close()

			}

			// TODO: Save all undelivered msges

			// TODO: Save users state for chat to repo

			return

		case client := <-r.Manager.Add:

			slog.Info(
				"client joined room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

			r.Add(client)
		case client := <-r.Manager.Logout:

			slog.Info(
				"client left room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

			r.Logout(client)
		case client := <-r.Manager.Kick:

			slog.Info(
				"client kicked from room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

			r.Kick(client)
		case msg := <-r.Broadcast:

			slog.Debug(
				"render message",
				slog.Any("message", msg),
			)

			data, err := msg.Render()
			if err != nil {
				slog.Error(
					"failed to render message",
					slog.Any("error", err),
				)

				continue
			}

			// TODO: Add to History

			for c := range r.ActiveUsers {
				slog.With(
					slog.String("client_id", c.ID),
				)

				slog.Debug(
					"sending message",
				)

				go func() {

					c.recieve <- data

				}()
			}
		}

	}
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	r.ActiveUsers[client] = struct{}{}
	r.mu.Unlock()
}

func (r *Room) Logout(client *Client) {
	r.mu.Lock()
	delete(r.ActiveUsers, client)
	r.mu.Unlock()
}

func (r *Room) Kick(client *Client) {
	r.mu.Lock()
	delete(r.ActiveUsers, client)
	r.mu.Unlock()
}
