// Code generated by protoc-gen-go. DO NOT EDIT.
// source: desc_test_complex.proto

package testprotos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/protoc-gen-go/descriptor"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Test_Nested__NestedNested_EEE int32

const (
	Test_Nested__NestedNested_OK Test_Nested__NestedNested_EEE = 0
	Test_Nested__NestedNested_V1 Test_Nested__NestedNested_EEE = 1
	Test_Nested__NestedNested_V2 Test_Nested__NestedNested_EEE = 2
	Test_Nested__NestedNested_V3 Test_Nested__NestedNested_EEE = 3
	Test_Nested__NestedNested_V4 Test_Nested__NestedNested_EEE = 4
	Test_Nested__NestedNested_V5 Test_Nested__NestedNested_EEE = 5
	Test_Nested__NestedNested_V6 Test_Nested__NestedNested_EEE = 6
)

var Test_Nested__NestedNested_EEE_name = map[int32]string{
	0: "OK",
	1: "V1",
	2: "V2",
	3: "V3",
	4: "V4",
	5: "V5",
	6: "V6",
}
var Test_Nested__NestedNested_EEE_value = map[string]int32{
	"OK": 0,
	"V1": 1,
	"V2": 2,
	"V3": 3,
	"V4": 4,
	"V5": 5,
	"V6": 6,
}

func (x Test_Nested__NestedNested_EEE) Enum() *Test_Nested__NestedNested_EEE {
	p := new(Test_Nested__NestedNested_EEE)
	*p = x
	return p
}
func (x Test_Nested__NestedNested_EEE) String() string {
	return proto.EnumName(Test_Nested__NestedNested_EEE_name, int32(x))
}
func (x *Test_Nested__NestedNested_EEE) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Test_Nested__NestedNested_EEE_value, data, "Test_Nested__NestedNested_EEE")
	if err != nil {
		return err
	}
	*x = Test_Nested__NestedNested_EEE(value)
	return nil
}
func (Test_Nested__NestedNested_EEE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor3, []int{1, 1, 0, 0}
}

type Validator_Action int32

const (
	Validator_LOGIN Validator_Action = 0
	Validator_READ  Validator_Action = 1
	Validator_WRITE Validator_Action = 2
)

var Validator_Action_name = map[int32]string{
	0: "LOGIN",
	1: "READ",
	2: "WRITE",
}
var Validator_Action_value = map[string]int32{
	"LOGIN": 0,
	"READ":  1,
	"WRITE": 2,
}

func (x Validator_Action) Enum() *Validator_Action {
	p := new(Validator_Action)
	*p = x
	return p
}
func (x Validator_Action) String() string {
	return proto.EnumName(Validator_Action_name, int32(x))
}
func (x *Validator_Action) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Validator_Action_value, data, "Validator_Action")
	if err != nil {
		return err
	}
	*x = Validator_Action(value)
	return nil
}
func (Validator_Action) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{3, 0} }

type Simple struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Id               *uint64 `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Simple) Reset()                    { *m = Simple{} }
func (m *Simple) String() string            { return proto.CompactTextString(m) }
func (*Simple) ProtoMessage()               {}
func (*Simple) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *Simple) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Simple) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

type Test struct {
	Foo                          *string          `protobuf:"bytes,1,opt,name=foo,json=|foo|" json:"foo,omitempty"`
	Array                        []int32          `protobuf:"varint,2,rep,name=array" json:"array,omitempty"`
	S                            *Simple          `protobuf:"bytes,3,opt,name=s" json:"s,omitempty"`
	R                            []*Simple        `protobuf:"bytes,4,rep,name=r" json:"r,omitempty"`
	M                            map[string]int32 `protobuf:"bytes,5,rep,name=m" json:"m,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	B                            []byte           `protobuf:"bytes,6,opt,name=b,def=\\000\\001\\002\\003\\004\\005\\006\\007fubar!" json:"b,omitempty"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *Test) Reset()                    { *m = Test{} }
func (m *Test) String() string            { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()               {}
func (*Test) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

var extRange_Test = []proto.ExtensionRange{
	{100, 200},
	{300, 350},
	{500, 550},
}

func (*Test) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Test
}

var Default_Test_B []byte = []byte("\x00\x01\x02\x03\x04\x05\x06\afubar!")

func (m *Test) GetFoo() string {
	if m != nil && m.Foo != nil {
		return *m.Foo
	}
	return ""
}

func (m *Test) GetArray() []int32 {
	if m != nil {
		return m.Array
	}
	return nil
}

func (m *Test) GetS() *Simple {
	if m != nil {
		return m.S
	}
	return nil
}

func (m *Test) GetR() []*Simple {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *Test) GetM() map[string]int32 {
	if m != nil {
		return m.M
	}
	return nil
}

func (m *Test) GetB() []byte {
	if m != nil && m.B != nil {
		return m.B
	}
	return append([]byte(nil), Default_Test_B...)
}

type Test_Nested struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Test_Nested) Reset()                    { *m = Test_Nested{} }
func (m *Test_Nested) String() string            { return proto.CompactTextString(m) }
func (*Test_Nested) ProtoMessage()               {}
func (*Test_Nested) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1, 1} }

var E_Test_Nested_Fooblez = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.MessageOptions)(nil),
	ExtensionType: (*int32)(nil),
	Field:         20003,
	Name:          "foo.bar.Test.Nested.fooblez",
	Tag:           "varint,20003,opt,name=fooblez",
	Filename:      "desc_test_complex.proto",
}

type Test_Nested__NestedNested struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Test_Nested__NestedNested) Reset()                    { *m = Test_Nested__NestedNested{} }
func (m *Test_Nested__NestedNested) String() string            { return proto.CompactTextString(m) }
func (*Test_Nested__NestedNested) ProtoMessage()               {}
func (*Test_Nested__NestedNested) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1, 1, 0} }

var E_Test_Nested_XNestedNested_XGarblez = &proto.ExtensionDesc{
	ExtendedType:  (*Test)(nil),
	ExtensionType: (*string)(nil),
	Field:         100,
	Name:          "foo.bar.Test.Nested._NestedNested._garblez",
	Tag:           "bytes,100,opt,name=_garblez,json=Garblez",
	Filename:      "desc_test_complex.proto",
}

type Test_Nested__NestedNested_NestedNestedNested struct {
	Test             *Test  `protobuf:"bytes,1,opt,name=Test" json:"Test,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Test_Nested__NestedNested_NestedNestedNested) Reset() {
	*m = Test_Nested__NestedNested_NestedNestedNested{}
}
func (m *Test_Nested__NestedNested_NestedNestedNested) String() string {
	return proto.CompactTextString(m)
}
func (*Test_Nested__NestedNested_NestedNestedNested) ProtoMessage() {}
func (*Test_Nested__NestedNested_NestedNestedNested) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{1, 1, 0, 0}
}

func (m *Test_Nested__NestedNested_NestedNestedNested) GetTest() *Test {
	if m != nil {
		return m.Test
	}
	return nil
}

type Another struct {
	Test             *Test                          `protobuf:"bytes,1,opt,name=test" json:"test,omitempty"`
	Fff              *Test_Nested__NestedNested_EEE `protobuf:"varint,2,opt,name=fff,enum=foo.bar.Test_Nested__NestedNested_EEE,def=1" json:"fff,omitempty"`
	XXX_unrecognized []byte                         `json:"-"`
}

func (m *Another) Reset()                    { *m = Another{} }
func (m *Another) String() string            { return proto.CompactTextString(m) }
func (*Another) ProtoMessage()               {}
func (*Another) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

const Default_Another_Fff Test_Nested__NestedNested_EEE = Test_Nested__NestedNested_V1

func (m *Another) GetTest() *Test {
	if m != nil {
		return m.Test
	}
	return nil
}

func (m *Another) GetFff() Test_Nested__NestedNested_EEE {
	if m != nil && m.Fff != nil {
		return *m.Fff
	}
	return Default_Another_Fff
}

type Validator struct {
	Authenticated    *bool                   `protobuf:"varint,1,opt,name=authenticated" json:"authenticated,omitempty"`
	Permission       []*Validator_Permission `protobuf:"bytes,2,rep,name=permission" json:"permission,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *Validator) Reset()                    { *m = Validator{} }
func (m *Validator) String() string            { return proto.CompactTextString(m) }
func (*Validator) ProtoMessage()               {}
func (*Validator) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *Validator) GetAuthenticated() bool {
	if m != nil && m.Authenticated != nil {
		return *m.Authenticated
	}
	return false
}

func (m *Validator) GetPermission() []*Validator_Permission {
	if m != nil {
		return m.Permission
	}
	return nil
}

type Validator_Permission struct {
	Action           *Validator_Action `protobuf:"varint,1,opt,name=action,enum=foo.bar.Validator_Action" json:"action,omitempty"`
	Entity           *string           `protobuf:"bytes,2,opt,name=entity" json:"entity,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Validator_Permission) Reset()                    { *m = Validator_Permission{} }
func (m *Validator_Permission) String() string            { return proto.CompactTextString(m) }
func (*Validator_Permission) ProtoMessage()               {}
func (*Validator_Permission) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3, 0} }

func (m *Validator_Permission) GetAction() Validator_Action {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return Validator_LOGIN
}

func (m *Validator_Permission) GetEntity() string {
	if m != nil && m.Entity != nil {
		return *m.Entity
	}
	return ""
}

var E_Label = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.ExtensionRangeOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         20000,
	Name:          "foo.bar.label",
	Tag:           "bytes,20000,opt,name=label",
	Filename:      "desc_test_complex.proto",
}

var E_Rept = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.MessageOptions)(nil),
	ExtensionType: ([]*Test)(nil),
	Field:         20002,
	Name:          "foo.bar.rept",
	Tag:           "bytes,20002,rep,name=rept",
	Filename:      "desc_test_complex.proto",
}

var E_Eee = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.MessageOptions)(nil),
	ExtensionType: (*Test_Nested__NestedNested_EEE)(nil),
	Field:         20010,
	Name:          "foo.bar.eee",
	Tag:           "varint,20010,opt,name=eee,enum=foo.bar.Test_Nested__NestedNested_EEE",
	Filename:      "desc_test_complex.proto",
}

var E_A = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.MessageOptions)(nil),
	ExtensionType: (*Another)(nil),
	Field:         20020,
	Name:          "foo.bar.a",
	Tag:           "bytes,20020,opt,name=a",
	Filename:      "desc_test_complex.proto",
}

var E_Validator = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.MethodOptions)(nil),
	ExtensionType: (*Validator)(nil),
	Field:         12345,
	Name:          "foo.bar.validator",
	Tag:           "bytes,12345,opt,name=validator",
	Filename:      "desc_test_complex.proto",
}

func init() {
	proto.RegisterType((*Simple)(nil), "foo.bar.Simple")
	proto.RegisterType((*Test)(nil), "foo.bar.Test")
	proto.RegisterType((*Test_Nested)(nil), "foo.bar.Test.Nested")
	proto.RegisterType((*Test_Nested__NestedNested)(nil), "foo.bar.Test.Nested._NestedNested")
	proto.RegisterType((*Test_Nested__NestedNested_NestedNestedNested)(nil), "foo.bar.Test.Nested._NestedNested.NestedNestedNested")
	proto.RegisterType((*Another)(nil), "foo.bar.Another")
	proto.RegisterType((*Validator)(nil), "foo.bar.Validator")
	proto.RegisterType((*Validator_Permission)(nil), "foo.bar.Validator.Permission")
	proto.RegisterEnum("foo.bar.Test_Nested__NestedNested_EEE", Test_Nested__NestedNested_EEE_name, Test_Nested__NestedNested_EEE_value)
	proto.RegisterEnum("foo.bar.Validator_Action", Validator_Action_name, Validator_Action_value)
	proto.RegisterExtension(E_Test_Nested_Fooblez)
	proto.RegisterExtension(E_Test_Nested_XNestedNested_XGarblez)
	proto.RegisterExtension(E_Label)
	proto.RegisterExtension(E_Rept)
	proto.RegisterExtension(E_Eee)
	proto.RegisterExtension(E_A)
	proto.RegisterExtension(E_Validator)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for TestTestService service

type TestTestServiceClient interface {
	UserAuth(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error)
	Get(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error)
}

type testTestServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestTestServiceClient(cc *grpc.ClientConn) TestTestServiceClient {
	return &testTestServiceClient{cc}
}

func (c *testTestServiceClient) UserAuth(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error) {
	out := new(Test)
	err := grpc.Invoke(ctx, "/foo.bar.TestTestService/UserAuth", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testTestServiceClient) Get(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error) {
	out := new(Test)
	err := grpc.Invoke(ctx, "/foo.bar.TestTestService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TestTestService service

type TestTestServiceServer interface {
	UserAuth(context.Context, *Test) (*Test, error)
	Get(context.Context, *Test) (*Test, error)
}

func RegisterTestTestServiceServer(s *grpc.Server, srv TestTestServiceServer) {
	s.RegisterService(&_TestTestService_serviceDesc, srv)
}

func _TestTestService_UserAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestTestServiceServer).UserAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/foo.bar.TestTestService/UserAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestTestServiceServer).UserAuth(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestTestService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestTestServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/foo.bar.TestTestService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestTestServiceServer).Get(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

var _TestTestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "foo.bar.TestTestService",
	HandlerType: (*TestTestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserAuth",
			Handler:    _TestTestService_UserAuth_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _TestTestService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "desc_test_complex.proto",
}

func init() { proto.RegisterFile("desc_test_complex.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 1023 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x41, 0x6f, 0x1b, 0x45,
	0x14, 0xce, 0xec, 0xac, 0x37, 0xbb, 0x2f, 0x4d, 0x3a, 0x1a, 0x02, 0x18, 0x4b, 0xa5, 0xae, 0x41,
	0x91, 0x15, 0xa1, 0xcd, 0xc6, 0x49, 0x0a, 0x32, 0x2d, 0x52, 0x50, 0xad, 0x28, 0x82, 0x26, 0xb0,
	0x2d, 0xa9, 0x04, 0x96, 0xc2, 0xda, 0x3b, 0x6b, 0x6f, 0x59, 0xef, 0x58, 0xbb, 0xe3, 0xa8, 0x09,
	0x3d, 0x11, 0x89, 0x1f, 0x90, 0x53, 0x0f, 0x39, 0x80, 0x39, 0x20, 0x01, 0x47, 0x2e, 0x1c, 0xb9,
	0x95, 0x4a, 0x95, 0xb8, 0x72, 0x42, 0x72, 0x7e, 0x02, 0x70, 0x47, 0x33, 0xbb, 0x71, 0xea, 0x36,
	0x28, 0xe2, 0xe0, 0x7d, 0x3b, 0xf3, 0xde, 0xf7, 0xcd, 0xbc, 0xf7, 0xbe, 0xb7, 0x86, 0x57, 0x7d,
	0x96, 0xb6, 0x77, 0x05, 0x4b, 0xc5, 0x6e, 0x9b, 0xf7, 0xfa, 0x11, 0x7b, 0x60, 0xf7, 0x13, 0x2e,
	0x38, 0x9d, 0x0e, 0x38, 0xb7, 0x5b, 0x5e, 0x52, 0x2a, 0x77, 0x38, 0xef, 0x44, 0x6c, 0x49, 0x6d,
	0xb7, 0x06, 0xc1, 0x92, 0x44, 0x24, 0x61, 0x5f, 0xf0, 0x24, 0x0b, 0xad, 0xbc, 0x05, 0xc6, 0x9d,
	0x50, 0x62, 0x29, 0x05, 0x3d, 0xf6, 0x7a, 0xac, 0x88, 0xca, 0xa8, 0x6a, 0xb9, 0xea, 0x9d, 0xce,
	0x81, 0x16, 0xfa, 0x45, 0xad, 0x8c, 0xaa, 0xba, 0xab, 0x85, 0x7e, 0xe5, 0xa9, 0x0e, 0xfa, 0x5d,
	0x96, 0x0a, 0x4a, 0x01, 0x07, 0x9c, 0xe7, 0xb1, 0x85, 0x87, 0x01, 0xe7, 0x0f, 0xe9, 0x3c, 0x14,
	0xbc, 0x24, 0xf1, 0xf6, 0x8b, 0x5a, 0x19, 0x57, 0x0b, 0x6e, 0xb6, 0xa0, 0x57, 0x00, 0xa5, 0x45,
	0x5c, 0x46, 0xd5, 0x99, 0xda, 0x65, 0x3b, 0xbf, 0x97, 0x9d, 0x1d, 0xe9, 0xa2, 0x54, 0xba, 0x93,
	0xa2, 0x5e, 0xc6, 0xe7, 0xba, 0x13, 0x5a, 0x01, 0xd4, 0x2b, 0x16, 0x94, 0x7b, 0x7e, 0xec, 0x96,
	0x37, 0xb0, 0x6f, 0x37, 0x62, 0x91, 0xec, 0xbb, 0xa8, 0x47, 0x57, 0x01, 0xb5, 0x8a, 0x46, 0x19,
	0x55, 0x2f, 0xd5, 0x17, 0x9a, 0x8e, 0xe3, 0x34, 0x1d, 0x67, 0xb9, 0xe9, 0x38, 0xb5, 0xa6, 0xe3,
	0xac, 0x34, 0x1d, 0x67, 0xb5, 0xe9, 0x38, 0x6b, 0x4d, 0xc7, 0xb9, 0xde, 0x74, 0x9c, 0xb7, 0x83,
	0x41, 0xcb, 0x4b, 0xae, 0xb9, 0xa8, 0x55, 0x5a, 0x05, 0x23, 0xa3, 0xa0, 0x04, 0xf0, 0x17, 0x6c,
	0x3f, 0xcf, 0x45, 0xbe, 0xca, 0x4c, 0xf6, 0xbc, 0x68, 0xc0, 0x54, 0xe6, 0x05, 0x37, 0x5b, 0xd4,
	0xb5, 0x77, 0x50, 0xe9, 0x5b, 0x0d, 0x8c, 0x2d, 0x96, 0x0a, 0xe6, 0x97, 0xfe, 0x40, 0x30, 0xbb,
	0x9b, 0xbd, 0xe7, 0x3b, 0x9b, 0x40, 0x9f, 0x5d, 0x67, 0x4f, 0x7a, 0x2d, 0x2b, 0x99, 0xe2, 0x9f,
	0xa9, 0xcd, 0x4e, 0x64, 0xe1, 0x2a, 0x57, 0x9d, 0x1c, 0x8d, 0xac, 0x4b, 0x80, 0xbb, 0x9c, 0x0f,
	0x0d, 0x3d, 0xed, 0x73, 0x5e, 0xb9, 0x09, 0xb8, 0xd1, 0x68, 0x50, 0x03, 0xb4, 0xed, 0x0f, 0xc8,
	0x94, 0xb4, 0x3b, 0xcb, 0x04, 0x29, 0x5b, 0x23, 0x9a, 0xb2, 0x2b, 0x04, 0x2b, 0xbb, 0x4a, 0x74,
	0x65, 0xd7, 0x48, 0x41, 0xd9, 0xeb, 0xc4, 0xa8, 0x55, 0xc1, 0xdc, 0xed, 0x78, 0x49, 0x2b, 0x62,
	0x07, 0x74, 0xf2, 0xc4, 0xa2, 0xaf, 0xd2, 0x9c, 0xde, 0xc8, 0xbc, 0xf5, 0xf9, 0xa3, 0x91, 0x35,
	0x03, 0xb8, 0x23, 0x8f, 0xc6, 0x2d, 0xce, 0x1f, 0x8d, 0xac, 0x7f, 0xb6, 0x6a, 0xef, 0x82, 0x94,
	0x90, 0x82, 0x5f, 0xb5, 0x33, 0x0d, 0xd9, 0xa7, 0x1a, 0xb2, 0x6f, 0xb3, 0x34, 0xf5, 0x3a, 0x6c,
	0xbb, 0x2f, 0x42, 0x1e, 0xa7, 0xc5, 0xef, 0x8e, 0x91, 0xaa, 0xd2, 0x29, 0x62, 0xb1, 0x60, 0xfa,
	0xe4, 0x37, 0xb4, 0x48, 0xcc, 0x9f, 0x34, 0xf2, 0xa7, 0x56, 0x32, 0xbf, 0x1a, 0x59, 0xfa, 0x7d,
	0xef, 0xe0, 0x60, 0x91, 0x98, 0x7f, 0x63, 0xf2, 0xbd, 0x7e, 0xb6, 0x53, 0x79, 0xa4, 0xc1, 0xf4,
	0x7a, 0xcc, 0x45, 0x97, 0x25, 0xb2, 0x4e, 0xe2, 0xbf, 0xeb, 0x24, 0x5d, 0xf4, 0x06, 0xe0, 0x20,
	0x08, 0x54, 0x57, 0xe6, 0x6a, 0x0b, 0x93, 0x7a, 0xc8, 0xaa, 0x6d, 0x4f, 0x74, 0xc4, 0x6e, 0x34,
	0x1a, 0x75, 0x6d, 0x67, 0xd9, 0x95, 0xb0, 0xfa, 0x8f, 0xe8, 0x68, 0x64, 0xbd, 0x01, 0xd8, 0x6b,
	0xb5, 0x09, 0x22, 0x1a, 0xc1, 0xa5, 0x69, 0xa5, 0x62, 0xf2, 0x65, 0x05, 0x03, 0x0a, 0xe4, 0x23,
	0xad, 0x60, 0xf2, 0x18, 0x1f, 0x8d, 0xac, 0xab, 0x80, 0x7d, 0x16, 0x10, 0x4c, 0x34, 0x82, 0x4a,
	0x26, 0xe0, 0x96, 0x97, 0x90, 0x5f, 0x35, 0x19, 0xd2, 0x51, 0x71, 0x47, 0x23, 0xab, 0xa0, 0x42,
	0x7e, 0x1f, 0x59, 0x68, 0x78, 0x62, 0x69, 0x64, 0x6a, 0x78, 0x62, 0xbd, 0x0c, 0x2f, 0x2d, 0xe6,
	0x9c, 0xfe, 0x62, 0x8e, 0x7a, 0x2c, 0xdd, 0xd3, 0x50, 0x00, 0xdc, 0x5b, 0xe8, 0x0d, 0x4f, 0x2c,
	0x00, 0xb3, 0x64, 0x80, 0xbe, 0xcf, 0x23, 0x3e, 0x3c, 0xb1, 0x4c, 0x30, 0x4a, 0x3a, 0x79, 0x7a,
	0x68, 0x0c, 0x4f, 0x2c, 0x1d, 0x34, 0x82, 0x72, 0xab, 0x55, 0xfe, 0x42, 0x60, 0xed, 0x78, 0x51,
	0xe8, 0x7b, 0x82, 0x27, 0xf4, 0x4d, 0x98, 0xf5, 0x06, 0xa2, 0xcb, 0x62, 0x11, 0xb6, 0x3d, 0xc1,
	0x7c, 0x55, 0x25, 0xd3, 0x9d, 0xdc, 0xa4, 0x37, 0x01, 0xfa, 0x2c, 0xe9, 0x85, 0x69, 0x1a, 0xf2,
	0x58, 0x8d, 0xe1, 0x4c, 0xed, 0xca, 0xb8, 0x4c, 0x63, 0x36, 0xfb, 0xa3, 0x71, 0x90, 0xfb, 0x0c,
	0xa0, 0x74, 0x0f, 0xe0, 0xcc, 0x43, 0x97, 0xc1, 0xf0, 0xda, 0xb2, 0xc5, 0xea, 0xac, 0xb9, 0xda,
	0x6b, 0xe7, 0x10, 0xad, 0xab, 0x00, 0x37, 0x0f, 0xa4, 0xaf, 0x80, 0x21, 0x2f, 0x23, 0xf6, 0x55,
	0x8b, 0x2c, 0x37, 0x5f, 0x55, 0xaa, 0x60, 0x64, 0x91, 0xd4, 0x82, 0xc2, 0x87, 0xdb, 0x1b, 0x9b,
	0x5b, 0x64, 0x8a, 0x9a, 0xa0, 0xbb, 0x8d, 0xf5, 0x5b, 0x04, 0xc9, 0xcd, 0x7b, 0xee, 0xe6, 0xdd,
	0x06, 0xd1, 0x6a, 0x5f, 0x23, 0xb8, 0x2c, 0xdb, 0x29, 0x7f, 0x77, 0x58, 0xb2, 0x17, 0xb6, 0x19,
	0xbd, 0x01, 0xe6, 0x27, 0x29, 0x4b, 0xd6, 0x07, 0xa2, 0xfb, 0x9c, 0x98, 0x4b, 0x93, 0xcb, 0x0a,
	0x7d, 0x72, 0x68, 0xcc, 0x99, 0x88, 0x82, 0x39, 0x45, 0x8d, 0x76, 0x14, 0xb2, 0x58, 0xd0, 0x35,
	0xc0, 0x1b, 0x4c, 0x5c, 0x00, 0x24, 0x4f, 0x0e, 0x8d, 0x4b, 0x26, 0xa2, 0xa6, 0x89, 0xa8, 0x3e,
	0x48, 0x59, 0x52, 0x7f, 0x0f, 0x0a, 0x91, 0xd7, 0x62, 0x11, 0x5d, 0x78, 0x41, 0xff, 0x8d, 0x07,
	0x82, 0xc5, 0xaa, 0x78, 0x5e, 0x7c, 0x36, 0x06, 0xdf, 0x1c, 0xe7, 0x1f, 0x43, 0x05, 0xab, 0xdf,
	0x02, 0x3d, 0x61, 0x7d, 0x71, 0xf1, 0xf8, 0x0c, 0x8f, 0x91, 0xea, 0xd3, 0xf3, 0x82, 0x97, 0xe8,
	0xfa, 0x67, 0x80, 0x19, 0x63, 0x17, 0x93, 0xfc, 0x70, 0x8c, 0xfe, 0xcf, 0x4c, 0xb8, 0x92, 0xb5,
	0xbe, 0x0e, 0xc8, 0xbb, 0x98, 0xfa, 0xe7, 0xe3, 0x6c, 0x20, 0xc9, 0x98, 0x3a, 0x1f, 0x58, 0x17,
	0x79, 0xf5, 0x8f, 0xc1, 0xda, 0x1b, 0x6b, 0xf4, 0xf5, 0x73, 0xa8, 0x44, 0x97, 0xfb, 0xa7, 0x4c,
	0xbf, 0x7c, 0xae, 0x88, 0xe8, 0x8b, 0x3a, 0x72, 0xcf, 0x58, 0xde, 0x5f, 0xf9, 0x74, 0xb9, 0x13,
	0x8a, 0xee, 0xa0, 0x65, 0xb7, 0x79, 0x6f, 0xe9, 0x7e, 0x77, 0xd0, 0xeb, 0x67, 0x7f, 0x5f, 0x09,
	0x0b, 0x22, 0xd6, 0x16, 0x4b, 0x61, 0x2c, 0x58, 0x12, 0x7b, 0xd1, 0x92, 0xfc, 0x24, 0x28, 0x4f,
	0xfa, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x31, 0xe4, 0xc0, 0x0a, 0x07, 0x00, 0x00,
}
