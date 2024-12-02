package chat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/meraiku/logging"

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

	log := logging.WithAttrs(
		ctx,
		logging.String("chat_id", chat.ID),
		logging.String("chat_name", chat.Name),
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

func (s *Service) GetActiveChats(ctx context.Context) ([]models.Chat, error) {
	log := logging.L(ctx)

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
		logging.Duration("duration", time.Since(start)),
		logging.Int("chats_count", len(chats)),
	)

	return chats, nil
}

func (s *Service) GetAllChats(ctx context.Context) ([]models.Chat, error) {
	log := logging.L(ctx)

	log.Debug("getting all chats")

	chats, err := s.chatRepo.GetAllChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	log.Debug(
		"all chats got",
		logging.Int("chats_count", len(chats)),
	)

	return chats, nil
}

func (s *Service) GetMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	log := logging.WithAttrs(
		ctx,
		logging.String("chat_id", chatID),
	)

	log.Debug("get chat from repository")

	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by id: %w", err)
	}

	log.Debug("getting chat messages")

	room, ok := s.ActiveChats[chatID]
	if !ok {
		log.Debug("direct search in repository")
		msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
		if err != nil {
			return nil, fmt.Errorf("get chat messages: %w", err)
		}

		log.Debug(
			"messages got",
			logging.Int("messages_count", len(msgs)),
		)

		return msgs, nil
	}

	log.Debug("stashing messages")

	if err := room.StashHistory(ctx); err != nil {
		return nil, fmt.Errorf("stash history: %w", err)
	}

	log.Debug("getting messages from repository")

	msgs, err := s.msgService.GetChatMessages(ctx, chat.ID)
	if err != nil {
		return nil, fmt.Errorf("get chat messages: %w", err)
	}

	log.Debug(
		"messages got",
		logging.Int("messages_count", len(msgs)),
	)

	return msgs, nil
}

func (s *Service) ConnectByID(
	w http.ResponseWriter,
	r *http.Request,
	chatID string,
	user *models.User,
) error {

	log := logging.WithAttrs(
		r.Context(),
		logging.String("chat_id", chatID),
		logging.String("user_id", user.ID),
		logging.String("username", user.Username),
	)

	ctx := logging.ContextWithLogger(r.Context(), log)

	r = r.WithContext(ctx)

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

func (s *Service) stop(cancel context.CancelFunc) {
	log := logging.L(s.ctx)

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
						logging.String("room_id", r.ID),
					)

					r.Manager.Close <- struct{}{}
				}
			}
		}
	}
}
