package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Database() *mongo.Database
}

type MongoConfig struct {
	URL      string `env:"MONGO_URL" required:"true"`
	Database string `env:"MONGO_DATABASE" required:"true"`
}

type mongoClient struct {
	config   MongoConfig
	client   *mongo.Client
	database *mongo.Database
	logger   *zap.Logger
}

func NewClient(ctx context.Context, config MongoConfig, logger *zap.Logger) (MongoClient, error) {
	client := &mongoClient{
		config: config,
		logger: logger,
	}

	err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *mongoClient) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.config.URL))
	if err != nil {
		return err
	}

	m.client = client

	m.logger.Info("connected to MongoDB", zap.String("url", m.config.URL), zap.String("database", m.config.Database))

	return nil
}

func (m *mongoClient) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	err := m.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	m.logger.Info("disconnected from MongoDB")

	return nil
}

func (m *mongoClient) Database() *mongo.Database {
	if m.client == nil {
		return nil
	}

	if m.database == nil {
		m.database = m.client.Database(m.config.Database)
	}

	return m.database
}
