package main

import (
	"fmt"
	"io/ioutil"

	pb "github.com/hacktmz/GrpcMockService/pbs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:30001"
	file    = "./pbs/hello/hello.proto"
	file2   = "./pbs/hello/hello3.proto"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewMockServiceClient(conn)

	protoFileBytes, e := ioutil.ReadFile(file)
	if e != nil {
		fmt.Printf("ReadFile(%s): %s", file, e)
		return
	}
	// 调用方法
	req := &pb.ParserRequest{Protofile: string(protoFileBytes)}
	res, err := c.Startparser(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(res)

	protoFileBytes2, e := ioutil.ReadFile(file2)
	if e != nil {
		fmt.Printf("ReadFile(%s): %s", file, e)
		return
	}
	// 调用方法
	req2 := &pb.ParserRequest{Protofile: string(protoFileBytes2)}
	res2, err := c.Startparser(context.Background(), req2)

	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(res2)
	///////////////////////////////////////////////////////
	reqstart := &pb.MockRequest{
		Port: 30003,
		Headers: []*pb.ProtoHeader{
			{
				Filename:  "hello.proto",
				Protofile: string(protoFileBytes),
				Protojson: res.GetProtoformart(),
			},
			{
				Filename:  "hello3.proto",
				Protofile: string(protoFileBytes2),
				Protojson: res2.GetProtoformart(),
			},
		},
	}

	res3, err := c.Startmock(context.Background(), reqstart)
	fmt.Println(res3)
	///////////////////////////////////////////////////////
	/*
		reqstop := &pb.MockStopRequest{
			Ip: "127.0.0.1",
		}

		res3, err := c.Stop(context.Background(), reqstop)
		fmt.Println(res3)
	*/
}
