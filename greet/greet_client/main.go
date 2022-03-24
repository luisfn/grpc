package main

import (
	"context"
	"github.com/luisfn/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)

	doStream(c)
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

func doStream(c greetpb.GreetServiceClient) {
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
