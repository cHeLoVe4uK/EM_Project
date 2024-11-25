package chat_repository

import (
	"context"
	"errors"
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	chatsCollection = "chats"
)

var (
	ErrChatNotFound = errors.New("chat with this id not found")
)

type ChatsRepo struct {
	collection *mongo.Collection
}

func NewChatsRepo(db *mongo.Database) *ChatsRepo {
	return &ChatsRepo{
		collection: db.Collection(chatsCollection),
	}
}

func (r *ChatsRepo) GetChatByID(ctx context.Context, chatID string) (models.Chat, error) {
	var chat models.Chat

	objectID, _ := primitive.ObjectIDFromHex(chatID)
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&chat); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return chat, ErrChatNotFound
		}
		return chat, err
	}

	return chat, nil
}

func (r *ChatsRepo) CreateChat(ctx context.Context, chat models.Chat) (string, error) {
	res, err := r.collection.InsertOne(ctx, bson.M{"name": chat.Name})
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func (r *ChatsRepo) UpdateChat(ctx context.Context, chat models.Chat) error {
	objectID, _ := primitive.ObjectIDFromHex(chat.ID)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": chat.Name}})
	return err
}

func (r *ChatsRepo) DeleteChat(ctx context.Context, chatID string) error {
	objectID, _ := primitive.ObjectIDFromHex(chatID)
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrChatNotFound
	}

	return nil
}
