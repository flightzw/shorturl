package service

import (
	"context"

	v1 "github.com/flightzw/shorturl/api/shorturl/v1"
	"github.com/flightzw/shorturl/internal/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShorturlService is a Shorturl service.
type ShorturlService struct {
	v1.UnimplementedShorturlServer

	uc *biz.ShorturlUsecase
}

// NewShorturlService new a Shorturl service.
func NewShorturlService(uc *biz.ShorturlUsecase) *ShorturlService {
	return &ShorturlService{uc: uc}
}
func (svc *ShorturlService) GetShorturl(ctx context.Context, req *v1.GetShorturlRequest) (*v1.GetShorturlReply, error) {
	data, err := svc.uc.GetShorturl(ctx, req.Longurl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &v1.GetShorturlReply{
		Shorturl: data,
	}, nil
}

func (svc *ShorturlService) GetLongurl(ctx context.Context, req *v1.GetLongurlRequest) (*v1.GetLongurlReply, error) {
	data, err := svc.uc.GetLongurl(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &v1.GetLongurlReply{
		Longurl: data,
	}, nil
}
