package mongoimpl

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Configuration struct {
	Username string `envconfig:"MONGO_USER" default:"user"`
	Password string `envconfig:"MONGO_PASS" default:"pass"`
	URI      string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	Database string `envconfig:"MONGO_DATABASE" default:"test"`
}

type Client interface {
	FindSeos(ctx context.Context, pageId int) (MongoCursor, error)
	FindProducts(ctx context.Context, pageId int) (MongoCursor, error)
	CloseMongoClient() error
}

type MongoCursor interface {
	Next(ctx context.Context) bool
	All(ctx context.Context, results interface{}) error
	Decode(val interface{}) error
	Close(ctx context.Context) error
	Err() error
}

type ClientImpl struct {
	config      *Configuration
	mongoClient *mongo.Client
}

func InitMongoFromEnv() (*ClientImpl, error) {
	config, err := loadConfigurationFromEnv()
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return &ClientImpl{
		config:      config,
		mongoClient: mongoClient,
	}, nil
}

func loadConfigurationFromEnv() (*Configuration, error) {
	config := &Configuration{}
	err := envconfig.Process("", config)
	return config, err
}

func (c ClientImpl) FindSeos(ctx context.Context, pageId int) (MongoCursor, error) {
	return c.findInCollectionByPageId(ctx, pageId, "seos")
}

func (c ClientImpl) FindProducts(ctx context.Context, pageId int) (MongoCursor, error) {
	return c.findInCollectionByPageId(ctx, pageId, "products")

}

func (c ClientImpl) findInCollectionByPageId(ctx context.Context, pageId int, collection string) (MongoCursor, error) {
	return c.mongoClient.
		Database(c.config.Database).
		Collection(collection).
		Find(ctx, bson.D{{"page_id", pageId}})
}

func (c ClientImpl) CloseMongoClient() error {
	return c.mongoClient.Disconnect(context.TODO())
}
