package domain

import (
	"context"
	"time"
)

type Event struct {
	ID   string     `bson:"_id,omitempty" json:"id,omitempty"`
	Name string     `bson:"name" json:"name"`
	Date *time.Time `bson:"date" json:"date"`
}

type EventService interface {
	GetAll(ctx context.Context) ([]Event, error)
	Get(ctx context.Context, id string) (*Event, error)
	Create(ctx context.Context, event Event) (string, error)
	Update(ctx context.Context, event Event) error
	Delete(ctx context.Context, id string) error
}

type EventRepository interface {
	FindAll(ctx context.Context) ([]Event, error)
	FindByID(ctx context.Context, id string) (*Event, error)
	Insert(ctx context.Context, event Event) (string, error)
	Update(ctx context.Context, event Event) error
	Delete(ctx context.Context, id string) error
}
