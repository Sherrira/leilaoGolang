package mongodb

import (
	"context"
	"os"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "MONGODB_URL"
	MONGODB_DB  = "MONGODB_DB"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	mongoURL := os.Getenv(MONGODB_URL)
	mongoDatabase := os.Getenv(MONGODB_DB)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("Error connecting to MongoDB", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error pinging to MongoDB", err)
		return nil, err
	}

	return client.Database(mongoDatabase), nil
}

func NewMongoDBConnectionIntegrationTest(ctx context.Context) (*mongo.Database, error) {
	mongoURL := "mongodb://admin:admin@localhost:27017/auctions_test?authSource=admin"
	mongoDatabase := "auctions_test"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("Error connecting to MongoDB_TEST", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error pinging to MongoDB_TEST", err)
		return nil, err
	}

	return client.Database(mongoDatabase), nil
}
