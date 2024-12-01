package chat

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	models.Chat

	ActiveUsers map[*Client]struct{}
	Users       map[string]bool

	Broadcast chan *MessageDTO
	History   *History

	saveMsgsChan chan struct{}

	Manager *RoomManger

	msgService MessageService
	mu         sync.RWMutex
}

type RoomManger struct {
	Add    chan *Client
	Logout chan *Client
	Kick   chan *Client

	Close chan struct{}
}

func (s *Service) newRoom(chat models.Chat) (*Room, error) {
	r := &Room{
		Chat:         chat,
		ActiveUsers:  make(map[*Client]struct{}),
		Users:        make(map[string]bool),
		Broadcast:    make(chan *MessageDTO, 100),
		saveMsgsChan: make(chan struct{}),
		msgService:   s.msgService,
		Manager: &RoomManger{
			Add:    make(chan *Client),
			Logout: make(chan *Client),
			Kick:   make(chan *Client),
			Close:  make(chan struct{}),
		},
	}

	history, err := r.newHistory(s.ctx, r.ID)
	if err != nil {
		return nil, err
	}

	r.History = history

	go r.controlHistory(r.saveMsgsChan)

	return r, nil
}

func (r *Room) Run(ctx context.Context) {

	slog.With(
		slog.String("room_id", r.ID),
	)

	for {
		select {
		case <-ctx.Done():

			r.Stop()

			return
		case <-r.Manager.Close:

			r.Stop()

			return
		case client := <-r.Manager.Add:

			r.mu.Lock()
			r.ActiveUsers[client] = struct{}{}
			r.mu.Unlock()

			slog.Info(
				"client joined room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)
		case client := <-r.Manager.Logout:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			slog.Info(
				"client left room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

		case client := <-r.Manager.Kick:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			slog.Info(
				"client kicked from room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

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

			pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
			if err != nil {
				slog.Error(
					"failed to prepare message",
					slog.Any("error", err),
				)

				continue
			}

			r.History.Add(*msg)

			for c := range r.ActiveUsers {

				go func() {

					if err := c.Send(pm); err != nil {
						slog.Error(
							"failed to send message",
							slog.Any("error", err),
						)
					}

				}()
			}
		}

	}
}

func (r *Room) Add(client *Client) {
	r.Manager.Add <- client

	msgs := r.History.Read()
	if len(msgs) == 0 {
		return
	}

	slog.Debug(
		"send history",
		slog.Int("messages", len(msgs)),
	)

	data, err := json.Marshal(msgs)
	if err != nil {
		slog.Warn(
			"marshal history",
			slog.Any("error", err),
		)
		return
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		slog.Warn(
			"prepare history",
			slog.Any("error", err),
		)
		return
	}

	client.Send(pm)
}

func (r *Room) Logout(client *Client) {
	r.Manager.Logout <- client
}

func (r *Room) Kick(client *Client) {
	r.Manager.Kick <- client
}

func (r *Room) Stop() {

	slog.Info(
		"closing room",
	)

	systemMsg := MessageDTO{
		ID:        uuid.NewString(),
		Author:    "System",
		ChatID:    r.ID,
		Content:   "Room closed",
		IsEdited:  false,
		Timestamp: time.Now(),
	}

	data, err := systemMsg.Render()
	if err != nil {
		slog.Warn(
			"render system message",
			slog.Any("error", err),
		)
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		slog.Warn(
			"prepare system message",
			slog.Any("error", err),
		)
		return
	}

	for client := range r.ActiveUsers {

		slog.Debug(
			"close connection",
			slog.String("client_id", client.ID),
		)

		r.Manager.Kick <- client

		go func() {

			defer func() {
				if err := client.conn.Close(); err != nil {
					slog.Error(
						"failed to close connection",
						slog.Any("error", err),
					)
				}
			}()

			client.Send(pm)

		}()

	}

	r.saveMsgsChan <- struct{}{}

	// TODO: Save users state for chat to repo

}

func (r *Room) controlHistory(saveChan chan struct{}) {
	tick := time.NewTicker(10 * time.Minute)
	defer tick.Stop()

	ctx := context.TODO()

	for {
		select {
		case <-tick.C:

			if err := r.StashHistory(ctx); err != nil {
				slog.Error(
					"failed to stash history by timer",
					slog.Any("error", err),
				)
			}

		case <-saveChan:
			retry := 0
			for {

				if err := r.StashHistory(ctx); err != nil {
					slog.Error(
						"failed to stash history by save chan",
						slog.Any("error", err),
						slog.Int("retry", retry),
					)

					retry++
					if retry > 3 {
						break
					}

					time.Sleep(time.Second)
					continue
				}

				break
			}

		}

	}
}

func (r *Room) StashHistory(ctx context.Context) error {
	msgs := ToMessageBatch(r.History.ReadNew())

	if err := r.msgService.SaveMessages(ctx, msgs); err != nil {
		return err
	}

	r.History.MarkReaded()

	return nil
}
