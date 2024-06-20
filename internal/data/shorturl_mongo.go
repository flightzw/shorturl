package data

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/flightzw/shorturl/internal/biz"
	"github.com/flightzw/shorturl/internal/data/model"
)

var _ biz.ShorturlRepo = (*shorturlMongoRepo)(nil)

var shorturlCacheKey = func(url string) string {
	return fmt.Sprintf("shorturl-cache-key:%X", md5.Sum([]byte(url)))
}

type count struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}

// repo *shorturlMongoRepo github.com/flightzw/shorturl/internal/biz.ShorturlRepo
type shorturlMongoRepo struct {
	data *Data
	log  *log.Helper
}

// NewShorturlMongoRepo .
func NewShorturlMongoRepo(data *Data, logger log.Logger) biz.ShorturlRepo {
	return &shorturlMongoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func getShorturlID(ctx context.Context, db *mongo.Database) (int64, error) {
	result := db.Collection("counts").
		FindOneAndUpdate(ctx,
			bson.D{{Key: "_id", Value: "shorturl_id"}},
			bson.D{{Key: "$inc", Value: bson.D{{Key: "value", Value: 1}}}},
		)
	if result.Err() != nil {
		return 0, errors.Wrap(result.Err(), "db.Collection.FindOneAndUpdate")
	}
	data := &count{}
	if err := result.Decode(data); err != nil {
		return 0, errors.Wrap(err, "result.Decode")
	}
	return data.Value + 1, nil
}

func (repo *shorturlMongoRepo) CreateShorturl(ctx context.Context, data *model.Shorturl) (*model.Shorturl, error) {
	cacheKey := shorturlCacheKey(data.URL)
	data1, err := repo.getShorturlFromCache(ctx, cacheKey)
	if data1 != nil && data1.URL == data.URL {
		return data1, nil
	}
	if err != nil {
		repo.log.Info("method", "repo.getShorturlFromCache", "error", err)
	}

	db := repo.data.mgoDB
	dataID, err := getShorturlID(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "getShorturlID")
	}
	data.ID = dataID
	_, err = db.Collection("shorturls").InsertOne(ctx, data)
	if err != nil {
		return nil, errors.Wrap(err, "db.Collection.InsertOne")
	}
	if err = repo.setShorturlCache(ctx, cacheKey, dataID); err != nil {
		repo.log.Warn("method", "repo.setShorturlCache", "error", err)
	}
	return data, nil
}

func (repo *shorturlMongoRepo) UpdateShorturl(ctx context.Context, id int64) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlMongoRepo) GetShorturl(ctx context.Context, id int64) (*model.Shorturl, error) {
	db := repo.data.mgoDB
	data := &model.Shorturl{}
	err := db.Collection("shorturls").FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(data)
	if err != nil {
		return nil, errors.Wrap(err, "db.Collection.FindOne")
	}
	return data, nil
}

func (repo *shorturlMongoRepo) ListShorturl(ctx context.Context) ([]*model.Shorturl, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlMongoRepo) DeleteShorturl(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlMongoRepo) getShorturlFromCache(ctx context.Context, key string) (*model.Shorturl, error) {
	idStr, err := repo.data.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis.Get")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "strconv.ParseInt")
	}
	data, err := repo.GetShorturl(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repo.GetShorturl")
	}
	return data, nil
}

func (repo *shorturlMongoRepo) setShorturlCache(ctx context.Context, key string, id int64) error {
	err := repo.data.redisClient.Set(ctx, key, id, 720*time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Set")
	}
	return nil
}
