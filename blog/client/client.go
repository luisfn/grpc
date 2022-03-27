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

	//createBlog(c)
	readBlog(c, "623f9159909936423278264d")
}

func createBlog(c pb.BlogServiceClient) {
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

func readBlog(c pb.BlogServiceClient, id string) {
	req := &pb.ReadBlogRequest{
		Id: id,
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to read blog entry %v", err)
	}
	log.Printf("Read blog response: %v\n", res)
}
