// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/antha-lang/antha/driver/antha_human_v1/human.proto

/*
Package antha_human_v1 is a generated protocol buffer package.

It is generated from these files:
	github.com/antha-lang/antha/driver/antha_human_v1/human.proto

It has these top-level messages:
*/
package antha_human_v1

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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Human service

type HumanClient interface {
}

type humanClient struct {
	cc *grpc.ClientConn
}

func NewHumanClient(cc *grpc.ClientConn) HumanClient {
	return &humanClient{cc}
}

// Server API for Human service

type HumanServer interface {
}

func RegisterHumanServer(s *grpc.Server, srv HumanServer) {
	s.RegisterService(&_Human_serviceDesc, srv)
}

var _Human_serviceDesc = grpc.ServiceDesc{
	ServiceName: "antha.human.v1.Human",
	HandlerType: (*HumanServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "github.com/antha-lang/antha/driver/antha_human_v1/human.proto",
}

func init() {
	proto.RegisterFile("github.com/antha-lang/antha/driver/antha_human_v1/human.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 96 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x4d, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcc, 0x2b, 0xc9, 0x48, 0xd4, 0xcd, 0x49, 0xcc,
	0x4b, 0x87, 0x30, 0xf5, 0x53, 0x8a, 0x32, 0xcb, 0x52, 0x8b, 0x20, 0x9c, 0xf8, 0x8c, 0xd2, 0xdc,
	0xc4, 0xbc, 0xf8, 0x32, 0x43, 0x7d, 0x30, 0x43, 0xaf, 0xa0, 0x28, 0xbf, 0x24, 0x5f, 0x88, 0x0f,
	0x2c, 0xa7, 0x07, 0x11, 0x2a, 0x33, 0x34, 0x62, 0xe7, 0x62, 0xf5, 0x00, 0xb1, 0x93, 0xd8, 0xc0,
	0xf2, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0x4f, 0x5a, 0x11, 0x60, 0x00, 0x00, 0x00,
}
