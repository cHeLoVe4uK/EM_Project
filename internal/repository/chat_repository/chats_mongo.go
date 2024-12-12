package chat_repository

import (
	"context"
	"errors"
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	chatsCollection = "chats"
)

var (
	ErrInvalidChatID = errors.New("invalid chat id")
	ErrChatNotFound  = errors.New("chat with this id not found")
)

type ChatsRepo struct {
	collection *mongo.Collection
}

func NewChatsRepo(db *mongo.Database) *ChatsRepo {
	return &ChatsRepo{
		collection: db.Collection(chatsCollection),
	}
}

func (r *ChatsRepo) GetAllChats(ctx context.Context) ([]models.Chat, error) {
	var chats []models.Chat

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatsRepo) GetChatByID(ctx context.Context, chatID string) (models.Chat, error) {
	var chat models.Chat

	objectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return chat, ErrInvalidChatID
	}

	if err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&chat); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return chat, ErrChatNotFound
		}
		return chat, err
	}

	return chat, nil
}

func (r *ChatsRepo) CreateChat(ctx context.Context, chat models.Chat) (string, error) {
	res, err := r.collection.InsertOne(ctx, chat)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *ChatsRepo) UpdateChat(ctx context.Context, chat models.Chat) error {
	objectID, err := primitive.ObjectIDFromHex(chat.ID)
	if err != nil {
		return ErrInvalidChatID
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": chat.Name}})
	return err
}

func (r *ChatsRepo) DeleteChat(ctx context.Context, chatID string) error {
	objectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return ErrInvalidChatID
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
