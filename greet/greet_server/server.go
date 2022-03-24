package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/luisfn/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet request received with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	msg := fmt.Sprintf("Hello %s %s", firstName, lastName)

	res := &greetpb.GreetResponse{
		Result: msg,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Many Times request received with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("%d %s %s", i, firstName, lastName)

		res := &greetpb.GreetManyTimesResponse{
			Result: msg,
		}

		err := stream.Send(res)

		if err != nil {
			log.Fatalf("Failed to stream message %s", msg)
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println("Listening on tcp/0.0.0.0:50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
