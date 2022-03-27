package main

import (
	"context"
	"log"

	"github.com/luisfn/blog/pb"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := pb.NewBlogServiceClient(cc)

	req := &pb.CreateBlogRequest{
		Blog: &pb.Blog{
			Author:  "Daniel Larusso",
			Title:   "Karate Kid 3",
			Content: "Marmelada",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to create blog entry %v", err)
	}
	log.Printf("Created blog response: %v\n", res)
}
