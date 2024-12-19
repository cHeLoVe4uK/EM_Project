package mongo

import (
	"context"
	"errors"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/msg_repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	msgCollection = "messages"
)

type Repository struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) (*Repository, error) {
	collection := db.Collection(msgCollection)

	index := mongo.IndexModel{
		Keys:    bson.M{"id": "text"},
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

func (r *Repository) SaveMessages(ctx context.Context, message []models.Message) error {

	msgs := FromMessageBatch(message)

	_, err := r.collection.InsertMany(ctx, msgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetChatMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	var msgs []Message

	filter := bson.M{"chat_id": chatID}

	opt := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(100)

	cursor, err := r.collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}

	out := make([]Message, len(msgs))
	j := len(msgs) - 1
	for i := range msgs {
		out[i] = msgs[j]
		j--
	}

	return ToMessageBatch(out), nil
}

func (r *Repository) Update(ctx context.Context, msg models.Message) error {

	var repoMsg Message

	if err := r.collection.FindOne(ctx, bson.M{"id": msg.ID}).Decode(&repoMsg); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return msg_repository.ErrMessageNotFound
		}

		return err
	}

	if msg.AuthorID != repoMsg.AuthorID {
		return msg_repository.ErrNotAllowed
	}

	filter := bson.M{"id": msg.ID}

	update := bson.M{"$set": bson.M{"content": msg.Content, "is_edited": msg.IsEdited}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return msg_repository.ErrMessageNotFound
		}
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, msg models.Message) error {

	var repoMsg Message

	if err := r.collection.FindOne(ctx, bson.M{"id": msg.ID}).Decode(&repoMsg); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return msg_repository.ErrMessageNotFound
		}

		return err
	}

	if msg.AuthorID != repoMsg.AuthorID {
		return msg_repository.ErrNotAllowed
	}

	filter := bson.M{"id": msg.ID}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return msg_repository.ErrMessageNotFound
		}
		return err
	}

	return nil
}
