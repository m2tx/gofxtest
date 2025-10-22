package mongo_test

import (
	"testing"
	"time"

	"github.com/m2tx/gofxtest/domain"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEventRepository_FindAll(t *testing.T) {
	rep := mongo.NewEventRepository(client, zap.NewNop())

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
	rep := mongo.NewEventRepository(client, zap.NewNop())

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
	rep := mongo.NewEventRepository(client, zap.NewNop())

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should insert event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestEventRepository_Delete(t *testing.T) {
	rep := mongo.NewEventRepository(client, zap.NewNop())

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should delete event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		err = rep.Delete(ctx, id)
		assert.NoError(t, err)

		res, err := rep.FindByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestEventRepository_Update(t *testing.T) {
	rep := mongo.NewEventRepository(client, zap.NewNop())

	defer client.Database().Collection("events").Drop(ctx)

	t.Run("should delete event successfully", func(t *testing.T) {
		event := domain.Event{
			Name: "Test Event",
			Date: &time.Time{},
		}

		id, err := rep.Insert(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		date := time.Now().UTC().Add(100 * time.Hour)
		event.ID = id
		event.Name = "Updated Event"
		event.Date = &date
		err = rep.Update(ctx, event)
		assert.NoError(t, err)

		res, err := rep.FindByID(ctx, id)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, "Updated Event", res.Name)
		assert.Equal(t, date.Day(), res.Date.Day())
		assert.Equal(t, date.Month(), res.Date.Month())
		assert.Equal(t, date.Year(), res.Date.Year())
		assert.Equal(t, date.Hour(), res.Date.Hour())
		assert.Equal(t, date.Minute(), res.Date.Minute())
		assert.Equal(t, date.Second(), res.Date.Second())
	})
}
