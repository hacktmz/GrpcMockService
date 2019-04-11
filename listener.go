package main

import (
	//"time"
	"context"
	//"io"
	//"io/ioutil"
	"fmt"
	"net"

	"log"

	//"golang.org/x/net/context"
	//proto "github.com/golang/protobuf/proto"
	//otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	//"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	//"gopkg.in/natefinch/lumberjack.v2"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/protoc-gen-go/descriptor"
	//"github.com/jhump/protoreflect/desc"
	//"github.com/jhump/protoreflect/dynamic"
	pb "github.com/hacktmz/GrpcMockService/pbs"
	"google.golang.org/grpc/reflection"
)

const (
	cfgport = ":30001"
)

type server struct{}

func (s *server) Startmock(ctx context.Context, in *pb.MockRequest) (*pb.MockResponse, error) {
	log.Println("Startmock in port %v ", in.Port)
	if len(in.Headers) < 1 {
		log.Println("headers is nil")
		return &pb.MockResponse{
			Message: "headers is nil",
		}, nil
	}

	ch <- in
	log.Println("send success %v", in.Port)
	select {
	case errStr := <-ch_return_err:
		log.Println("errStr= %s", errStr)
		return &pb.MockResponse{
			Message: errStr,
		}, nil
	}
}

func (s *server) Startparser(ctx context.Context, in *pb.ParserRequest) (*pb.ParserResponse, error) {
	protoFileBytes := in.Protofile
	b, e := Proto2Json([]byte(protoFileBytes))
	if e != nil {
		return nil, e
	}
	return &pb.ParserResponse{
		Protoformart: string(b),
		Error:        "ok",
	}, nil
}

func (s *server) Stopmock(ctx context.Context, in *pb.StopRequest) (*pb.MockResponse, error) {
	log.Println("Stopmock in %s", in.GetIp())
	ch_stop <- true
	select {
	case <-ch_notify_stop:
		return &pb.MockResponse{
			Message: "ok",
		}, nil
	}
}

func Listener() {

	lis, err := net.Listen("tcp", cfgport)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMockServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	grpclog.Println(fmt.Sprintf("Listen on %s", cfgport))
}
