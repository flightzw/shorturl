package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/flightzw/shorturl/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisClient, NewShorturlMongoRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	mgoDB       *mongo.Database
	redisClient *redis.Client
}

func NewRedisClient(conf *conf.Data, logger log.Logger) *redis.Client {
	log := log.NewHelper(log.With(logger, "module", "shorturl/data"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, redisClient *redis.Client) (*Data, func(), error) {
	mongoDB, cleanup, err := NewMongoDB("mongodb://"+c.Mongodb.Addr, c.Mongodb.DbName)
	if err != nil {
		return nil, nil, err
	}
	cleanup = func() {
		cleanup()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		mgoDB:       mongoDB,
		redisClient: redisClient,
	}, cleanup, nil
}

func NewMongoDB(url, dbName string) (*mongo.Database, func(), error) {
	logOptions := options.Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	clientOptions := options.Client().ApplyURI(url).SetLoggerOptions(logOptions)

	mgoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "mongo connect failed")
	}
	err = mgoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, nil, errors.Wrap(err, "mongo connect ping failed")
	}
	cleanup := func() {
		if err := mgoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
	return mgoClient.Database(dbName), cleanup, nil
}
