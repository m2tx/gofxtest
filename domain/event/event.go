package event

import (
	"context"

	"github.com/m2tx/gofxtest/domain"
)

func NewEventService(repository domain.EventRepository) domain.EventService {
	return &eventService{
		repository: repository,
	}
}

type eventService struct {
	repository domain.EventRepository
}

func (s *eventService) GetAll(ctx context.Context) ([]domain.Event, error) {
	return s.repository.FindAll(ctx)
}

func (s *eventService) Get(ctx context.Context, id string) (*domain.Event, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *eventService) Create(ctx context.Context, event domain.Event) (string, error) {
	return s.repository.Insert(ctx, event)
}

func (s *eventService) Update(ctx context.Context, event domain.Event) error {
	return s.repository.Update(ctx, event)
}

func (s *eventService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
