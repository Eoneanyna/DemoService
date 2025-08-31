package data

import (
	"context"
	"demoserveice/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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

func (r *greeterRepo) CreateGreeter(ctx context.Context, g *biz.Greeter) (error, biz.NewsService) {
	//news := db.Ecmsnews{
	//	Id:      1,
	//	ClassId: 1,
	//	Content: "hello word ",
	//}
	//servic := biz.NewsService{
	//	AirctId: news.Id,
	//	Content: news.Content,
	//}
	//err := errors.New("serverce err")
	return nil, biz.NewsService{}
}

func (r *greeterRepo) UpdateGreeter(ctx context.Context, g *biz.Greeter) error {
	return nil
}
