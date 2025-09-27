package event_test

import (
	"context"
	"testing"

	"github.com/m2tx/gofxtest/domain"
	"github.com/m2tx/gofxtest/domain/event"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

func TestEventService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &mongo.EventRepositoryMock{
			InsertFn: func(ctx context.Context, event domain.Event) (string, error) {
				return "id", nil
			},
		}

		service := event.NewEventService(mock)

		id, err := service.Create(ctx, domain.Event{})
		assert.NoError(t, err)
		assert.Equal(t, 1, mock.InsertCount)
		assert.Equal(t, "id", id)
	})
	t.Run("Error", func(t *testing.T) {
		mock := &mongo.EventRepositoryMock{
			InsertFn: func(ctx context.Context, event domain.Event) (string, error) {
				return "", assert.AnError
			},
		}

		service := event.NewEventService(mock)

		id, err := service.Create(ctx, domain.Event{})
		assert.Error(t, err)
		assert.Equal(t, 1, mock.InsertCount)
		assert.Empty(t, id)
	})
}
