package service

import (
	"context"
	"demoserveice/internal/biz"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.cqrb.cn/shangyou_mic/testpg/api/demoserveice/v1"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	fmt.Println("SayHello Received: %v", in.GetName())
	err, data := s.uc.Create(ctx, &biz.Greeter{Hello: "222"})
	if err != nil {
		return nil, errors.New(500, err.Error(), err.Error())
	}
	//if in.GetName() == "error" {
	//	return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	return &v1.HelloReply{Message: "Hello " + data.Content}, nil
}
