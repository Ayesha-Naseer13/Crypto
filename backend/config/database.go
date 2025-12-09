package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB(ctx context.Context) (*mongo.Database, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database("crypto_wallet_1"), nil
}

func DisconnectDB(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}
