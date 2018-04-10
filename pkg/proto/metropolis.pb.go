// Code generated by protoc-gen-go. DO NOT EDIT.
// source: metropolis.proto

/*
Package metropolis is a generated protocol buffer package.

It is generated from these files:
	metropolis.proto

It has these top-level messages:
	VerRequest
	VerResponse
*/
package metropolis

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

type VerFlag int32

const (
	VerFlag_AUTH      VerFlag = 0
	VerFlag_SIGNATURE VerFlag = 1
)

var VerFlag_name = map[int32]string{
	0: "AUTH",
	1: "SIGNATURE",
}
var VerFlag_value = map[string]int32{
	"AUTH":      0,
	"SIGNATURE": 1,
}

func (x VerFlag) String() string {
	return proto.EnumName(VerFlag_name, int32(x))
}
func (VerFlag) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type VerStatus int32

const (
	VerStatus_SUCCESS VerStatus = 0
	VerStatus_ERROR   VerStatus = 1
)

var VerStatus_name = map[int32]string{
	0: "SUCCESS",
	1: "ERROR",
}
var VerStatus_value = map[string]int32{
	"SUCCESS": 0,
	"ERROR":   1,
}

func (x VerStatus) String() string {
	return proto.EnumName(VerStatus_name, int32(x))
}
func (VerStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type VerRequest struct {
	SignedXml string  `protobuf:"bytes,1,opt,name=signedXml" json:"signedXml,omitempty"`
	Flag      VerFlag `protobuf:"varint,2,opt,name=flag,enum=metropolis.VerFlag" json:"flag,omitempty"`
}

func (m *VerRequest) Reset()                    { *m = VerRequest{} }
func (m *VerRequest) String() string            { return proto.CompactTextString(m) }
func (*VerRequest) ProtoMessage()               {}
func (*VerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *VerRequest) GetSignedXml() string {
	if m != nil {
		return m.SignedXml
	}
	return ""
}

func (m *VerRequest) GetFlag() VerFlag {
	if m != nil {
		return m.Flag
	}
	return VerFlag_AUTH
}

type VerResponse struct {
	Status      VerStatus `protobuf:"varint,1,opt,name=status,enum=metropolis.VerStatus" json:"status,omitempty"`
	Message     string    `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Description string    `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *VerResponse) Reset()                    { *m = VerResponse{} }
func (m *VerResponse) String() string            { return proto.CompactTextString(m) }
func (*VerResponse) ProtoMessage()               {}
func (*VerResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *VerResponse) GetStatus() VerStatus {
	if m != nil {
		return m.Status
	}
	return VerStatus_SUCCESS
}

func (m *VerResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *VerResponse) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*VerRequest)(nil), "metropolis.VerRequest")
	proto.RegisterType((*VerResponse)(nil), "metropolis.VerResponse")
	proto.RegisterEnum("metropolis.VerFlag", VerFlag_name, VerFlag_value)
	proto.RegisterEnum("metropolis.VerStatus", VerStatus_name, VerStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MetropolisService service

type MetropolisServiceClient interface {
	VerifySignature(ctx context.Context, in *VerRequest, opts ...grpc.CallOption) (*VerResponse, error)
}

type metropolisServiceClient struct {
	cc *grpc.ClientConn
}

func NewMetropolisServiceClient(cc *grpc.ClientConn) MetropolisServiceClient {
	return &metropolisServiceClient{cc}
}

func (c *metropolisServiceClient) VerifySignature(ctx context.Context, in *VerRequest, opts ...grpc.CallOption) (*VerResponse, error) {
	out := new(VerResponse)
	err := grpc.Invoke(ctx, "/metropolis.MetropolisService/VerifySignature", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MetropolisService service

type MetropolisServiceServer interface {
	VerifySignature(context.Context, *VerRequest) (*VerResponse, error)
}

func RegisterMetropolisServiceServer(s *grpc.Server, srv MetropolisServiceServer) {
	s.RegisterService(&_MetropolisService_serviceDesc, srv)
}

func _MetropolisService_VerifySignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetropolisServiceServer).VerifySignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/metropolis.MetropolisService/VerifySignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetropolisServiceServer).VerifySignature(ctx, req.(*VerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MetropolisService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "metropolis.MetropolisService",
	HandlerType: (*MetropolisServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifySignature",
			Handler:    _MetropolisService_VerifySignature_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metropolis.proto",
}

func init() { proto.RegisterFile("metropolis.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x4b, 0x4f, 0x3a, 0x31,
	0x14, 0xc5, 0x99, 0xff, 0x1f, 0xc1, 0xb9, 0x44, 0x18, 0xaf, 0x41, 0x27, 0xc6, 0x05, 0x19, 0x17,
	0x12, 0x12, 0x59, 0xe0, 0x27, 0x00, 0xc4, 0xc7, 0xc2, 0x47, 0x5a, 0x20, 0xba, 0x2c, 0x70, 0x99,
	0x34, 0x0e, 0xd3, 0xb1, 0x2d, 0x46, 0xfd, 0xf4, 0xc6, 0x8a, 0x30, 0xc1, 0x65, 0x7f, 0xf7, 0xf4,
	0x9c, 0xfb, 0x80, 0x60, 0x41, 0x56, 0xab, 0x4c, 0x25, 0xd2, 0xb4, 0x33, 0xad, 0xac, 0x42, 0xd8,
	0x90, 0x88, 0x03, 0x8c, 0x49, 0x33, 0x7a, 0x5d, 0x92, 0xb1, 0x78, 0x02, 0xbe, 0x91, 0x71, 0x4a,
	0xb3, 0xa7, 0x45, 0x12, 0x7a, 0x0d, 0xaf, 0xe9, 0xb3, 0x0d, 0xc0, 0x33, 0x28, 0xce, 0x13, 0x11,
	0x87, 0xff, 0x1a, 0x5e, 0xb3, 0xda, 0x39, 0x68, 0xe7, 0x8c, 0xc7, 0xa4, 0xaf, 0x12, 0x11, 0x33,
	0x27, 0x88, 0xde, 0xa1, 0xe2, 0x4c, 0x4d, 0xa6, 0x52, 0x43, 0x78, 0x0e, 0x25, 0x63, 0x85, 0x5d,
	0x1a, 0x67, 0x59, 0xed, 0xd4, 0xb7, 0x7e, 0x72, 0x57, 0x64, 0x2b, 0x11, 0x86, 0x50, 0x5e, 0x90,
	0x31, 0x22, 0x26, 0x97, 0xe4, 0xb3, 0xdf, 0x27, 0x36, 0xa0, 0x32, 0x23, 0x33, 0xd5, 0x32, 0xb3,
	0x52, 0xa5, 0xe1, 0x7f, 0x57, 0xcd, 0xa3, 0x56, 0x04, 0xe5, 0x55, 0x2b, 0xb8, 0x0b, 0xc5, 0xee,
	0x68, 0x78, 0x13, 0x14, 0x70, 0x0f, 0x7c, 0x7e, 0x7b, 0x7d, 0xdf, 0x1d, 0x8e, 0xd8, 0x20, 0xf0,
	0x5a, 0xa7, 0xe0, 0xaf, 0x43, 0xb1, 0x02, 0x65, 0x3e, 0xea, 0xf7, 0x07, 0x9c, 0x07, 0x05, 0xf4,
	0x61, 0x67, 0xc0, 0xd8, 0x03, 0x0b, 0xbc, 0xce, 0x33, 0xec, 0xdf, 0xad, 0x9b, 0xe4, 0xa4, 0xdf,
	0xe4, 0x94, 0xf0, 0x12, 0x6a, 0x63, 0xd2, 0x72, 0xfe, 0xc1, 0x65, 0x9c, 0x0a, 0xbb, 0xd4, 0x84,
	0x87, 0x5b, 0xb3, 0xac, 0x36, 0x79, 0x7c, 0xf4, 0x87, 0xff, 0x2c, 0x23, 0x2a, 0xf4, 0x9a, 0x50,
	0x7f, 0xf9, 0x6c, 0x4f, 0x85, 0x9e, 0x08, 0x43, 0x39, 0x59, 0xaf, 0xb6, 0x49, 0x7c, 0xfc, 0x3e,
	0xd4, 0xa4, 0xe4, 0xee, 0x75, 0xf1, 0x15, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x93, 0x86, 0x8e, 0xc3,
	0x01, 0x00, 0x00,
}