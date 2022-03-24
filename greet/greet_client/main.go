package main

import (
	"context"
	"fmt"
	"github.com/luisfn/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	//doServerStream(c)
	doClientStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mambo",
			LastName:  "Jambo",
		},
	}
	resp, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("could not send request: %v", req)
	}

	log.Printf("Response: %v", resp)
}

func doServerStream(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Johnny",
			LastName:  "Leroy",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to send stream request")
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to process response: %v", err)
		}

		log.Println(msg)
	}
}

func doClientStream(c greetpb.GreetServiceClient) {

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream with server: %v", err)
	}

	names := []string{
		"John",
		"Reacher",
		"Belle",
		"Mambo",
		"Jambo",
	}

	for _, name := range names {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: name,
			},
		}

		err := stream.Send(req)

		if err != nil {
			log.Fatalf("Failed to stream request %v", req)
		}

		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Failed to get response from server: %v", err)
	}

	fmt.Println(resp.String())
}
