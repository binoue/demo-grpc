package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/ymmt2005/demo-grpc/proto"
	"google.golang.org/grpc"
)

const address = "localhost:50052"

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTimeServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Report(context.Background(), &pb.ReportRequest{
		Interval: ptypes.DurationProto(3 * time.Second),
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		tr, err := r.Recv()
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Println(tr.GetMessage())
	}
}
