package log

import (
	"google.golang.org/grpc"
)

const ReceiveTopic = "log_service_receive"

var mqSender = &MqClient{}

type MqClient struct {
	client *grpc.ClientConn
}

func MqSender(mqConn *grpc.ClientConn) *MqClient {
	mqSender.client = mqConn
	return mqSender
}

//func (m *MqClient) Push(ctx context.Context, data *receive.LogRequest) error {
//	if mqSender.client == nil {
//		return errors.New("mq client is nil")
//	}
//	mq := pb.NewProduceClient(mqSender.client)
//	msg, _ := json.Marshal(data)
//	_, err := mq.Produce(ctx, &pb.ProduceRequest{
//		Topic: ReceiveTopic,
//		Msg:   string(msg),
//		Key:   fmt.Sprintf("%d:%s", data.SystemId, data.ObjectId),
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}
