package main

import (
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/ymmt2005/demo-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type timeServer struct{}

func (s timeServer) Report(req *pb.ReportRequest, stream pb.TimeService_ReportServer) error {
	duration := time.Second
	if dur := req.GetInterval(); dur != nil {
		d, err := ptypes.Duration(dur)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid duration: %s", dur.String())
		}
		duration = d
	}
	format := time.RFC3339
	switch req.GetFormat() {
	case pb.ReportRequest_RFC3339:
	case pb.ReportRequest_RFC822:
		format = time.RFC822
	case pb.ReportRequest_KITCHEN:
		format = time.Kitchen
	}

	tick := time.NewTicker(duration)
	defer tick.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-tick.C:
		}

		now := time.Now()
		pbt, _ := ptypes.TimestampProto(now)
		err := stream.Send(&pb.ReportResponse{
			Message:   now.Format(format),
			Timestamp: pbt,
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTimeServiceServer(s, &timeServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
