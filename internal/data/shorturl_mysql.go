package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/flightzw/shorturl/internal/biz"
	"github.com/flightzw/shorturl/internal/data/model"
)

// repo *shorturlRepo github.com/flightzw/shorturl/internal/biz.ShorturlRepo
type shorturlRepo struct {
	data *Data
	log  *log.Helper
}

// NewShorturlRepo .
func NewShorturlRepo(data *Data, logger log.Logger) biz.ShorturlRepo {
	return &shorturlRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *shorturlRepo) CreateShorturl(ctx context.Context, data *model.Shorturl) (*model.Shorturl, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlRepo) UpdateShorturl(ctx context.Context, id int64) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlRepo) GetShorturl(ctx context.Context, id int64) (*model.Shorturl, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlRepo) ListShorturl(ctx context.Context) ([]*model.Shorturl, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *shorturlRepo) DeleteShorturl(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}
