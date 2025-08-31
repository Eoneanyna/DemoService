package biz

import (
	"context"
	"demoserveice/internal/data"
	"demoserveice/internal/data/models/db"
	"github.com/go-kratos/kratos/v2/log"
)

type News struct {
	Id         int32
	Title      string
	Content    string
	ViewCount  int64
	CreateTime int64
}

type NewsRepo interface {
	GetNewsById(ctx context.Context, id int32) (db.News, error)
	CreateNews(ctx context.Context, news *data.CreateNewsReq) (data.CreateNewsResp, error)
}

type NewsUsecase struct {
	repo NewsRepo
	log  *log.Helper
}

func NewNewsUsecase(repo NewsRepo, logger log.Logger) *NewsUsecase {
	return &NewsUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetNewsById 根据ID获取新闻详情
func (uc *NewsUsecase) GetNewsById(ctx context.Context, id int32) (*News, error) {
	news, err := uc.repo.GetNewsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &News{
		Id:         news.Id,
		Title:      news.Title,
		Content:    news.Content,
		ViewCount:  news.ViewCount,
		CreateTime: news.CreateTime.Unix(),
	}, nil
}

// CreateNews 根据ID获取新闻详情
func (uc *NewsUsecase) CreateNews(ctx context.Context, req *data.CreateNewsReq) (data.CreateNewsResp, error) {
	news, err := uc.repo.CreateNews(ctx, &data.CreateNewsReq{Content: req.Content, Title: req.Title})
	if err != nil {
		return data.CreateNewsResp{}, err
	}

	return news, nil
}
