package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/meraiku/logging"
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

	roomCtx context.Context
	mu      sync.RWMutex
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
	ctx := context.Background()

	log := logging.WithAttrs(
		ctx,
		logging.String("room_id", r.ID),
	)

	ctx = logging.ContextWithLogger(ctx, log)

	r.roomCtx = ctx

	log.Debug("read message history")

	history, err := r.newHistory(s.ctx, r.ID)
	if err != nil {
		log.Error("create history", logging.Err(err))

		return nil, fmt.Errorf("create history: %w", err)
	}

	r.History = history

	go r.controlHistory(r.saveMsgsChan)

	log.Debug("message history readed")

	return r, nil
}

// Запускает работу чат-комнаты
func (r *Room) Run(ctx context.Context) {

	log := logging.L(r.roomCtx)

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
				logging.String("client_id", client.ID),
				logging.String("username", client.Username),
			)

		case client := <-r.Manager.Logout:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			log.Info(
				"client left room",
				logging.String("client_id", client.ID),
				logging.String("username", client.Username),
			)

		case client := <-r.Manager.Kick:

			r.mu.Lock()
			delete(r.ActiveUsers, client)
			r.mu.Unlock()

			log.Info(
				"client kicked from room",
				logging.String("client_id", client.ID),
				logging.String("username", client.Username),
			)

		case msg := <-r.Broadcast:

			log.Debug("render message")

			data, err := msg.Render()
			if err != nil {
				log.Error(
					"failed to render message",
					logging.Err(err),
				)

				continue
			}

			pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
			if err != nil {
				log.Error(
					"failed to prepare message",
					logging.Err(err),
				)
				continue
			}

			log.Debug("add message to history")

			r.History.Add(*msg)

			for c := range r.ActiveUsers {

				log.Debug(
					"send message to client",
					logging.String("client_id", c.ID),
				)

				go func() {

					if err := c.Send(pm); err != nil {
						log.Warn(
							"failed to send message",
							logging.Err(err),
						)
					}
				}()

			}
		}

	}
}

// Добавляет клиента в чат
func (r *Room) Add(client *Client) {
	log := logging.WithAttrs(
		r.roomCtx,
		logging.String("client_id", client.ID),
		logging.String("username", client.Username),
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
		log.Error(
			"marshal history",
			logging.Err(err),
		)
		return
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		log.Error(
			"prepare history",
			logging.Err(err),
		)
		return
	}

	log.Debug(
		"send history",
		logging.Int("messages_count", len(msgs)),
	)

	if err := client.Send(pm); err != nil {
		log.Warn(
			"send history to client",
			logging.Err(err),
		)
	}

	msg := fmt.Sprintf("%s has been joined chat!", r.Name)
	if err := r.SendSystemMessage(msg); err != nil {
		log.Warn(
			"send system message",
			logging.Err(err),
		)
	}

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

	log := logging.L(r.roomCtx)

	log.Info("recieve signal to close chat room")

	r.saveMsgsChan <- struct{}{}

	if err := r.SendSystemMessage("Chat room is closed!"); err != nil {
		log.Warn(
			"send system message",
			logging.Err(err),
		)
	}

	for client := range r.ActiveUsers {
		log := logging.WithAttrs(
			r.roomCtx,
			logging.String("client_id", client.ID),
		)

		log.Debug("close connection")

		r.Manager.Kick <- client

		go func() {
			if err := client.conn.Close(); err != nil {
				log.Warn(
					"close connection",
					logging.Err(err),
				)
			}
		}()
	}

	// TODO: Save users state for chat to repo

}

// Воркер, который слушает события и сохраняет историю чата в репозиторий
func (r *Room) controlHistory(saveChan chan struct{}) {
	tick := time.NewTicker(10 * time.Minute)
	defer tick.Stop()

	log := logging.L(r.roomCtx)

	for {
		select {
		case <-tick.C:

			log.Debug("stash history by timer")

			_ = r.StashHistory(r.roomCtx)

		case <-saveChan:

			retry := 0

			for {

				log.Debug("stash history by save chan", logging.Int("retry", retry))

				if err := r.StashHistory(r.roomCtx); err != nil {

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
	log := logging.L(ctx)

	msgs := ToMessageBatch(r.History.ReadNew())

	if err := r.msgService.SaveMessages(ctx, msgs); err != nil {
		return err
	}

	r.History.MarkReaded()

	log.Info("history stashed successfully")

	return nil
}

// Отправляет системное сообщение всем пользователям чата
func (r *Room) SendSystemMessage(msg string) error {
	systemMsg := MessageDTO{
		ID:         uuid.NewString(),
		AuthorName: "System",
		ChatID:     r.ID,
		Content:    msg,
		IsEdited:   false,
		CreatedAt:  time.Now().UTC(),
	}

	data, err := systemMsg.Render()
	if err != nil {
		return err
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	for client := range r.ActiveUsers {

		_ = client.Send(pm)
	}

	return nil
}

// Обновляет содержимое сообщения в базе и в истории, если это сообщение там находится
func (r *Room) UpdateMessage(ctx context.Context, msg models.Message) error {
	log := logging.L(ctx)

	log.Debug("stash history before update")

	if err := r.StashHistory(ctx); err != nil {
		return err
	}

	if err := r.msgService.UpdateMessageContent(ctx, msg); err != nil {
		return err
	}

	go func() {
		_ = r.History.UpdateMessage(ctx, msg)
	}()

	return nil
}

func (r *Room) DeleteMessage(ctx context.Context, msg models.Message) error {
	log := logging.L(ctx)

	log.Debug("stash history before delete")

	if err := r.StashHistory(ctx); err != nil {
		return err
	}

	if err := r.msgService.DeleteMessage(ctx, msg); err != nil {
		return err
	}

	go func() {
		log.Debug("delete message from history")

		_ = r.History.DeleteMessage(ctx, msg)
	}()

	return nil
}
