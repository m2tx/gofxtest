package mongo_test

import (
	"testing"
	"time"

	"github.com/m2tx/gofxtest/domain"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_FindAll(t *testing.T) {
	rep := mongo.NewEventRepository(client)

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should find event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)

		events, err := rep.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, id, events[0].ID)
	})
}

func TestEventRepository_FindByID(t *testing.T) {
	rep := mongo.NewEventRepository(client)

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should find event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)

		res, err := rep.FindByID(ctx, id)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, id, res.ID)
	})
}

func TestEventRepository_Insert(t *testing.T) {
	rep := mongo.NewEventRepository(client)

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should upsert event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}
