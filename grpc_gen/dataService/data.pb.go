// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.9.0
// source: proto/data.proto

package dataService

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetStateReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStateReq) Reset() {
	*x = GetStateReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateReq) ProtoMessage() {}

func (x *GetStateReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateReq.ProtoReflect.Descriptor instead.
func (*GetStateReq) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{0}
}

type GetStateResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsLeader bool `protobuf:"varint,1,opt,name=IsLeader,proto3" json:"IsLeader,omitempty"`
}

func (x *GetStateResp) Reset() {
	*x = GetStateResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateResp) ProtoMessage() {}

func (x *GetStateResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateResp.ProtoReflect.Descriptor instead.
func (*GetStateResp) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{1}
}

func (x *GetStateResp) GetIsLeader() bool {
	if x != nil {
		return x.IsLeader
	}
	return false
}

type GetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *GetReq) Reset() {
	*x = GetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReq) ProtoMessage() {}

func (x *GetReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReq.ProtoReflect.Descriptor instead.
func (*GetReq) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{2}
}

func (x *GetReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	Value   string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *GetResp) Reset() {
	*x = GetResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResp) ProtoMessage() {}

func (x *GetResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResp.ProtoReflect.Descriptor instead.
func (*GetResp) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{3}
}

func (x *GetResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *GetResp) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// 用于判断幂等性，clientId + seqId唯一标识一个客户端的一次请求
type ReqInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	SeqId    string `protobuf:"bytes,2,opt,name=SeqId,proto3" json:"SeqId,omitempty"`
}

func (x *ReqInfo) Reset() {
	*x = ReqInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqInfo) ProtoMessage() {}

func (x *ReqInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqInfo.ProtoReflect.Descriptor instead.
func (*ReqInfo) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{4}
}

func (x *ReqInfo) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *ReqInfo) GetSeqId() string {
	if x != nil {
		return x.SeqId
	}
	return ""
}

type SetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Info  *ReqInfo `protobuf:"bytes,3,opt,name=Info,proto3" json:"Info,omitempty"`
}

func (x *SetReq) Reset() {
	*x = SetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetReq) ProtoMessage() {}

func (x *SetReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetReq.ProtoReflect.Descriptor instead.
func (*SetReq) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{5}
}

func (x *SetReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetReq) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SetReq) GetInfo() *ReqInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type SetResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *SetResp) Reset() {
	*x = SetResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResp) ProtoMessage() {}

func (x *SetResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResp.ProtoReflect.Descriptor instead.
func (*SetResp) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{6}
}

func (x *SetResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type DelReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key  string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Info *ReqInfo `protobuf:"bytes,3,opt,name=Info,proto3" json:"Info,omitempty"`
}

func (x *DelReq) Reset() {
	*x = DelReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelReq) ProtoMessage() {}

func (x *DelReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelReq.ProtoReflect.Descriptor instead.
func (*DelReq) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{7}
}

func (x *DelReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DelReq) GetInfo() *ReqInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type DelResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	Value   string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *DelResp) Reset() {
	*x = DelResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_data_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelResp) ProtoMessage() {}

func (x *DelResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_data_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelResp.ProtoReflect.Descriptor instead.
func (*DelResp) Descriptor() ([]byte, []int) {
	return file_proto_data_proto_rawDescGZIP(), []int{8}
}

func (x *DelResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *DelResp) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_proto_data_proto protoreflect.FileDescriptor

var file_proto_data_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22,
	0x0d, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x22, 0x2a,
	0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1a,
	0x0a, 0x08, 0x49, 0x73, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x49, 0x73, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0x1a, 0x0a, 0x06, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x39, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x3b, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x65, 0x71, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x65, 0x71, 0x49, 0x64, 0x22, 0x5a,
	0x0a, 0x06, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x28, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x23, 0x0a, 0x07, 0x53, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22,
	0x44, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x04, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x04, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x39, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x32, 0xec, 0x01, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x32, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x1a, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x03, 0x44, 0x65, 0x6c, 0x12,
	0x13, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x19, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x42,
	0x0e, 0x5a, 0x0c, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_data_proto_rawDescOnce sync.Once
	file_proto_data_proto_rawDescData = file_proto_data_proto_rawDesc
)

func file_proto_data_proto_rawDescGZIP() []byte {
	file_proto_data_proto_rawDescOnce.Do(func() {
		file_proto_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_data_proto_rawDescData)
	})
	return file_proto_data_proto_rawDescData
}

var file_proto_data_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_data_proto_goTypes = []interface{}{
	(*GetStateReq)(nil),  // 0: dataService.GetStateReq
	(*GetStateResp)(nil), // 1: dataService.GetStateResp
	(*GetReq)(nil),       // 2: dataService.GetReq
	(*GetResp)(nil),      // 3: dataService.GetResp
	(*ReqInfo)(nil),      // 4: dataService.ReqInfo
	(*SetReq)(nil),       // 5: dataService.SetReq
	(*SetResp)(nil),      // 6: dataService.SetResp
	(*DelReq)(nil),       // 7: dataService.DelReq
	(*DelResp)(nil),      // 8: dataService.DelResp
}
var file_proto_data_proto_depIdxs = []int32{
	4, // 0: dataService.SetReq.Info:type_name -> dataService.ReqInfo
	4, // 1: dataService.DelReq.Info:type_name -> dataService.ReqInfo
	2, // 2: dataService.dataService.Get:input_type -> dataService.GetReq
	5, // 3: dataService.dataService.Set:input_type -> dataService.SetReq
	7, // 4: dataService.dataService.Del:input_type -> dataService.DelReq
	0, // 5: dataService.dataService.GetState:input_type -> dataService.GetStateReq
	3, // 6: dataService.dataService.Get:output_type -> dataService.GetResp
	6, // 7: dataService.dataService.Set:output_type -> dataService.SetResp
	8, // 8: dataService.dataService.Del:output_type -> dataService.DelResp
	1, // 9: dataService.dataService.GetState:output_type -> dataService.GetStateResp
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_data_proto_init() }
func file_proto_data_proto_init() {
	if File_proto_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_data_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_data_proto_goTypes,
		DependencyIndexes: file_proto_data_proto_depIdxs,
		MessageInfos:      file_proto_data_proto_msgTypes,
	}.Build()
	File_proto_data_proto = out.File
	file_proto_data_proto_rawDesc = nil
	file_proto_data_proto_goTypes = nil
	file_proto_data_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DataServiceClient is the client API for DataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataServiceClient interface {
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
	Set(ctx context.Context, in *SetReq, opts ...grpc.CallOption) (*SetResp, error)
	Del(ctx context.Context, in *DelReq, opts ...grpc.CallOption) (*DelResp, error)
	GetState(ctx context.Context, in *GetStateReq, opts ...grpc.CallOption) (*GetStateResp, error)
}

type dataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataServiceClient(cc grpc.ClientConnInterface) DataServiceClient {
	return &dataServiceClient{cc}
}

func (c *dataServiceClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/dataService.dataService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) Set(ctx context.Context, in *SetReq, opts ...grpc.CallOption) (*SetResp, error) {
	out := new(SetResp)
	err := c.cc.Invoke(ctx, "/dataService.dataService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) Del(ctx context.Context, in *DelReq, opts ...grpc.CallOption) (*DelResp, error) {
	out := new(DelResp)
	err := c.cc.Invoke(ctx, "/dataService.dataService/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataServiceClient) GetState(ctx context.Context, in *GetStateReq, opts ...grpc.CallOption) (*GetStateResp, error) {
	out := new(GetStateResp)
	err := c.cc.Invoke(ctx, "/dataService.dataService/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServiceServer is the server API for DataService service.
type DataServiceServer interface {
	Get(context.Context, *GetReq) (*GetResp, error)
	Set(context.Context, *SetReq) (*SetResp, error)
	Del(context.Context, *DelReq) (*DelResp, error)
	GetState(context.Context, *GetStateReq) (*GetStateResp, error)
}

// UnimplementedDataServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDataServiceServer struct {
}

func (*UnimplementedDataServiceServer) Get(context.Context, *GetReq) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedDataServiceServer) Set(context.Context, *SetReq) (*SetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (*UnimplementedDataServiceServer) Del(context.Context, *DelReq) (*DelResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (*UnimplementedDataServiceServer) GetState(context.Context, *GetStateReq) (*GetStateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}

func RegisterDataServiceServer(s *grpc.Server, srv DataServiceServer) {
	s.RegisterService(&_DataService_serviceDesc, srv)
}

func _DataService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataService.dataService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataService.dataService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).Set(ctx, req.(*SetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataService.dataService/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).Del(ctx, req.(*DelReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataService_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServiceServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataService.dataService/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServiceServer).GetState(ctx, req.(*GetStateReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dataService.dataService",
	HandlerType: (*DataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _DataService_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _DataService_Set_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _DataService_Del_Handler,
		},
		{
			MethodName: "GetState",
			Handler:    _DataService_GetState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/data.proto",
}
