package chat_repository

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	chatsCollection = "chats"
)

var (
	ErrInvalidChatID = errors.New("invalid chat id")
	ErrChatNotFound  = errors.New("chat not found")
)

type Repository struct {
	collection *mongo.Collection
}

func NewChatsRepo(db *mongo.Database) (*Repository, error) {
	collection := db.Collection(chatsCollection)

	index := mongo.IndexModel{
		Keys:    bson.M{"chat_id": "text"},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return nil, err
	}

	return &Repository{
		collection: collection,
	}, nil
}

func (r *Repository) GetAllChats(ctx context.Context) ([]models.Chat, error) {
	var chats []Chat

	filter := bson.M{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return ToChatBatch(chats), nil
}

func (r *Repository) GetChatByID(ctx context.Context, chatID string) (models.Chat, error) {
	var chat Chat

	if err := r.collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&chat); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return models.Chat{}, ErrChatNotFound
		}

		return models.Chat{}, err
	}

	return ToChat(chat), nil
}

func (r *Repository) CreateChat(ctx context.Context, chat models.Chat) (string, error) {
	repoChat := FromChat(chat)

	_, err := r.collection.InsertOne(ctx, repoChat)
	if err != nil {
		return "", err
	}

	return chat.ID, nil
}

func (r *Repository) UpdateChat(ctx context.Context, chat models.Chat) error {

	repoChat := FromChat(chat)

	_, err := r.collection.UpdateOne(ctx, bson.M{"chat_id": repoChat.ChatID}, bson.M{"$set": bson.M{"name": repoChat.Name}})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteChat(ctx context.Context, chatID string) error {

	res, err := r.collection.DeleteOne(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrChatNotFound
	}

	return nil
}
