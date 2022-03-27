package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/luisfn/blog/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	collection *mongo.Collection
)

type blogItem struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Author  string             `bson:"author"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

type server struct{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: "root",
		Password: "pass",
	}

	log.Println("Connecting Mongo DB")
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").
		SetAuth(credential)

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	collection = client.Database("db").Collection("blog")

	defer func() {
		log.Println("Disconnecting Mongo DB")
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		log.Println("Listening on tcp/0.0.0.0:50051")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Blocking until interrupt
	<-ch

	s.Stop()
	lis.Close()
	log.Println("Server Stopped")
}

func (*server) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	log.Printf("Creating blog entry: %v", req)

	blog := req.GetBlog()
	data := blogItem{
		Author:  blog.GetAuthor(),
		Title:   blog.GetTitle(),
		Content: blog.GetContent(),
	}

	dbRes, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := dbRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Can not parte to OID: %v", err))
	}

	return &pb.CreateBlogResponse{
		Blog: &pb.Blog{
			Id:      oid.Hex(),
			Author:  blog.GetAuthor(),
			Title:   blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	log.Printf("Retrieving blog entry: %v", req)

	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	data := &blogItem{}
	if err := collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Entry %s not found: %v", id, err))
	}

	return &pb.ReadBlogResponse{
		Blog: &pb.Blog{
			Id:      data.Id.Hex(),
			Author:  data.Author,
			Title:   data.Title,
			Content: data.Content,
		},
	}, nil
}
