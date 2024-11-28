package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var (
	ErrRoomClosed = errors.New("room closed")

	ErrChatNotFound      = errors.New("room not found")
	ErrChatClosed        = errors.New("room closed")
	ErrChatAlreadyActive = errors.New("chat already active")
	ErrChatAlreadyExists = errors.New("chat already exists")

	ErrClientNotAvailable = errors.New("client not available")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	closeCheck     = 60 * time.Minute
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type MessageService interface {
}

type ChatRepository interface {
	GetChatByID(ctx context.Context, chatID string) (models.Chat, error)
	CreateChat(ctx context.Context, chat models.Chat) (string, error)
	UpdateChat(ctx context.Context, chat models.Chat) error
	DeleteChat(ctx context.Context, chatID string) error
}

type Service struct {
	ActiveChats map[string]*Room
	msgService  MessageService
	chatRepo    ChatRepository

	ctx context.Context
	mu  sync.RWMutex
}

func NewService(
	ctx context.Context,
	msgService MessageService,
	chatRepo ChatRepository,
) *Service {

	ctx, cancel := context.WithCancel(ctx)

	s := &Service{
		ActiveChats: map[string]*Room{},
		msgService:  msgService,
		chatRepo:    chatRepo,
		ctx:         ctx,
		mu:          sync.RWMutex{},
	}

	go s.stop(cancel)

	return s
}

func (s *Service) CreateChat(ctx context.Context, chat models.Chat) (string, error) {

	chat.ID = uuid.New().String()

	slog := slog.With(
		slog.String("chat_id", chat.ID),
		slog.String("chat_name", chat.Name),
	)

	slog.Debug(
		"creating chat",
	)

	chatID, err := s.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return "", err
	}

	slog.Debug(
		"creating room",
	)

	room := s.newRoom(chat)
	s.ActiveChats[chatID] = room

	go room.Run(s.ctx)

	slog.Debug(
		"room created",
	)

	return chatID, nil
}

func (s *Service) ConnectByID(
	w http.ResponseWriter,
	r *http.Request,
	chatID string,
	user *models.User,
) error {

	slog.With(
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("chat_id", chatID),
	)

	chat, err := s.chatRepo.GetChatByID(r.Context(), chatID)
	if err != nil {
		return fmt.Errorf("get chat by id: %w", err)
	}

	slog.Debug(
		"upgrading connection",
	)

	client := NewClient(*user)

	if err := s.connect(client, w, r); err != nil {
		return err
	}

	slog.Debug(
		"adding client to chat",
	)

	var room *Room

	s.mu.Lock()

	room, ok := s.ActiveChats[chatID]
	if !ok {
		room = s.newRoom(chat)
		s.ActiveChats[chatID] = room
	}

	s.mu.Unlock()

	room.Add(client)

	slog.Debug(
		"starting session",
	)

	if err := client.StartSession(r.Context(), client.conn, room); err != nil {
		room.Kick(client)
		return err
	}

	return nil
}

func (s *Service) connect(
	client *Client,
	w http.ResponseWriter,
	r *http.Request,
) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrading connection: %v", err)
	}

	client.addConnection(conn)

	return nil
}

func (s *Service) stop(cancel context.CancelFunc) {
	t := time.NewTicker(closeCheck)
	defer t.Stop()

	for {
		select {
		case <-s.ctx.Done():

			cancel()

			return
		case <-t.C:
			for _, r := range s.ActiveChats {
				if len(r.ActiveUsers) == 0 {
					r.Manager.Close <- struct{}{}
				}
			}
		}
	}
}
