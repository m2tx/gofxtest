package event_test

import (
	"context"
	"testing"

	"github.com/m2tx/gofxtest/domain"
	"github.com/m2tx/gofxtest/domain/event"
	"github.com/m2tx/gofxtest/internal/queue"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

func TestEventService_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			FindByIDFn: func(ctx context.Context, id string) (*domain.Event, error) {
				assert.Equal(t, "id", id)
				return &domain.Event{ID: id}, nil
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		event, err := service.Get(ctx, "id")
		assert.NoError(t, err)
		assert.Equal(t, 1, repositoryMock.FindByIDCount)
		assert.NotEmpty(t, event)
	})
	t.Run("Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			FindByIDFn: func(ctx context.Context, id string) (*domain.Event, error) {
				return nil, assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		event, err := service.Get(ctx, "id")
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.FindByIDCount)
		assert.Empty(t, event)
	})
}

func TestEventService_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			FindAllFn: func(ctx context.Context) ([]domain.Event, error) {
				return []domain.Event{{ID: "id"}}, nil
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		events, err := service.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, repositoryMock.FindAllCount)
		assert.Len(t, events, 1)
	})
	t.Run("Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			FindAllFn: func(ctx context.Context) ([]domain.Event, error) {
				return nil, assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		events, err := service.GetAll(ctx)
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.FindAllCount)
		assert.Len(t, events, 0)
	})
}

func TestEventService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			InsertFn: func(ctx context.Context, event domain.Event) (string, error) {
				return "id", nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventCreatedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return nil
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		id, err := service.Create(ctx, domain.Event{})
		assert.NoError(t, err)
		assert.Equal(t, 1, repositoryMock.InsertCount)
		assert.Equal(t, 1, queueMock.PublishCount)
		assert.Equal(t, "id", id)
	})
	t.Run("Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			InsertFn: func(ctx context.Context, event domain.Event) (string, error) {
				return "", assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		id, err := service.Create(ctx, domain.Event{})
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.InsertCount)
		assert.Empty(t, id)
	})
	t.Run("Publish Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			InsertFn: func(ctx context.Context, event domain.Event) (string, error) {
				return "id", nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventCreatedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		id, err := service.Create(ctx, domain.Event{})
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.InsertCount)
		assert.Equal(t, 1, queueMock.PublishCount)
		assert.Equal(t, "id", id)
	})
}

func TestEventService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventDeletedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return nil
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		err := service.Delete(ctx, "id")
		assert.NoError(t, err)
		assert.Equal(t, 1, repositoryMock.DeleteCount)
		assert.Equal(t, 1, queueMock.PublishCount)
	})
	t.Run("Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				return assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		err := service.Delete(ctx, "id")
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.DeleteCount)
	})
	t.Run("Publish Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventDeletedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		err := service.Delete(ctx, "id")
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.DeleteCount)
		assert.Equal(t, 1, queueMock.PublishCount)
	})
}

func TestEventService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			UpdateFn: func(ctx context.Context, event domain.Event) error {
				return nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventUpdatedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return nil
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		err := service.Update(ctx, domain.Event{ID: "id"})
		assert.NoError(t, err)
		assert.Equal(t, 1, repositoryMock.UpdateCount)
		assert.Equal(t, 1, queueMock.PublishCount)
	})
	t.Run("Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			UpdateFn: func(ctx context.Context, event domain.Event) error {
				return assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, nil)

		err := service.Update(ctx, domain.Event{ID: "id"})
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.UpdateCount)
	})
	t.Run("Publish Error", func(t *testing.T) {
		repositoryMock := &mongo.EventRepositoryMock{
			UpdateFn: func(ctx context.Context, event domain.Event) error {
				return nil
			},
		}

		queueMock := &queue.QueueMock{
			PublishFn: func(ctx context.Context, topic string, msg any) error {
				assert.Equal(t, event.EventUpdatedTopic, topic)
				e, ok := msg.(domain.Event)
				assert.True(t, ok)
				assert.Equal(t, "id", e.ID)

				return assert.AnError
			},
		}

		service := event.NewEventService(repositoryMock, queueMock)

		err := service.Update(ctx, domain.Event{ID: "id"})
		assert.Error(t, err)
		assert.Equal(t, 1, repositoryMock.UpdateCount)
		assert.Equal(t, 1, queueMock.PublishCount)
	})
}
