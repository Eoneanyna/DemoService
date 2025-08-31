package biz

//
//type Greeter struct {
//	Hello string
//}
//type NewsService struct {
//	AirctId int32
//	Content string
//}
//type GreeterRepo interface {
//}
//
//type GreeterUsecase struct {
//	repo GreeterRepo
//	log  *log.Helper
//}
//
//func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
//	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
//}
//
//func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) (error, NewsService) {
//
//	return uc.repo.CreateGreeter(ctx, g)
//}
//
//func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) error {
//	return uc.repo.UpdateGreeter(ctx, g)
//}
