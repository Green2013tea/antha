// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/antha-lang/antha/driver/antha_driver_v1/driver.proto

/*
Package antha_driver_v1 is a generated protocol buffer package.

It is generated from these files:
	github.com/antha-lang/antha/driver/antha_driver_v1/driver.proto

It has these top-level messages:
	TypeRequest
	TypeReply
	HttpHeader
	HttpCall
*/
package antha_driver_v1

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

type TypeRequest struct {
}

func (m *TypeRequest) Reset()                    { *m = TypeRequest{} }
func (m *TypeRequest) String() string            { return proto.CompactTextString(m) }
func (*TypeRequest) ProtoMessage()               {}
func (*TypeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type TypeReply struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *TypeReply) Reset()                    { *m = TypeReply{} }
func (m *TypeReply) String() string            { return proto.CompactTextString(m) }
func (*TypeReply) ProtoMessage()               {}
func (*TypeReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TypeReply) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type HttpHeader struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *HttpHeader) Reset()                    { *m = HttpHeader{} }
func (m *HttpHeader) String() string            { return proto.CompactTextString(m) }
func (*HttpHeader) ProtoMessage()               {}
func (*HttpHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *HttpHeader) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *HttpHeader) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Remote Http call
type HttpCall struct {
	Url     string        `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Method  string        `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	Body    []byte        `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Headers []*HttpHeader `protobuf:"bytes,4,rep,name=headers" json:"headers,omitempty"`
}

func (m *HttpCall) Reset()                    { *m = HttpCall{} }
func (m *HttpCall) String() string            { return proto.CompactTextString(m) }
func (*HttpCall) ProtoMessage()               {}
func (*HttpCall) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *HttpCall) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *HttpCall) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *HttpCall) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *HttpCall) GetHeaders() []*HttpHeader {
	if m != nil {
		return m.Headers
	}
	return nil
}

func init() {
	proto.RegisterType((*TypeRequest)(nil), "antha.driver.v1.TypeRequest")
	proto.RegisterType((*TypeReply)(nil), "antha.driver.v1.TypeReply")
	proto.RegisterType((*HttpHeader)(nil), "antha.driver.v1.HttpHeader")
	proto.RegisterType((*HttpCall)(nil), "antha.driver.v1.HttpCall")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Driver service

type DriverClient interface {
	DriverType(ctx context.Context, in *TypeRequest, opts ...grpc.CallOption) (*TypeReply, error)
}

type driverClient struct {
	cc *grpc.ClientConn
}

func NewDriverClient(cc *grpc.ClientConn) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) DriverType(ctx context.Context, in *TypeRequest, opts ...grpc.CallOption) (*TypeReply, error) {
	out := new(TypeReply)
	err := grpc.Invoke(ctx, "/antha.driver.v1.Driver/DriverType", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Driver service

type DriverServer interface {
	DriverType(context.Context, *TypeRequest) (*TypeReply, error)
}

func RegisterDriverServer(s *grpc.Server, srv DriverServer) {
	s.RegisterService(&_Driver_serviceDesc, srv)
}

func _Driver_DriverType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).DriverType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/antha.driver.v1.Driver/DriverType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).DriverType(ctx, req.(*TypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Driver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "antha.driver.v1.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DriverType",
			Handler:    _Driver_DriverType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/antha-lang/antha/driver/antha_driver_v1/driver.proto",
}

func init() {
	proto.RegisterFile("github.com/antha-lang/antha/driver/antha_driver_v1/driver.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xcd, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x89, 0xa9, 0xd1, 0x4e, 0x15, 0x65, 0x11, 0x09, 0x55, 0x30, 0xe4, 0x94, 0x8b, 0x5b,
	0x5a, 0xd1, 0xab, 0x07, 0x45, 0x7a, 0x94, 0xe0, 0xbd, 0x6c, 0xcc, 0xd0, 0x14, 0x36, 0xd9, 0x75,
	0xb3, 0x09, 0xec, 0xc9, 0x7f, 0x5d, 0xf6, 0xc3, 0x0f, 0x14, 0x6f, 0xbf, 0xf7, 0xe6, 0xbd, 0x64,
	0x67, 0xe0, 0x7e, 0xbb, 0xd3, 0xcd, 0x50, 0xd1, 0x57, 0xd1, 0x2e, 0x58, 0xa7, 0x1b, 0x76, 0xcd,
	0x59, 0xb7, 0xf5, 0xb8, 0xa8, 0xd5, 0x6e, 0x44, 0xe5, 0xc5, 0xc6, 0x8b, 0xcd, 0xb8, 0x0c, 0x36,
	0x95, 0x4a, 0x68, 0x41, 0x4e, 0xdc, 0x94, 0x06, 0x6f, 0x5c, 0xe6, 0xc7, 0x30, 0x7b, 0x31, 0x12,
	0x4b, 0x7c, 0x1b, 0xb0, 0xd7, 0xf9, 0x15, 0x4c, 0xbd, 0x94, 0xdc, 0x10, 0x02, 0x13, 0x6d, 0x24,
	0xa6, 0x51, 0x16, 0x15, 0xd3, 0xd2, 0x71, 0x7e, 0x07, 0xb0, 0xd6, 0x5a, 0xae, 0x91, 0xd5, 0xa8,
	0x6c, 0xa2, 0x63, 0xed, 0x57, 0xc2, 0x32, 0x39, 0x83, 0xfd, 0x91, 0xf1, 0x01, 0xd3, 0x3d, 0x67,
	0x7a, 0x91, 0xbf, 0xc3, 0xa1, 0xed, 0x3d, 0x30, 0xce, 0xc9, 0x29, 0xc4, 0x83, 0xe2, 0xa1, 0x64,
	0x91, 0x9c, 0x43, 0xd2, 0xa2, 0x6e, 0x44, 0x1d, 0x4a, 0x41, 0xd9, 0xef, 0x57, 0xa2, 0x36, 0x69,
	0x9c, 0x45, 0xc5, 0x51, 0xe9, 0x98, 0xdc, 0xc2, 0x41, 0xe3, 0xfe, 0xde, 0xa7, 0x93, 0x2c, 0x2e,
	0x66, 0xab, 0x0b, 0xfa, 0x6b, 0x29, 0xfa, 0xfd, 0xc2, 0xf2, 0x33, 0xbb, 0x7a, 0x86, 0xe4, 0xd1,
	0x05, 0xc8, 0x13, 0x80, 0x27, 0xbb, 0x29, 0xb9, 0xfc, 0xd3, 0xfe, 0x71, 0x8f, 0xf9, 0xfc, 0x9f,
	0xa9, 0xe4, 0xa6, 0x4a, 0xdc, 0x49, 0x6f, 0x3e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xbf, 0x5e,
	0xee, 0x95, 0x01, 0x00, 0x00,
}
