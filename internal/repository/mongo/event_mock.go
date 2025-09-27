package mongo

import (
	"context"

	"github.com/m2tx/gofxtest/domain"
)

func NewEventRepositoryMock() domain.EventRepository {
	return &EventRepositoryMock{}
}

type EventRepositoryMock struct {
	FindAllFn     func(ctx context.Context) ([]domain.Event, error)
	FindAllCount  int
	FindByIDFn    func(ctx context.Context, id string) (*domain.Event, error)
	FindByIDCount int
	InsertFn      func(ctx context.Context, event domain.Event) (string, error)
	InsertCount   int
	UpdateFn      func(ctx context.Context, event domain.Event) error
	UpdateCount   int
	DeleteFn      func(ctx context.Context, id string) error
	DeleteCount   int
}

func (mock *EventRepositoryMock) FindAll(ctx context.Context) ([]domain.Event, error) {
	mock.FindAllCount++
	return mock.FindAllFn(ctx)
}

func (mock *EventRepositoryMock) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	mock.FindByIDCount++
	return mock.FindByIDFn(ctx, id)
}

func (mock *EventRepositoryMock) Insert(ctx context.Context, event domain.Event) (string, error) {
	mock.InsertCount++
	return mock.InsertFn(ctx, event)
}

func (mock *EventRepositoryMock) Update(ctx context.Context, event domain.Event) error {
	mock.UpdateCount++
	return mock.UpdateFn(ctx, event)
}

func (mock *EventRepositoryMock) Delete(ctx context.Context, id string) error {
	mock.DeleteCount++
	return mock.DeleteFn(ctx, id)
}
