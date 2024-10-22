package data

import (
	"context"
	"demoserveice/internal/biz"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gitlab.cqrb.cn/shangyou_mic/kit/errors_ez"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) CreateGreeter(ctx context.Context, g *biz.Greeter) error {
	err := errors.New("serverce err")
	return errors_ez.Wap(err)
}

func (r *greeterRepo) UpdateGreeter(ctx context.Context, g *biz.Greeter) error {
	return nil
}
