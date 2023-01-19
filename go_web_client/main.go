package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	pb "go_web/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Reply struct {
	Per        int32
	Page_order int32
	Page_size  int32
}

func main() {
	// 基于证书的安全连接
	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 建议单个测试
	// C2S(conn)
	// S2C(conn)
	CS(conn)
}

// 客户端发送stream数据：C2S
func C2S(conn *grpc.ClientConn) {
	client := pb.NewRPCClient(conn)
	stream, err := client.C2S(context.Background())
	mg := "C2Smessage"
	for i := 0; i < 10; i++ {
		mg = mg + strconv.Itoa(i)
		if err := stream.Send(&pb.C2SRequest{Message: mg}); err != nil {
			log.Fatalf("%v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("C2S:C<-", reply.Message)
}

// 服务端端发送stream数据：S2C
func S2C(conn *grpc.ClientConn) {
	client := pb.NewRPCClient(conn)
	stream2, err := client.S2C(context.Background(), &pb.S2CRequest{Message: "S2Cmessage"})
	if err != nil {
	}
	for {
		in, err := stream2.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println("S2C:C<-", in.Message)
	}
}

func CS(conn *grpc.ClientConn) {
	client := pb.NewRPCClient(conn)
	// 互相发送stream数据：CS
	stream3, err := client.CS(context.Background())

	if err != nil {
		log.Fatalf("%v", err)
	}

	// 堵塞主协程，防止结束
	waitc := make(chan struct{})

	// 收发互不影响
	go func() {
		// 互相收发10条后结束收发

		for flag := 10; flag >= 1; flag-- {
			time.Sleep(time.Second)

			err := stream3.Send(&pb.CSRequest{Message: "SC:C->S" + strconv.Itoa(flag)})
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
		// 结束战斗
		close(waitc)
	}()

	go func() {
		for {
			time.Sleep(time.Second)

			in, err := stream3.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Println(in.Message)
		}
	}()

	// 堵塞主goroutinue
	<-waitc
}
