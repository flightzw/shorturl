package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"

	"github.com/flightzw/shorturl/internal/conf"
	"github.com/flightzw/shorturl/internal/data/model"
)

var (
// ErrUserNotFound is user not found.
)

// ShorturlRepo is a Greater repo.
type ShorturlRepo interface {
	CreateShorturl(ctx context.Context, data *model.Shorturl) (*model.Shorturl, error) // 创建
	UpdateShorturl(ctx context.Context, id int64)                                      // 更新
	GetShorturl(ctx context.Context, id int64) (*model.Shorturl, error)                // 按 id 查询
	ListShorturl(ctx context.Context) ([]*model.Shorturl, int64, error)                // 分页查询
	DeleteShorturl(ctx context.Context, id int64) error                                // 删除
}

// ShorturlUsecase is a Shorturl usecase.
type ShorturlUsecase struct {
	repo   ShorturlRepo
	log    *log.Helper
	config *conf.Data
}

// NewShorturlUsecase new a Shorturl usecase.
func NewShorturlUsecase(repo ShorturlRepo, logger log.Logger, config *conf.Data) *ShorturlUsecase {
	return &ShorturlUsecase{
		repo:   repo,
		log:    log.NewHelper(logger),
		config: config,
	}
}

func (svc *ShorturlUsecase) GetShorturl(ctx context.Context, longurl string) (url string, err error) {
	data, err := svc.repo.CreateShorturl(ctx, model.NewShorturl(longurl))
	if err != nil {
		return "", errors.Wrap(err, "repo.CreateShorturl")
	}
	code, err := decimalToBase62(data.ID)
	if err != nil {
		return "", errors.Wrap(err, "decimalToBase62")
	}
	return svc.config.Shorturl.UrlPrefix + code, nil
}

func (svc *ShorturlUsecase) GetLongurl(ctx context.Context, code string) (longurl string, err error) {
	id, err := base62ToDecimal(code)
	if err != nil {
		return "", errors.Wrap(err, "base62ToDecimal")
	}
	data, err := svc.repo.GetShorturl(ctx, id)
	if err != nil {
		return "", errors.Wrap(err, "repo.GetShorturl")
	}
	fmt.Println(data)
	return data.URL, nil
}
