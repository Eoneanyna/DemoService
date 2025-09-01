package service

import (
	"context"
	v1 "demoserveice/api/news/v1"
	"demoserveice/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type NewsService struct {
	v1.UnimplementedNewsServiceServer

	uc  *biz.NewsUsecase
	log *log.Helper
}

func NewNewsService(uc *biz.NewsUsecase, logger log.Logger) *NewsService {
	return &NewsService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// GetNewsById 根据ID获取新闻详情
func (s *NewsService) GetNewsById(ctx context.Context, req *v1.GetNewsByIdRequest) (*v1.GetNewsByIdResponse, error) {
	log.Infof("GetNewsById req: %v", req)
	news, err := s.uc.GetNewsById(ctx, &biz.GetNewsByIdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &v1.GetNewsByIdResponse{
		Id:         news.Id,
		Title:      news.Title,
		Content:    news.Content,
		CreateTime: news.CreateTime,
	}, nil
}

// CreateNews 根据ID获取新闻详情
func (s *NewsService) CreateNews(ctx context.Context, req *v1.CreateNewsRequest) (*v1.CreateNewsResponse, error) {
	log.Infof("CreateNews req: %v", req)
	news, err := s.uc.CreateNews(ctx, biz.CreateNewsReq{Title: req.Title, Content: req.Content})
	if err != nil {
		return nil, err
	}

	return &v1.CreateNewsResponse{
		Id: news.Id,
	}, nil
}
