package service

import (
	"fmt"
	pb "gitlab.cqrb.cn/shangyou_mic/testpg/api/demoserveice/v1"
	"strconv"
	"time"
)

type DataStreamService struct {
	pb.UnimplementedDataStreamServer
}

func NewDataStreamService() *DataStreamService {
	return &DataStreamService{}
}

func (s *DataStreamService) Subscribe(req *pb.SubscribeRequest, conn pb.DataStream_SubscribeServer) error {

	topic := req.GetTopic()
	fmt.Printf("Client subscribed to topic: %s", topic)

	for i := 0; i < 10; i++ {
		fmt.Println("sending data" + strconv.Itoa(i))
		err := conn.Send(&pb.SubscribeResponse{
			Data: fmt.Sprintf("Data %d from topic %s", i, topic),
		})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second) // 模拟实时推送
	}
	return nil

}
