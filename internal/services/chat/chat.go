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
		return nil, fmt.Errorf("create history: %w", err)
	}

	r.History = history

	go r.controlHistory(r.saveMsgsChan)

	return r, nil
}

func (r *Room) Run(ctx context.Context) {

	log := logging.WithAttrs(ctx, logging.String("room_id", r.ID))

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
						log.Error(
							"failed to send message",
							logging.Err(err),
						)
					}

				}()
			}
		}

	}
}

func (r *Room) Add(client *Client) {
	ctx := context.TODO()
	log := logging.WithAttrs(
		ctx,
		logging.String("room_id", r.ID),
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
		log.Warn(
			"marshal history",
			logging.Err(err),
		)
		return
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		log.Warn(
			"prepare history",
			logging.Err(err),
		)
		return
	}

	log.Debug(
		"send history",
		logging.Int("messages_count", len(msgs)),
	)

	client.Send(pm)
}

func (r *Room) Logout(client *Client) {
	r.Manager.Logout <- client
}

func (r *Room) Kick(client *Client) {
	r.Manager.Kick <- client
}

func (r *Room) Stop() {
	ctx := context.TODO()

	log := logging.WithAttrs(ctx, logging.String("room_id", r.ID))

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
			logging.Err(err),
		)
	}

	pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data)
	if err != nil {
		log.Warn(
			"prepare system message",
			logging.Err(err),
		)
		return
	}

	for client := range r.ActiveUsers {
		log := logging.WithAttrs(
			ctx,
			logging.String("client_id", client.ID),
		)

		log.Debug("closing connection")

		r.Manager.Kick <- client

		go func() {

			defer func() {
				if err := client.conn.Close(); err != nil {
					log.Error(
						"failed to close connection",
						logging.Err(err),
					)
				}
			}()

			_ = client.Send(pm)

		}()

	}

	r.saveMsgsChan <- struct{}{}

	// TODO: Save users state for chat to repo

}

func (r *Room) controlHistory(saveChan chan struct{}) {
	tick := time.NewTicker(10 * time.Minute)
	defer tick.Stop()

	ctx := context.TODO()

	log := logging.WithAttrs(ctx, logging.String("room_id", r.ID))

	for {
		select {
		case <-tick.C:

			log.Debug("stashing history by timer")

			if err := r.StashHistory(ctx); err != nil {
				log.Error(
					"failed to stash history by timer",
					logging.Err(err),
				)
			}

		case <-saveChan:

			log.Debug("stashing history by save chan")

			retry := 0

			for {

				if err := r.StashHistory(ctx); err != nil {
					log.Error(
						"failed to stash history by save chan",
						logging.Err(err),
						logging.Int("retry", retry),
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
	log := logging.L(ctx)

	msgs := ToMessageBatch(r.History.ReadNew())

	if err := r.msgService.SaveMessages(ctx, msgs); err != nil {
		return fmt.Errorf("save messages: %w", err)
	}

	r.History.MarkReaded()

	log.Info("history stashed successfully")

	return nil
}
