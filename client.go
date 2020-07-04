package mongodb

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/Bezunca/mongo_connection/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalClient *mongo.Client = nil
var globalContext context.Context = nil
var globalConfigs *config.MongoConfigs = nil

func New(configs *config.MongoConfigs, tlsConfig *tls.Config) (*mongo.Client, error) {
	url, err := configs.Url()
	if err != nil {
		return nil, err
	}

	globalConfigs = configs
	globalClient, err = mongo.NewClient(
		options.Client().
			ApplyURI(url).
			SetAuth(configs.Credentials()).
			SetTLSConfig(tlsConfig),
		)
	if err != nil {
		return nil, err
	}

	globalContext = context.Background()

	if err = globalClient.Connect(globalContext); err != nil {
		return nil, err
	}

	return globalClient, nil
}

func Get() *mongo.Client {
	if globalClient == nil {
		panic("MongoDB client must be initialized before used!")
	}

	return globalClient
}

func Close() error {
	ctx, cancel := context.WithTimeout(globalContext, 5*time.Second)
	defer cancel()
	return Get().Disconnect(ctx)
}

func TimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(globalContext, globalConfigs.Timeout)
}