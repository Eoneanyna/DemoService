package biz

import (
	"context"
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
	GetNewsById(ctx context.Context, req *GetNewsByIdReq) (GetNewsByIdResp, error)
	CreateNews(ctx context.Context, news *CreateNewsReq) (CreateNewsResp, error)
}

type NewsUsecase struct {
	repo NewsRepo
	log  *log.Helper
}

func NewNewsUsecase(repo NewsRepo, logger log.Logger) *NewsUsecase {
	return &NewsUsecase{repo: repo, log: log.NewHelper(logger)}
}

type GetNewsByIdReq struct {
	Id int32 `json:"id"`
}

type GetNewsByIdResp struct {
	Id         int32  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ViewCount  int64  `json:"view_count"`
	CreateTime int64  `json:"create_time"`
}

// GetNewsById 根据ID获取新闻详情
func (uc *NewsUsecase) GetNewsById(ctx context.Context, req *GetNewsByIdReq) (GetNewsByIdResp, error) {
	resp, err := uc.repo.GetNewsById(ctx, req)
	if err != nil {
		return GetNewsByIdResp{}, err
	}

	return resp, nil
}

type CreateNewsReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type CreateNewsResp struct {
	Id int32 `json:"id"`
}

// CreateNews 根据ID获取新闻详情
func (uc *NewsUsecase) CreateNews(ctx context.Context, req CreateNewsReq) (CreateNewsResp, error) {
	resp, err := uc.repo.CreateNews(ctx, &CreateNewsReq{Content: req.Content, Title: req.Title})
	if err != nil {
		return CreateNewsResp{}, err
	}

	return resp, nil
}
