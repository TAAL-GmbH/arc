// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: internal/callbacker/callbacker_api/callbacker_api.proto

package callbacker_api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Note: Values of the statuses have a difference between them in case
// it's necessary to add another values in between. This will allow to
// create new statuses without changing the values in the existing db.
type Status int32

const (
	Status_UNKNOWN                Status = 0
	Status_QUEUED                 Status = 10
	Status_RECEIVED               Status = 20
	Status_STORED                 Status = 30
	Status_ANNOUNCED_TO_NETWORK   Status = 40
	Status_REQUESTED_BY_NETWORK   Status = 50
	Status_SENT_TO_NETWORK        Status = 60
	Status_ACCEPTED_BY_NETWORK    Status = 70
	Status_SEEN_IN_ORPHAN_MEMPOOL Status = 80
	Status_SEEN_ON_NETWORK        Status = 90
	Status_DOUBLE_SPEND_ATTEMPTED Status = 100
	Status_REJECTED               Status = 110
	Status_MINED                  Status = 120
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0:   "UNKNOWN",
		10:  "QUEUED",
		20:  "RECEIVED",
		30:  "STORED",
		40:  "ANNOUNCED_TO_NETWORK",
		50:  "REQUESTED_BY_NETWORK",
		60:  "SENT_TO_NETWORK",
		70:  "ACCEPTED_BY_NETWORK",
		80:  "SEEN_IN_ORPHAN_MEMPOOL",
		90:  "SEEN_ON_NETWORK",
		100: "DOUBLE_SPEND_ATTEMPTED",
		110: "REJECTED",
		120: "MINED",
	}
	Status_value = map[string]int32{
		"UNKNOWN":                0,
		"QUEUED":                 10,
		"RECEIVED":               20,
		"STORED":                 30,
		"ANNOUNCED_TO_NETWORK":   40,
		"REQUESTED_BY_NETWORK":   50,
		"SENT_TO_NETWORK":        60,
		"ACCEPTED_BY_NETWORK":    70,
		"SEEN_IN_ORPHAN_MEMPOOL": 80,
		"SEEN_ON_NETWORK":        90,
		"DOUBLE_SPEND_ATTEMPTED": 100,
		"REJECTED":               110,
		"MINED":                  120,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_internal_callbacker_callbacker_api_callbacker_api_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescGZIP(), []int{0}
}

// swagger:model HealthResponse
type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescGZIP(), []int{0}
}

func (x *HealthResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

// swagger:model SendCallbackRequest
type SendCallbackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CallbackEndpoints []*CallbackEndpoint `protobuf:"bytes,1,rep,name=callback_endpoints,json=callbackEndpoints,proto3" json:"callback_endpoints,omitempty"`
	Txid              string              `protobuf:"bytes,2,opt,name=txid,proto3" json:"txid,omitempty"`
	Status            Status              `protobuf:"varint,3,opt,name=status,proto3,enum=callbacker_api.Status" json:"status,omitempty"`
	MerklePath        string              `protobuf:"bytes,4,opt,name=merkle_path,json=merklePath,proto3" json:"merkle_path,omitempty"`
	ExtraInfo         string              `protobuf:"bytes,5,opt,name=extra_info,json=extraInfo,proto3" json:"extra_info,omitempty"`
	CompetingTxs      []string            `protobuf:"bytes,6,rep,name=competing_txs,json=competingTxs,proto3" json:"competing_txs,omitempty"`
	BlockHash         string              `protobuf:"bytes,7,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockHeight       uint64              `protobuf:"varint,8,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (x *SendCallbackRequest) Reset() {
	*x = SendCallbackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendCallbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCallbackRequest) ProtoMessage() {}

func (x *SendCallbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCallbackRequest.ProtoReflect.Descriptor instead.
func (*SendCallbackRequest) Descriptor() ([]byte, []int) {
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescGZIP(), []int{1}
}

func (x *SendCallbackRequest) GetCallbackEndpoints() []*CallbackEndpoint {
	if x != nil {
		return x.CallbackEndpoints
	}
	return nil
}

func (x *SendCallbackRequest) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

func (x *SendCallbackRequest) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *SendCallbackRequest) GetMerklePath() string {
	if x != nil {
		return x.MerklePath
	}
	return ""
}

func (x *SendCallbackRequest) GetExtraInfo() string {
	if x != nil {
		return x.ExtraInfo
	}
	return ""
}

func (x *SendCallbackRequest) GetCompetingTxs() []string {
	if x != nil {
		return x.CompetingTxs
	}
	return nil
}

func (x *SendCallbackRequest) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *SendCallbackRequest) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

// swagger:model CallbackEndpoint
type CallbackEndpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CallbackEndpoint) Reset() {
	*x = CallbackEndpoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallbackEndpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallbackEndpoint) ProtoMessage() {}

func (x *CallbackEndpoint) ProtoReflect() protoreflect.Message {
	mi := &file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallbackEndpoint.ProtoReflect.Descriptor instead.
func (*CallbackEndpoint) Descriptor() ([]byte, []int) {
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescGZIP(), []int{2}
}

func (x *CallbackEndpoint) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CallbackEndpoint) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_internal_callbacker_callbacker_api_callbacker_api_proto protoreflect.FileDescriptor

var file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDesc = []byte{
	0x0a, 0x37, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x5f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x63, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x22, 0xd1, 0x02, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x6c, 0x6c,
	0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4f, 0x0a, 0x12, 0x63,
	0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63,
	0x6b, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x11, 0x63, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x78, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64,
	0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x50, 0x61, 0x74,
	0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x74, 0x72, 0x61, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x78,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x54, 0x78, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x3a, 0x0a, 0x10, 0x43, 0x61, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x2a, 0x83, 0x02, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b,
	0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x51,
	0x55, 0x45, 0x55, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x43, 0x45, 0x49,
	0x56, 0x45, 0x44, 0x10, 0x14, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x4f, 0x52, 0x45, 0x44, 0x10,
	0x1e, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x4e, 0x4e, 0x4f, 0x55, 0x4e, 0x43, 0x45, 0x44, 0x5f, 0x54,
	0x4f, 0x5f, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x10, 0x28, 0x12, 0x18, 0x0a, 0x14, 0x52,
	0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x45, 0x44, 0x5f, 0x42, 0x59, 0x5f, 0x4e, 0x45, 0x54, 0x57,
	0x4f, 0x52, 0x4b, 0x10, 0x32, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x4f,
	0x5f, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x10, 0x3c, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43,
	0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x5f, 0x42, 0x59, 0x5f, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52,
	0x4b, 0x10, 0x46, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x45, 0x4e, 0x5f, 0x49, 0x4e, 0x5f, 0x4f,
	0x52, 0x50, 0x48, 0x41, 0x4e, 0x5f, 0x4d, 0x45, 0x4d, 0x50, 0x4f, 0x4f, 0x4c, 0x10, 0x50, 0x12,
	0x13, 0x0a, 0x0f, 0x53, 0x45, 0x45, 0x4e, 0x5f, 0x4f, 0x4e, 0x5f, 0x4e, 0x45, 0x54, 0x57, 0x4f,
	0x52, 0x4b, 0x10, 0x5a, 0x12, 0x1a, 0x0a, 0x16, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x5f, 0x53,
	0x50, 0x45, 0x4e, 0x44, 0x5f, 0x41, 0x54, 0x54, 0x45, 0x4d, 0x50, 0x54, 0x45, 0x44, 0x10, 0x64,
	0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x6e, 0x12, 0x09,
	0x0a, 0x05, 0x4d, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x78, 0x32, 0xa2, 0x01, 0x0a, 0x0d, 0x43, 0x61,
	0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x41, 0x50, 0x49, 0x12, 0x42, 0x0a, 0x06, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1e, 0x2e,
	0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4d, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x12,
	0x23, 0x2e, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x12,
	0x5a, 0x10, 0x2e, 0x3b, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x61,
	0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescOnce sync.Once
	file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescData = file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDesc
)

func file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescGZIP() []byte {
	file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescOnce.Do(func() {
		file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescData)
	})
	return file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDescData
}

var file_internal_callbacker_callbacker_api_callbacker_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_callbacker_callbacker_api_callbacker_api_proto_goTypes = []any{
	(Status)(0),                   // 0: callbacker_api.Status
	(*HealthResponse)(nil),        // 1: callbacker_api.HealthResponse
	(*SendCallbackRequest)(nil),   // 2: callbacker_api.SendCallbackRequest
	(*CallbackEndpoint)(nil),      // 3: callbacker_api.CallbackEndpoint
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 5: google.protobuf.Empty
}
var file_internal_callbacker_callbacker_api_callbacker_api_proto_depIdxs = []int32{
	4, // 0: callbacker_api.HealthResponse.timestamp:type_name -> google.protobuf.Timestamp
	3, // 1: callbacker_api.SendCallbackRequest.callback_endpoints:type_name -> callbacker_api.CallbackEndpoint
	0, // 2: callbacker_api.SendCallbackRequest.status:type_name -> callbacker_api.Status
	5, // 3: callbacker_api.CallbackerAPI.Health:input_type -> google.protobuf.Empty
	2, // 4: callbacker_api.CallbackerAPI.SendCallback:input_type -> callbacker_api.SendCallbackRequest
	1, // 5: callbacker_api.CallbackerAPI.Health:output_type -> callbacker_api.HealthResponse
	5, // 6: callbacker_api.CallbackerAPI.SendCallback:output_type -> google.protobuf.Empty
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_callbacker_callbacker_api_callbacker_api_proto_init() }
func file_internal_callbacker_callbacker_api_callbacker_api_proto_init() {
	if File_internal_callbacker_callbacker_api_callbacker_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*HealthResponse); i {
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
		file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SendCallbackRequest); i {
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
		file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CallbackEndpoint); i {
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
			RawDescriptor: file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_callbacker_callbacker_api_callbacker_api_proto_goTypes,
		DependencyIndexes: file_internal_callbacker_callbacker_api_callbacker_api_proto_depIdxs,
		EnumInfos:         file_internal_callbacker_callbacker_api_callbacker_api_proto_enumTypes,
		MessageInfos:      file_internal_callbacker_callbacker_api_callbacker_api_proto_msgTypes,
	}.Build()
	File_internal_callbacker_callbacker_api_callbacker_api_proto = out.File
	file_internal_callbacker_callbacker_api_callbacker_api_proto_rawDesc = nil
	file_internal_callbacker_callbacker_api_callbacker_api_proto_goTypes = nil
	file_internal_callbacker_callbacker_api_callbacker_api_proto_depIdxs = nil
}
