package chat

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

type Room struct {
	models.Chat
	Users map[*Client]bool

	Broadcast chan *MessageDTO
	History   [100]MessageDTO
	Manager   *RoomManger

	msgService MessageService
	mu         sync.RWMutex
}

type RoomManger struct {
	Add     chan *Client
	slogout chan *Client
	Kick    chan *Client
}

func (s *Service) newRoom(chat models.Chat) *Room {
	return &Room{
		Chat:       chat,
		Users:      map[*Client]bool{},
		Broadcast:  make(chan *MessageDTO, 100),
		msgService: s.msgService,
		Manager: &RoomManger{
			Add:     make(chan *Client),
			slogout: make(chan *Client),
			Kick:    make(chan *Client),
		},
	}
}

func (r *Room) Run(ctx context.Context) {
	slog.With(
		slog.String("room_id", r.ID),
	)

	for {
		select {
		case client := <-r.Manager.Add:
			r.Add(client)
		case client := <-r.Manager.slogout:
			r.slogout(client)
		case client := <-r.Manager.Kick:
			r.Kick(client)
		case msg := <-r.Broadcast:

			for c, ok := range r.Users {
				slog.With(
					slog.String("client_id", c.ID),
				)
				if ok {

					slog.Debug(
						"sending message",
						slog.Any("message", msg),
					)

					c.recieve <- msg.Render()
					continue
				}

				slog.Info(
					"client left room",
				)

				delete(r.Users, c)
				close(c.recieve)
			}
		}

	}
}

func (r *Room) Add(client *Client) {
	r.mu.Lock()
	r.Users[client] = true
	r.mu.Unlock()
}

func (r *Room) slogout(client *Client) {
	r.mu.Lock()
	r.Users[client] = false
	r.mu.Unlock()
}

func (r *Room) Kick(client *Client) {
	r.mu.Lock()
	delete(r.Users, client)
	r.mu.Unlock()
}
