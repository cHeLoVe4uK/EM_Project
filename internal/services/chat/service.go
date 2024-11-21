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
	mu          sync.RWMutex
}

func NewService(
	ctx context.Context,
	msgService MessageService,
	chatRepo ChatRepository,
) *Service {

	return &Service{
		ActiveChats: map[string]*Room{},
		msgService:  msgService,
		chatRepo:    chatRepo,
		mu:          sync.RWMutex{},
	}
}

func (s *Service) CreateChat(ctx context.Context, chat models.Chat) (string, error) {

	chatID, err := s.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return "", err
	}

	room := s.newRoom(chat)
	s.ActiveChats[chatID] = room

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

	return client.StartSession(r.Context(), client.conn, room)
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
