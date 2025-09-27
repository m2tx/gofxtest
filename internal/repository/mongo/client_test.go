package mongo_test

import (
	"context"
	"testing"

	"github.com/m2tx/gofxtest/internal/env"
	"github.com/m2tx/gofxtest/internal/repository/mongo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	ctx    context.Context
	client mongo.MongoClient
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	config, err := env.New[mongo.MongoConfig]()
	if err != nil {
		panic(err)
	}
	client, err = mongo.NewClient(ctx, config, zap.NewNop())
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	if code != 0 {
		panic(code)
	}
}

func TestMongoClient(t *testing.T) {
	t.Run("Connect and Disconnect", func(t *testing.T) {
		mongoConfig, err := env.New[mongo.MongoConfig]()
		assert.NoError(t, err)

		mongoConfig.Database = "testdb"

		client, err := mongo.NewClient(ctx, mongoConfig, zap.NewNop())
		assert.NotNil(t, client)
		assert.NoError(t, err)

		db := client.Database()
		assert.NotNil(t, db)
		assert.Equal(t, "testdb", db.Name())

		err = db.Drop(ctx)
		assert.NoError(t, err)

		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	})
}
