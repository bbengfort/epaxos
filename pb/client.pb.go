// Code generated by protoc-gen-go. DO NOT EDIT.
// source: client.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type ProposeRequest struct {
	Identity             string     `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	Op                   *Operation `protobuf:"bytes,2,opt,name=op,proto3" json:"op,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ProposeRequest) Reset()         { *m = ProposeRequest{} }
func (m *ProposeRequest) String() string { return proto.CompactTextString(m) }
func (*ProposeRequest) ProtoMessage()    {}
func (*ProposeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_014de31d7ac8c57c, []int{0}
}

func (m *ProposeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProposeRequest.Unmarshal(m, b)
}
func (m *ProposeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProposeRequest.Marshal(b, m, deterministic)
}
func (m *ProposeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposeRequest.Merge(m, src)
}
func (m *ProposeRequest) XXX_Size() int {
	return xxx_messageInfo_ProposeRequest.Size(m)
}
func (m *ProposeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProposeRequest proto.InternalMessageInfo

func (m *ProposeRequest) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

func (m *ProposeRequest) GetOp() *Operation {
	if m != nil {
		return m.Op
	}
	return nil
}

type ProposeReply struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Slot                 int64    `protobuf:"varint,4,opt,name=slot,proto3" json:"slot,omitempty"`
	Key                  string   `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProposeReply) Reset()         { *m = ProposeReply{} }
func (m *ProposeReply) String() string { return proto.CompactTextString(m) }
func (*ProposeReply) ProtoMessage()    {}
func (*ProposeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_014de31d7ac8c57c, []int{1}
}

func (m *ProposeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProposeReply.Unmarshal(m, b)
}
func (m *ProposeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProposeReply.Marshal(b, m, deterministic)
}
func (m *ProposeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposeReply.Merge(m, src)
}
func (m *ProposeReply) XXX_Size() int {
	return xxx_messageInfo_ProposeReply.Size(m)
}
func (m *ProposeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposeReply.DiscardUnknown(m)
}

var xxx_messageInfo_ProposeReply proto.InternalMessageInfo

func (m *ProposeReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ProposeReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ProposeReply) GetSlot() int64 {
	if m != nil {
		return m.Slot
	}
	return 0
}

func (m *ProposeReply) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ProposeReply) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*ProposeRequest)(nil), "pb.ProposeRequest")
	proto.RegisterType((*ProposeReply)(nil), "pb.ProposeReply")
}

func init() { proto.RegisterFile("client.proto", fileDescriptor_014de31d7ac8c57c) }

var fileDescriptor_014de31d7ac8c57c = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x8f, 0xcd, 0x4a, 0x04, 0x31,
	0x0c, 0x80, 0x69, 0xf7, 0xc7, 0xdd, 0x58, 0x45, 0x82, 0x87, 0xb2, 0x20, 0x94, 0x3d, 0xf5, 0x34,
	0x07, 0x7d, 0x0c, 0x0f, 0x4a, 0xdf, 0x60, 0x66, 0xcc, 0xa1, 0x58, 0x26, 0xb1, 0xed, 0x88, 0xe3,
	0xd3, 0xcb, 0x76, 0x70, 0x6f, 0xdf, 0x17, 0x92, 0x0f, 0x02, 0x66, 0x4c, 0x91, 0xa6, 0xda, 0x49,
	0xe6, 0xca, 0xa8, 0x65, 0x38, 0x19, 0x92, 0xfe, 0x87, 0xcb, 0x3a, 0x39, 0xbf, 0xc2, 0xfd, 0x7b,
	0x66, 0xe1, 0x42, 0x81, 0xbe, 0x66, 0x2a, 0x15, 0x4f, 0x70, 0x88, 0x1f, 0x34, 0xd5, 0x58, 0x17,
	0xab, 0x9c, 0xf2, 0xc7, 0x70, 0x75, 0x7c, 0x02, 0xcd, 0x62, 0xb5, 0x53, 0xfe, 0xf6, 0xf9, 0xae,
	0x93, 0xa1, 0x7b, 0x13, 0xca, 0x7d, 0x8d, 0x3c, 0x05, 0xcd, 0x72, 0xfe, 0x05, 0x73, 0x8d, 0x49,
	0x5a, 0xd0, 0xc2, 0x4d, 0x99, 0xc7, 0x91, 0x4a, 0x69, 0xa5, 0x43, 0xf8, 0x57, 0x7c, 0x84, 0x1d,
	0xe5, 0xcc, 0xb9, 0xb5, 0x8e, 0x61, 0x15, 0x44, 0xd8, 0x96, 0xc4, 0xd5, 0x6e, 0x9d, 0xf2, 0x9b,
	0xd0, 0x18, 0x1f, 0x60, 0xf3, 0x49, 0x8b, 0xdd, 0xb5, 0xbd, 0x0b, 0x5e, 0x6e, 0xbf, 0xfb, 0x34,
	0x93, 0xdd, 0x3b, 0xe5, 0x4d, 0x58, 0x65, 0xd8, 0xb7, 0x7f, 0x5e, 0xfe, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x28, 0x0f, 0x54, 0x39, 0xf1, 0x00, 0x00, 0x00,
}