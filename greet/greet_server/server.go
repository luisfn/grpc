package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/luisfn/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Long greet request received")

	msgs := []string{}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("Got %d messages", len(msgs))
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: strings.Join(msgs, " "),
			})
		}

		if err != nil {
			log.Fatalf("Failed to process request: %v", err)
			return err
		}

		log.Printf("Received: %v", msg)

		firstName := msg.GetGreeting().GetFirstName()

		msgs = append(msgs, fmt.Sprintf("hello %s", firstName))
	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("Starting BiDi Server")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Failed to process  BiDi request: %v", err)
			return err
		}

		log.Printf("Received request: %v\n", req)

		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()

		resp := &greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello %s %s", firstName, lastName),
		}

		err = stream.Send(resp)
		if err != nil {
			log.Fatalf("Failed to send BiDi response: %v", err)
			return err
		}
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

	reflection.Register(s)

	fmt.Println("Listening on tcp/0.0.0.0:50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
