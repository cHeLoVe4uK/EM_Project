package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Чат-комната
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

// Менеджер чат-комнат. Слушает события и управляет комнатой
type RoomManger struct {
	Add    chan *Client
	Logout chan *Client
	Kick   chan *Client

	Close chan struct{}
}

// Создает новую чат-комнату
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
		return nil, fmt.Errorf("create history: %w", err)
	}

	r.History = history

	go r.controlHistory(r.saveMsgsChan)

	return r, nil
}

// Запускает работу чат-комнаты
func (r *Room) Run(ctx context.Context) {

	log := slog.With(
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

			log.Info(
				"client joined room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)
		case client := <-r.Manager.Logout:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			log.Info(
				"client left room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

		case client := <-r.Manager.Kick:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			log.Info(
				"client kicked from room",
				slog.String("client_id", client.ID),
				slog.String("username", client.Username),
			)

		case msg := <-r.Broadcast:

			log.Debug("render message")

			data, err := msg.Render()
			if err != nil {
				log.Error(
					"failed to render message",
					slog.Any("error", err),
				)

				continue
			}

			pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
			if err != nil {
				log.Error(
					"failed to prepare message",
					slog.Any("error", err),
				)
				continue
			}

			log.Debug("add message to history")

			r.History.Add(*msg)

			for c := range r.ActiveUsers {

				log.Debug(
					"send message to client",
					slog.String("client_id", c.ID),
				)

				go func() {

					if err := c.Send(pm); err != nil {
						log.Error(
							"failed to send message",
							slog.Any("error", err),
						)
					}
				}()

			}
		}

	}
}

// Добавляет клиента в чат
func (r *Room) Add(client *Client) {
	log := slog.With(
		slog.String("room_id", r.ID),
		slog.String("client_id", client.ID),
		slog.String("username", client.Username),
	)

	r.Manager.Add <- client

	log.Debug("read message history")

	msgs := r.History.Read()
	if len(msgs) == 0 {
		log.Debug("history is empty")
		return
	}

	data, err := json.Marshal(msgs)
	if err != nil {
		log.Warn(
			"marshal history",
			slog.Any("error", err),
		)
		return
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		log.Warn(
			"prepare history",
			slog.Any("error", err),
		)
		return
	}

	log.Debug(
		"send history",
		slog.Int("messages_count", len(msgs)),
	)

	client.Send(pm)
}

// Убирает клиента из чата при выходе
func (r *Room) Logout(client *Client) {
	r.Manager.Logout <- client
}

// Исключает клиента из чата
func (r *Room) Kick(client *Client) {
	r.Manager.Kick <- client
}

// Останавливает работу чат-комнаты
func (r *Room) Stop() {

	log := slog.With(slog.String("room_id", r.ID))

	log.Info("recieve close signal")

	systemMsg := MessageDTO{
		ID:         uuid.NewString(),
		AuthorName: "System",
		ChatID:     r.ID,
		Content:    "Room closed",
		IsEdited:   false,
		CreatedAt:  time.Now().UTC(),
	}

	data, err := systemMsg.Render()
	if err != nil {
		log.Warn(
			"render system message",
			slog.Any("error", err),
		)
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		log.Warn(
			"prepare system message",
			slog.Any("error", err),
		)
		return
	}

	for client := range r.ActiveUsers {
		log := slog.With(slog.String("client_id", client.ID))

		log.Debug("closing connection")

		r.Manager.Kick <- client

		go func() {

			defer func() {
				if err := client.conn.Close(); err != nil {
					log.Error(
						"failed to close connection",
						slog.Any("error", err),
					)
				}
			}()

			_ = client.Send(pm)

		}()

	}

	r.saveMsgsChan <- struct{}{}

	// TODO: Save users state for chat to repo

}

// Воркер, который слушает события и сохраняет историю чата в репозиторий
func (r *Room) controlHistory(saveChan chan struct{}) {
	tick := time.NewTicker(10 * time.Minute)
	defer tick.Stop()

	ctx := context.TODO()

	log := slog.With(slog.String("room_id", r.ID))

	for {
		select {
		case <-tick.C:

			log.Debug("stashing history by timer")

			if err := r.StashHistory(ctx); err != nil {
				log.Error(
					"failed to stash history by timer",
					slog.Any("error", err),
				)
			}

		case <-saveChan:

			log.Debug("stashing history by save chan")

			retry := 0

			for {

				if err := r.StashHistory(ctx); err != nil {
					log.Error(
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

// Сохраняет все новые сообщения в истории чата в репозиторий
func (r *Room) StashHistory(ctx context.Context) error {
	log := slog.Default()

	msgs := ToMessageBatch(r.History.ReadNew())

	if err := r.msgService.SaveMessages(ctx, msgs); err != nil {
		return fmt.Errorf("save messages: %w", err)
	}

	r.History.MarkReaded()

	log.Info("history stashed successfully")

	return nil
}
