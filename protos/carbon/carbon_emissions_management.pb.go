// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/carbon/carbon_emissions_management.proto

package nrm

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Carbon Emissions Update
type ResourceUpdate struct {
	WeightGrams          int64                `protobuf:"varint,1,opt,name=weight_grams,json=weightGrams,proto3" json:"weight_grams,omitempty"`
	Start                *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	Finish               *timestamp.Timestamp `protobuf:"bytes,3,opt,name=finish,proto3" json:"finish,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ResourceUpdate) Reset()         { *m = ResourceUpdate{} }
func (m *ResourceUpdate) String() string { return proto.CompactTextString(m) }
func (*ResourceUpdate) ProtoMessage()    {}
func (*ResourceUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_fd4a03625c8f0292, []int{0}
}

func (m *ResourceUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourceUpdate.Unmarshal(m, b)
}
func (m *ResourceUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourceUpdate.Marshal(b, m, deterministic)
}
func (m *ResourceUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceUpdate.Merge(m, src)
}
func (m *ResourceUpdate) XXX_Size() int {
	return xxx_messageInfo_ResourceUpdate.Size(m)
}
func (m *ResourceUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceUpdate proto.InternalMessageInfo

func (m *ResourceUpdate) GetWeightGrams() int64 {
	if m != nil {
		return m.WeightGrams
	}
	return 0
}

func (m *ResourceUpdate) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *ResourceUpdate) GetFinish() *timestamp.Timestamp {
	if m != nil {
		return m.Finish
	}
	return nil
}

// Summary of the storage operation
type StorageSummary struct {
	Success              bool                 `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Attempts             int32                `protobuf:"varint,2,opt,name=attempts,proto3" json:"attempts,omitempty"`
	LastAttemptedAt      *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_attempted_at,json=lastAttemptedAt,proto3" json:"last_attempted_at,omitempty"`
	FailureMessage       string               `protobuf:"bytes,4,opt,name=failure_message,json=failureMessage,proto3" json:"failure_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *StorageSummary) Reset()         { *m = StorageSummary{} }
func (m *StorageSummary) String() string { return proto.CompactTextString(m) }
func (*StorageSummary) ProtoMessage()    {}
func (*StorageSummary) Descriptor() ([]byte, []int) {
	return fileDescriptor_fd4a03625c8f0292, []int{1}
}

func (m *StorageSummary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageSummary.Unmarshal(m, b)
}
func (m *StorageSummary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageSummary.Marshal(b, m, deterministic)
}
func (m *StorageSummary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageSummary.Merge(m, src)
}
func (m *StorageSummary) XXX_Size() int {
	return xxx_messageInfo_StorageSummary.Size(m)
}
func (m *StorageSummary) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageSummary.DiscardUnknown(m)
}

var xxx_messageInfo_StorageSummary proto.InternalMessageInfo

func (m *StorageSummary) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *StorageSummary) GetAttempts() int32 {
	if m != nil {
		return m.Attempts
	}
	return 0
}

func (m *StorageSummary) GetLastAttemptedAt() *timestamp.Timestamp {
	if m != nil {
		return m.LastAttemptedAt
	}
	return nil
}

func (m *StorageSummary) GetFailureMessage() string {
	if m != nil {
		return m.FailureMessage
	}
	return ""
}

func init() {
	proto.RegisterType((*ResourceUpdate)(nil), "nrm.ResourceUpdate")
	proto.RegisterType((*StorageSummary)(nil), "nrm.StorageSummary")
}

func init() {
	proto.RegisterFile("protos/carbon/carbon_emissions_management.proto", fileDescriptor_fd4a03625c8f0292)
}

var fileDescriptor_fd4a03625c8f0292 = []byte{
	// 359 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0xaf, 0xd3, 0x30,
	0x10, 0xc7, 0x31, 0xa5, 0x8f, 0x87, 0x1f, 0xea, 0x13, 0x66, 0x09, 0x59, 0x28, 0x5d, 0xc8, 0xd2,
	0xa4, 0x0a, 0x03, 0x73, 0x19, 0x60, 0x2a, 0x43, 0x5a, 0x16, 0x96, 0xe8, 0x92, 0x5e, 0x1d, 0xab,
	0xb1, 0x1d, 0xd9, 0x17, 0x21, 0x3e, 0x0a, 0x9f, 0x84, 0xaf, 0x87, 0x52, 0x27, 0x95, 0x3a, 0xf1,
	0xa6, 0xe4, 0x7e, 0xfe, 0xff, 0x75, 0xff, 0xbb, 0xe3, 0x59, 0xe7, 0x2c, 0x59, 0x9f, 0xd5, 0xe0,
	0x2a, 0x6b, 0xc6, 0x4f, 0x89, 0x5a, 0x79, 0xaf, 0xac, 0xf1, 0xa5, 0x06, 0x03, 0x12, 0x35, 0x1a,
	0x4a, 0x2f, 0x4a, 0x31, 0x33, 0x4e, 0xc7, 0xef, 0xa5, 0xb5, 0xb2, 0xc5, 0x60, 0xae, 0xfa, 0x53,
	0x46, 0x4a, 0xa3, 0x27, 0xd0, 0x5d, 0x50, 0xad, 0xfe, 0x30, 0xbe, 0x28, 0xd0, 0xdb, 0xde, 0xd5,
	0xf8, 0xa3, 0x3b, 0x02, 0xa1, 0xf8, 0xc0, 0x5f, 0xff, 0x42, 0x25, 0x1b, 0x2a, 0xa5, 0x03, 0xed,
	0x23, 0xb6, 0x64, 0xc9, 0xac, 0x78, 0x08, 0xec, 0xdb, 0x80, 0xc4, 0x86, 0xcf, 0x3d, 0x81, 0xa3,
	0xe8, 0xf9, 0x92, 0x25, 0x0f, 0x79, 0x9c, 0x86, 0x36, 0xe9, 0xd4, 0x26, 0x3d, 0x4c, 0x6d, 0x8a,
	0x20, 0x14, 0x39, 0xbf, 0x3b, 0x29, 0xa3, 0x7c, 0x13, 0xcd, 0xfe, 0x6b, 0x19, 0x95, 0xab, 0xbf,
	0x8c, 0x2f, 0xf6, 0x64, 0x1d, 0x48, 0xdc, 0xf7, 0x5a, 0x83, 0xfb, 0x2d, 0x22, 0xfe, 0xd2, 0xf7,
	0x75, 0x8d, 0x3e, 0xc4, 0xba, 0x2f, 0xa6, 0x52, 0xc4, 0xfc, 0x1e, 0x88, 0x50, 0x77, 0xe4, 0x2f,
	0xa9, 0xe6, 0xc5, 0xb5, 0x16, 0x5f, 0xf9, 0x9b, 0x16, 0x3c, 0x95, 0x23, 0xc0, 0x63, 0x09, 0xf4,
	0x84, 0x1c, 0x8f, 0x83, 0x69, 0x3b, 0x79, 0xb6, 0x24, 0x3e, 0xf2, 0xc7, 0x13, 0xa8, 0xb6, 0x77,
	0x58, 0x6a, 0xf4, 0x1e, 0x24, 0x46, 0x2f, 0x96, 0x2c, 0x79, 0x55, 0x2c, 0x46, 0xbc, 0x0b, 0x34,
	0x3f, 0xf0, 0x77, 0xdf, 0x81, 0x7a, 0x07, 0xed, 0xb4, 0xdb, 0xdd, 0xf5, 0x3c, 0xe2, 0x33, 0x9f,
	0x0f, 0x53, 0xa1, 0x78, 0x9b, 0x1a, 0xa7, 0xd3, 0xdb, 0xed, 0xc7, 0x01, 0xde, 0x8e, 0xbd, 0x7a,
	0x96, 0xb0, 0x0d, 0xfb, 0x92, 0xff, 0xdc, 0x48, 0x45, 0x4d, 0x5f, 0xa5, 0xb5, 0xd5, 0x99, 0xa9,
	0x3a, 0x70, 0x67, 0x74, 0x99, 0x71, 0x7a, 0x8d, 0xd4, 0xac, 0xeb, 0x56, 0xa1, 0xa1, 0xac, 0x3b,
	0xcb, 0x70, 0xee, 0xe1, 0xa1, 0xba, 0xbb, 0xfc, 0x7e, 0xfa, 0x17, 0x00, 0x00, 0xff, 0xff, 0xb2,
	0x29, 0x3e, 0x43, 0x3f, 0x02, 0x00, 0x00,
}