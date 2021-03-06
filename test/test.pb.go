// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test/test.proto

/*
Package test is a generated protocol buffer package.

It is generated from these files:
	test/test.proto

It has these top-level messages:
	PingRequest
	Pong
*/
package test

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

type PingRequest struct {
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Pong struct {
}

func (m *Pong) Reset()                    { *m = Pong{} }
func (m *Pong) String() string            { return proto.CompactTextString(m) }
func (*Pong) ProtoMessage()               {}
func (*Pong) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*PingRequest)(nil), "com.github.k2wanko.gaegrpc.test.PingRequest")
	proto.RegisterType((*Pong)(nil), "com.github.k2wanko.gaegrpc.test.Pong")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Test service

type TestClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*Pong, error)
}

type testClient struct {
	cc *grpc.ClientConn
}

func NewTestClient(cc *grpc.ClientConn) TestClient {
	return &testClient{cc}
}

func (c *testClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*Pong, error) {
	out := new(Pong)
	err := grpc.Invoke(ctx, "/com.github.k2wanko.gaegrpc.test.Test/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestServer interface {
	Ping(context.Context, *PingRequest) (*Pong, error)
}

func RegisterTestServer(s *grpc.Server, srv TestServer) {
	s.RegisterService(&_Test_serviceDesc, srv)
}

func _Test_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.github.k2wanko.gaegrpc.test.Test/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Test_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.github.k2wanko.gaegrpc.test.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Test_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test/test.proto",
}

func init() { proto.RegisterFile("test/test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x49, 0x2d, 0x2e,
	0xd1, 0x07, 0x11, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xf2, 0xc9, 0xf9, 0xb9, 0x7a, 0xe9,
	0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xd9, 0x46, 0xe5, 0x89, 0x79, 0xd9, 0xf9, 0x7a, 0xe9, 0x89,
	0xa9, 0xe9, 0x45, 0x05, 0xc9, 0x7a, 0x20, 0x65, 0x4a, 0xbc, 0x5c, 0xdc, 0x01, 0x99, 0x79, 0xe9,
	0x41, 0xa9, 0x85, 0xa5, 0x20, 0x2e, 0x1b, 0x17, 0x4b, 0x40, 0x7e, 0x5e, 0xba, 0x51, 0x32, 0x17,
	0x4b, 0x48, 0x6a, 0x71, 0x89, 0x50, 0x34, 0x17, 0x0b, 0x48, 0x5a, 0x48, 0x47, 0x8f, 0x80, 0x41,
	0x7a, 0x48, 0xa6, 0x48, 0xa9, 0x12, 0x56, 0x9d, 0x9f, 0x97, 0xee, 0xc4, 0x16, 0xc5, 0x02, 0xe2,
	0x24, 0xb1, 0x81, 0xdd, 0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x4d, 0x1d, 0x35, 0xbe,
	0x00, 0x00, 0x00,
}
