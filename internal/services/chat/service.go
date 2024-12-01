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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	closeCheck     = 60 * time.Minute
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type MessageService interface {
	GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error)
	SaveMessages(ctx context.Context, messages []models.Message) error
}

type ChatRepository interface {
	GetAllChats(ctx context.Context) ([]models.Chat, error)
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

	room, err := s.newRoom(chat)
	if err != nil {
		return "", err
	}

	s.ActiveChats[chatID] = room

	go room.Run(s.ctx)

	slog.Debug(
		"room created",
	)

	return chatID, nil
}

func (s *Service) GetActiveChats(ctx context.Context) ([]models.Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chats := make([]models.Chat, len(s.ActiveChats))

	i := 0

	for _, room := range s.ActiveChats {
		chats[i] = room.Chat
		i++
	}

	return chats, nil
}

func (s *Service) GetAllChats(ctx context.Context) ([]models.Chat, error) {
	chats, err := s.chatRepo.GetAllChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	return chats, nil
}

func (s *Service) GetMessages(ctx context.Context, chatID string) ([]models.Message, error) {

	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	room, ok := s.ActiveChats[chatID]
	if !ok {
		msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	if err := room.StashHistory(ctx); err != nil {
		return nil, err
	}

	msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *Service) ConnectByID(
	w http.ResponseWriter,
	r *http.Request,
	chatID string,
	user *models.User,
) error {

	log := slog.With(
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("chat_id", chatID),
	)

	chat, err := s.chatRepo.GetChatByID(r.Context(), chatID)
	if err != nil {
		return fmt.Errorf("get chat by id: %w", err)
	}

	log.Debug(
		"upgrading connection",
	)

	client := NewClient(*user)

	if err := s.connect(client, w, r); err != nil {
		return err
	}

	log.Debug(
		"adding client to chat",
	)

	var room *Room

	s.mu.Lock()

	room, ok := s.ActiveChats[chatID]
	if !ok {
		var err error

		log.Debug(
			"creating new room",
		)

		room, err = s.newRoom(chat)
		if err != nil {
			return err
		}
		s.ActiveChats[chatID] = room

		go room.Run(s.ctx)

		log.Debug(
			"room created",
		)
	}

	s.mu.Unlock()

	go room.Add(client)

	log.Debug(
		"starting session",
	)

	client.ChatRoom = room

	if err := client.StartSession(r.Context()); err != nil {
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
