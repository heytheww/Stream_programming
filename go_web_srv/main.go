package main

import (
	"fmt"
	pb "go_web/message"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRPCServer
}

// C2S
// 这个函数在哪一侧，stream参数的类型就是哪一侧，这里是server
func (s *server) C2S(stream pb.RPC_C2SServer) error {
	for {
		massage, err := stream.Recv()

		// 如果消息已经发完了
		if err == io.EOF {

			return stream.SendAndClose(&pb.C2SResponse{Message: "C2S服务端接收完毕"})
		}

		if err != nil {
			return err
		}

		// 服务端收到客户端发送的stream
		fmt.Println("C2S:S<-", massage.Message)
		massage = nil
	}
}

// S2C
func (s *server) S2C(in *pb.S2CRequest, stream pb.RPC_S2CServer) error {

	mg := "S2Cmessage"
	for i := 0; i < 10; i++ {
		mg = mg + strconv.Itoa(i)
		if err := stream.Send(&pb.S2CResponse{Message: mg}); err != nil {
			return err
		}
	}

	return nil
}

// CS

func (s *server) CS(stream pb.RPC_CSServer) error {
	// 一旦开始接收数据，就不停止，直到遇到EOF
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// 服务端收到客户端发送的stream
		fmt.Println("S2C:S<-C", in.Message)

		if err := stream.Send(&pb.CSResponse{Message: "SC:C<-S"}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:1234")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterRPCServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
