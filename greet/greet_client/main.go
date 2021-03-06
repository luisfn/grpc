package main

import (
	"context"
	"fmt"
	"github.com/luisfn/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type Person struct {
	FirstName string
	LastName  string
}

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	//doServerStream(c)
	//doClientStream(c)
	doBiDiStream(c)
}

func doBiDiStream(c greetpb.GreetServiceClient) {
	log.Println("Starting to send BiDi requests")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("could not create BiDi stream: %v", err)
		return
	}

	ch := make(chan struct{})

	people := []Person{
		{
			FirstName: "Johnny",
			LastName:  "Leroy",
		},
		{
			FirstName: "Jack",
			LastName:  "Reacher",
		},
		{
			FirstName: "Aloy",
			LastName:  "Rost",
		},
		{
			FirstName: "Pepeka",
			LastName:  "Ruiva",
		},
		{
			FirstName: "Talannah",
			LastName:  "Roth",
		},
	}

	//keep sending
	go func() {
		for _, person := range people {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: person.FirstName,
					LastName:  person.LastName,
				},
			}

			log.Printf("Sending message %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	//keep receiving
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("could not receive BiDi request: %v", err)
				break
			}

			log.Printf("Received: %v\n", res.GetResult())
		}

		close(ch)
	}()

	<-ch
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
