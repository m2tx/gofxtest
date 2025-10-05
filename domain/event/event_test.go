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
}
