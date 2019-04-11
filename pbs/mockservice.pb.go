// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mockservice.proto

/*
Package mock is a generated protocol buffer package.

It is generated from these files:
	mockservice.proto

It has these top-level messages:
	ParserRequest
	ParserResponse
	StopRequest
	MockRequest
	ProtoHeader
	MockResponse
*/
package mock

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ParserRequest struct {
	Protofile string `protobuf:"bytes,1,opt,name=protofile" json:"protofile,omitempty"`
}

func (m *ParserRequest) Reset()                    { *m = ParserRequest{} }
func (m *ParserRequest) String() string            { return proto.CompactTextString(m) }
func (*ParserRequest) ProtoMessage()               {}
func (*ParserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ParserRequest) GetProtofile() string {
	if m != nil {
		return m.Protofile
	}
	return ""
}

type ParserResponse struct {
	Protoformart string `protobuf:"bytes,1,opt,name=protoformart" json:"protoformart,omitempty"`
	Error        string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *ParserResponse) Reset()                    { *m = ParserResponse{} }
func (m *ParserResponse) String() string            { return proto.CompactTextString(m) }
func (*ParserResponse) ProtoMessage()               {}
func (*ParserResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ParserResponse) GetProtoformart() string {
	if m != nil {
		return m.Protoformart
	}
	return ""
}

func (m *ParserResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type StopRequest struct {
	Ip string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
}

func (m *StopRequest) Reset()                    { *m = StopRequest{} }
func (m *StopRequest) String() string            { return proto.CompactTextString(m) }
func (*StopRequest) ProtoMessage()               {}
func (*StopRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StopRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

// HelloRequest 请求结构
type MockRequest struct {
	Port    int32          `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	Headers []*ProtoHeader `protobuf:"bytes,2,rep,name=headers" json:"headers,omitempty"`
}

func (m *MockRequest) Reset()                    { *m = MockRequest{} }
func (m *MockRequest) String() string            { return proto.CompactTextString(m) }
func (*MockRequest) ProtoMessage()               {}
func (*MockRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MockRequest) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *MockRequest) GetHeaders() []*ProtoHeader {
	if m != nil {
		return m.Headers
	}
	return nil
}

type ProtoHeader struct {
	Filename  string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
	Protofile string `protobuf:"bytes,2,opt,name=protofile" json:"protofile,omitempty"`
	Protojson string `protobuf:"bytes,3,opt,name=protojson" json:"protojson,omitempty"`
}

func (m *ProtoHeader) Reset()                    { *m = ProtoHeader{} }
func (m *ProtoHeader) String() string            { return proto.CompactTextString(m) }
func (*ProtoHeader) ProtoMessage()               {}
func (*ProtoHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ProtoHeader) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *ProtoHeader) GetProtofile() string {
	if m != nil {
		return m.Protofile
	}
	return ""
}

func (m *ProtoHeader) GetProtojson() string {
	if m != nil {
		return m.Protojson
	}
	return ""
}

// HelloResponse 响应结构
type MockResponse struct {
	// bytes voice    = 1;
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *MockResponse) Reset()                    { *m = MockResponse{} }
func (m *MockResponse) String() string            { return proto.CompactTextString(m) }
func (*MockResponse) ProtoMessage()               {}
func (*MockResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *MockResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ParserRequest)(nil), "test.ParserRequest")
	proto.RegisterType((*ParserResponse)(nil), "test.ParserResponse")
	proto.RegisterType((*StopRequest)(nil), "test.StopRequest")
	proto.RegisterType((*MockRequest)(nil), "test.MockRequest")
	proto.RegisterType((*ProtoHeader)(nil), "test.ProtoHeader")
	proto.RegisterType((*MockResponse)(nil), "test.MockResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MockService service

type MockServiceClient interface {
	// 定义SayHello方法
	Startparser(ctx context.Context, in *ParserRequest, opts ...grpc.CallOption) (*ParserResponse, error)
	Startmock(ctx context.Context, in *MockRequest, opts ...grpc.CallOption) (*MockResponse, error)
	Stopmock(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*MockResponse, error)
}

type mockServiceClient struct {
	cc *grpc.ClientConn
}

func NewMockServiceClient(cc *grpc.ClientConn) MockServiceClient {
	return &mockServiceClient{cc}
}

func (c *mockServiceClient) Startparser(ctx context.Context, in *ParserRequest, opts ...grpc.CallOption) (*ParserResponse, error) {
	out := new(ParserResponse)
	err := grpc.Invoke(ctx, "/test.MockService/startparser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mockServiceClient) Startmock(ctx context.Context, in *MockRequest, opts ...grpc.CallOption) (*MockResponse, error) {
	out := new(MockResponse)
	err := grpc.Invoke(ctx, "/test.MockService/startmock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mockServiceClient) Stopmock(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*MockResponse, error) {
	out := new(MockResponse)
	err := grpc.Invoke(ctx, "/test.MockService/stopmock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MockService service

type MockServiceServer interface {
	// 定义SayHello方法
	Startparser(context.Context, *ParserRequest) (*ParserResponse, error)
	Startmock(context.Context, *MockRequest) (*MockResponse, error)
	Stopmock(context.Context, *StopRequest) (*MockResponse, error)
}

func RegisterMockServiceServer(s *grpc.Server, srv MockServiceServer) {
	s.RegisterService(&_MockService_serviceDesc, srv)
}

func _MockService_Startparser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MockServiceServer).Startparser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.MockService/Startparser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MockServiceServer).Startparser(ctx, req.(*ParserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MockService_Startmock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MockServiceServer).Startmock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.MockService/Startmock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MockServiceServer).Startmock(ctx, req.(*MockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MockService_Stopmock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MockServiceServer).Stopmock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.MockService/Stopmock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MockServiceServer).Stopmock(ctx, req.(*StopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MockService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "test.MockService",
	HandlerType: (*MockServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "startparser",
			Handler:    _MockService_Startparser_Handler,
		},
		{
			MethodName: "startmock",
			Handler:    _MockService_Startmock_Handler,
		},
		{
			MethodName: "stopmock",
			Handler:    _MockService_Stopmock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mockservice.proto",
}

func init() { proto.RegisterFile("mockservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x51, 0x4d, 0x4f, 0xc2, 0x40,
	0x10, 0x95, 0xf2, 0x3d, 0x45, 0x12, 0x46, 0x0e, 0x0d, 0xd1, 0x84, 0xec, 0x89, 0xc4, 0xc8, 0x01,
	0x3c, 0x79, 0xf4, 0x64, 0x4c, 0x34, 0xa4, 0xdc, 0xbc, 0x55, 0x1c, 0xb5, 0x62, 0xd9, 0x75, 0x67,
	0xf5, 0x97, 0xf9, 0x03, 0x4d, 0xf7, 0x03, 0x17, 0xe3, 0xad, 0xf3, 0xde, 0xbc, 0xce, 0xbe, 0xf7,
	0x60, 0x54, 0xc9, 0xcd, 0x96, 0x49, 0x7f, 0x95, 0x1b, 0x9a, 0x2b, 0x2d, 0x8d, 0xc4, 0x96, 0x21,
	0x36, 0xe2, 0x02, 0x8e, 0x57, 0x85, 0x66, 0xd2, 0x39, 0x7d, 0x7c, 0x12, 0x1b, 0x3c, 0x85, 0xbe,
	0xe5, 0x9f, 0xcb, 0x77, 0xca, 0x1a, 0xd3, 0xc6, 0xac, 0x9f, 0xff, 0x02, 0xe2, 0x16, 0x86, 0x61,
	0x9d, 0x95, 0xdc, 0x31, 0xa1, 0x80, 0x81, 0xa3, 0xa5, 0xae, 0x0a, 0x6d, 0xbc, 0xe4, 0x00, 0xc3,
	0x31, 0xb4, 0x49, 0x6b, 0xa9, 0xb3, 0xc4, 0x92, 0x6e, 0x10, 0x67, 0x90, 0xae, 0x8d, 0x54, 0xe1,
	0xf0, 0x10, 0x92, 0x52, 0x79, 0x79, 0x52, 0x2a, 0x71, 0x0f, 0xe9, 0x9d, 0xdc, 0x6c, 0x03, 0x8d,
	0xd0, 0x52, 0xd2, 0xff, 0xbf, 0x9d, 0xdb, 0x6f, 0x3c, 0x87, 0xee, 0x2b, 0x15, 0x4f, 0xa4, 0x39,
	0x4b, 0xa6, 0xcd, 0x59, 0xba, 0x18, 0xcd, 0x6b, 0x53, 0xf3, 0x55, 0x7d, 0xfc, 0xc6, 0x32, 0x79,
	0xd8, 0x10, 0x04, 0x69, 0x84, 0xe3, 0x04, 0x7a, 0xb5, 0xa3, 0x5d, 0x51, 0x05, 0x9b, 0xfb, 0xf9,
	0x30, 0x83, 0xe4, 0x4f, 0x06, 0x7b, 0xf6, 0x8d, 0xe5, 0x2e, 0x6b, 0x46, 0x6c, 0x0d, 0x88, 0x19,
	0x0c, 0xdc, 0xb3, 0x7d, 0x3e, 0x19, 0x74, 0x2b, 0x62, 0x2e, 0x5e, 0xc2, 0x99, 0x30, 0x2e, 0xbe,
	0x1b, 0xce, 0xe1, 0xda, 0xd5, 0x82, 0x57, 0x90, 0xb2, 0x29, 0xb4, 0x51, 0x36, 0x60, 0x3c, 0xf1,
	0x5e, 0xe2, 0x76, 0x26, 0xe3, 0x43, 0xd0, 0xdd, 0x10, 0x47, 0x78, 0x09, 0x7d, 0xab, 0xad, 0x6b,
	0x46, 0x9f, 0x42, 0x94, 0xde, 0x04, 0x63, 0x68, 0xaf, 0x5a, 0x42, 0x8f, 0x8d, 0x54, 0xb1, 0x28,
	0x6a, 0xe4, 0x7f, 0xd1, 0x75, 0xe7, 0xa1, 0x55, 0x0b, 0x1e, 0x3b, 0xd6, 0xf3, 0xf2, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x77, 0x10, 0x6b, 0x8b, 0x5b, 0x02, 0x00, 0x00,
}