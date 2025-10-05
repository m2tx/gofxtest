package mongo

import (
	"context"

	"github.com/m2tx/gofxtest/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type eventRepository struct {
	client     MongoClient
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewEventRepository(client MongoClient, logger *zap.Logger) domain.EventRepository {
	collection := client.Database().Collection("events")

	return &eventRepository{
		client:     client,
		collection: collection,
		logger:     logger,
	}
}

func (e *eventRepository) Delete(ctx context.Context, id string) error {
	where := bson.M{"_id": id}

	_, err := e.collection.DeleteOne(ctx, where)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventRepository) FindAll(ctx context.Context) ([]domain.Event, error) {
	cursor, err := e.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			e.logger.Error("error closing cursor event repository", zap.Error(err))
		}
	}()

	var events []domain.Event
	for cursor.Next(ctx) {
		var event domain.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e *eventRepository) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	var event domain.Event
	err := e.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *eventRepository) Insert(ctx context.Context, event domain.Event) (string, error) {
	event.ID = primitive.NewObjectID().Hex()
	result, err := e.collection.InsertOne(ctx, event)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}

func (e *eventRepository) Update(ctx context.Context, event domain.Event) error {
	where := bson.M{"_id": event.ID}

	set := bson.M{
		"$set": event,
	}

	_, err := e.collection.UpdateOne(ctx, where, set)
	if err != nil {
		return err
	}

	return nil
}
