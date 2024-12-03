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

// Интерфейс сервиса работы с сообщениями
type MessageService interface {
	GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error)
	SaveMessages(ctx context.Context, messages []models.Message) error
}

// Интерфейс репозитория чатов
type ChatRepository interface {
	GetAllChats(ctx context.Context) ([]models.Chat, error)
	GetChatByID(ctx context.Context, chatID string) (models.Chat, error)
	CreateChat(ctx context.Context, chat models.Chat) (string, error)
	UpdateChat(ctx context.Context, chat models.Chat) error
	DeleteChat(ctx context.Context, chatID string) error
}

// Сервис работы с чатами
type Service struct {
	ActiveChats map[string]*Room
	msgService  MessageService
	chatRepo    ChatRepository

	ctx context.Context
	mu  sync.RWMutex
}

// Создание нового сервиса
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

// Создание нового чата и запуск его работы
func (s *Service) CreateChat(ctx context.Context, chat models.Chat) (string, error) {

	log := slog.With(
		slog.String("chat_id", chat.ID),
		slog.String("chat_name", chat.Name),
	)

	log.Debug("creating chat")

	chatID, err := s.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return "", fmt.Errorf("create chat: %w", err)
	}

	log.Debug("creating new room")

	room, err := s.newRoom(chat)
	if err != nil {
		return "", fmt.Errorf("create room: %w", err)
	}

	s.ActiveChats[chatID] = room

	go room.Run(s.ctx)

	log.Debug("room created")

	return chatID, nil
}

// Получение активных чатов из сервиса
func (s *Service) GetActiveChats(ctx context.Context) ([]models.Chat, error) {
	log := slog.Default()

	s.mu.RLock()
	defer s.mu.RUnlock()

	chats := make([]models.Chat, len(s.ActiveChats))

	i := 0

	start := time.Now()

	log.Debug("getting active chats")

	for _, room := range s.ActiveChats {
		chats[i] = room.Chat
		i++
	}

	log.Debug(
		"active chats got",
		slog.Duration("duration", time.Since(start)),
		slog.Int("chats_count", len(chats)),
	)

	return chats, nil
}

// Получение всех чатов из репозитория
func (s *Service) GetAllChats(ctx context.Context) ([]models.Chat, error) {
	log := slog.Default()

	log.Debug("getting all chats")

	chats, err := s.chatRepo.GetAllChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	log.Debug(
		"all chats got",
		slog.Int("chats_count", len(chats)),
	)

	return chats, nil
}

// Получение истории сообщений чата
func (s *Service) GetMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	log := slog.With(slog.String("chat_id", chatID))

	log.Debug("get chat from repository")

	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by id: %w", err)
	}

	log.Debug("getting chat messages")

	room, ok := s.ActiveChats[chatID]
	if !ok {
		log.Debug("search for innactive room")

		msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	log.Debug("search for active room")

	log.Debug("stashing messages")

	if err := room.StashHistory(ctx); err != nil {
		return nil, fmt.Errorf("stash history: %w", err)
	}

	msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
	if err != nil {
		return nil, err
	}

	log.Debug(
		"messages got",
		slog.Int("messages_count", len(msgs)),
	)

	return msgs, nil
}

// Подключает клиента к чату по ID
func (s *Service) ConnectByID(
	w http.ResponseWriter,
	r *http.Request,
	chatID string,
	user *models.User,
) error {

	log := slog.With(
		slog.String("chat_id", chatID),
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
	)

	log.Debug("getting chat from repository")

	chat, err := s.chatRepo.GetChatByID(r.Context(), chatID)
	if err != nil {
		return fmt.Errorf("get chat by id: %w", err)
	}

	log.Debug("upgrading connection")

	client := NewClient(*user)

	if err := s.connect(client, w, r); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	log.Debug("find room")

	var room *Room

	s.mu.Lock()

	room, ok := s.ActiveChats[chatID]
	if !ok {
		var err error

		log.Debug("not in active, create new room")

		room, err = s.newRoom(chat)
		if err != nil {
			return fmt.Errorf("new room: %w", err)
		}
		s.ActiveChats[chatID] = room

		go room.Run(s.ctx)

		log.Debug("room created")
	}

	s.mu.Unlock()

	go room.Add(client)

	client.ChatRoom = room

	if err := client.StartSession(r.Context()); err != nil {

		room.Kick(client)

		return fmt.Errorf("start session: %w", err)
	}

	return nil
}

// Апгрейд http-соединения к websocket
func (s *Service) connect(
	client *Client,
	w http.ResponseWriter,
	r *http.Request,
) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrade connection: %w", err)
	}

	client.addConnection(conn)

	return nil
}

// Останавливает работу сервиса
func (s *Service) stop(cancel context.CancelFunc) {
	log := slog.Default()

	t := time.NewTicker(closeCheck)
	defer t.Stop()

	for {
		select {
		case <-s.ctx.Done():

			log.Info("service stopped, closing rooms")

			cancel()

			return
		case <-t.C:
			log.Debug("check inactive rooms")

			for _, r := range s.ActiveChats {
				if len(r.ActiveUsers) == 0 {

					log.Debug(
						"found inactive room, closing",
						slog.String("room_id", r.ID),
					)

					r.Manager.Close <- struct{}{}
				}
			}
		}
	}
}
