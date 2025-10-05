package event

import (
	"context"

	"github.com/m2tx/gofxtest/domain"
	"github.com/m2tx/gofxtest/internal/queue"
)

const (
	EventCreatedTopic = "event-created"
	EventUpdatedTopic = "event-updated"
	EventDeletedTopic = "event-deleted"
)

func NewEventService(repository domain.EventRepository, publisher queue.Publisher) domain.EventService {
	return &eventService{
		repository: repository,
		publisher:  publisher,
	}
}

type eventService struct {
	repository domain.EventRepository
	publisher  queue.Publisher
}

func (s *eventService) GetAll(ctx context.Context) ([]domain.Event, error) {
	return s.repository.FindAll(ctx)
}

func (s *eventService) Get(ctx context.Context, id string) (*domain.Event, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *eventService) Create(ctx context.Context, event domain.Event) (string, error) {
	id, err := s.repository.Insert(ctx, event)
	if err != nil {
		return "", err
	}

	event.ID = id

	err = s.publisher.Publish(ctx, EventCreatedTopic, event)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *eventService) Update(ctx context.Context, event domain.Event) error {
	err := s.repository.Update(ctx, event)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, EventUpdatedTopic, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, EventDeletedTopic, domain.Event{ID: id})
	if err != nil {
		return err
	}

	return nil
}
