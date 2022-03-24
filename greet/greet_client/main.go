package main

import (
	"fmt"
	"github.com/luisfn/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client %f", c)

}
