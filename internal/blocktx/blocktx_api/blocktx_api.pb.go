// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: internal/blocktx/blocktx_api/blocktx_api.proto

package blocktx_api

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

// swagger:model HealthResponse
type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok        bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Details   string                 `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[0]
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
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{0}
}

func (x *HealthResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *HealthResponse) GetDetails() string {
	if x != nil {
		return x.Details
	}
	return ""
}

func (x *HealthResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

// swagger:model Block {
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash         []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`                                     // Little endian
	PreviousHash []byte `protobuf:"bytes,2,opt,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"` // Little endian
	MerkleRoot   []byte `protobuf:"bytes,3,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`       // Little endian
	Height       uint64 `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	Orphaned     bool   `protobuf:"varint,5,opt,name=orphaned,proto3" json:"orphaned,omitempty"`
	Processed    bool   `protobuf:"varint,6,opt,name=processed,proto3" json:"processed,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetPreviousHash() []byte {
	if x != nil {
		return x.PreviousHash
	}
	return nil
}

func (x *Block) GetMerkleRoot() []byte {
	if x != nil {
		return x.MerkleRoot
	}
	return nil
}

func (x *Block) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Block) GetOrphaned() bool {
	if x != nil {
		return x.Orphaned
	}
	return false
}

func (x *Block) GetProcessed() bool {
	if x != nil {
		return x.Processed
	}
	return false
}

// swagger:model Transactions
type Transactions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
}

func (x *Transactions) Reset() {
	*x = Transactions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transactions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transactions) ProtoMessage() {}

func (x *Transactions) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transactions.ProtoReflect.Descriptor instead.
func (*Transactions) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{2}
}

func (x *Transactions) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

type TransactionBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash       []byte `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"` // Little endian
	BlockHeight     uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	TransactionHash []byte `protobuf:"bytes,3,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"` // Little endian
	MerklePath      string `protobuf:"bytes,4,opt,name=merklePath,proto3" json:"merklePath,omitempty"`
}

func (x *TransactionBlock) Reset() {
	*x = TransactionBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionBlock) ProtoMessage() {}

func (x *TransactionBlock) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionBlock.ProtoReflect.Descriptor instead.
func (*TransactionBlock) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{3}
}

func (x *TransactionBlock) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *TransactionBlock) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *TransactionBlock) GetTransactionHash() []byte {
	if x != nil {
		return x.TransactionHash
	}
	return nil
}

func (x *TransactionBlock) GetMerklePath() string {
	if x != nil {
		return x.MerklePath
	}
	return ""
}

type TransactionBlocks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TransactionBlocks []*TransactionBlock `protobuf:"bytes,1,rep,name=transaction_blocks,json=transactionBlocks,proto3" json:"transaction_blocks,omitempty"`
}

func (x *TransactionBlocks) Reset() {
	*x = TransactionBlocks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionBlocks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionBlocks) ProtoMessage() {}

func (x *TransactionBlocks) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionBlocks.ProtoReflect.Descriptor instead.
func (*TransactionBlocks) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{4}
}

func (x *TransactionBlocks) GetTransactionBlocks() []*TransactionBlock {
	if x != nil {
		return x.TransactionBlocks
	}
	return nil
}

// swagger:model Transaction
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash   []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`     // Little endian
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"` // This is the metamorph address:port
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{5}
}

func (x *Transaction) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Transaction) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

// swagger:model ClearData
type ClearData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RetentionDays int32 `protobuf:"varint,1,opt,name=retentionDays,proto3" json:"retentionDays,omitempty"`
}

func (x *ClearData) Reset() {
	*x = ClearData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClearData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearData) ProtoMessage() {}

func (x *ClearData) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearData.ProtoReflect.Descriptor instead.
func (*ClearData) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{6}
}

func (x *ClearData) GetRetentionDays() int32 {
	if x != nil {
		return x.RetentionDays
	}
	return 0
}

// swagger:model ClearDataResponse
type ClearDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rows int64 `protobuf:"varint,1,opt,name=rows,proto3" json:"rows,omitempty"`
}

func (x *ClearDataResponse) Reset() {
	*x = ClearDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClearDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearDataResponse) ProtoMessage() {}

func (x *ClearDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearDataResponse.ProtoReflect.Descriptor instead.
func (*ClearDataResponse) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{7}
}

func (x *ClearDataResponse) GetRows() int64 {
	if x != nil {
		return x.Rows
	}
	return 0
}

type TransactionAndSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash   []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *TransactionAndSource) Reset() {
	*x = TransactionAndSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionAndSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionAndSource) ProtoMessage() {}

func (x *TransactionAndSource) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionAndSource.ProtoReflect.Descriptor instead.
func (*TransactionAndSource) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{8}
}

func (x *TransactionAndSource) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *TransactionAndSource) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

type DelUnfinishedBlockProcessingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessedBy string `protobuf:"bytes,1,opt,name=processed_by,json=processedBy,proto3" json:"processed_by,omitempty"`
}

func (x *DelUnfinishedBlockProcessingRequest) Reset() {
	*x = DelUnfinishedBlockProcessingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelUnfinishedBlockProcessingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelUnfinishedBlockProcessingRequest) ProtoMessage() {}

func (x *DelUnfinishedBlockProcessingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelUnfinishedBlockProcessingRequest.ProtoReflect.Descriptor instead.
func (*DelUnfinishedBlockProcessingRequest) Descriptor() ([]byte, []int) {
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP(), []int{9}
}

func (x *DelUnfinishedBlockProcessingRequest) GetProcessedBy() string {
	if x != nil {
		return x.ProcessedBy
	}
	return ""
}

var File_internal_blocktx_blocktx_api_blocktx_api_proto protoreflect.FileDescriptor

var file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x74, 0x78, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x74, 0x0a, 0x0e, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0xb3, 0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x72,
	0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c,
	0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x6f, 0x72, 0x70, 0x68, 0x61, 0x6e, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x6f, 0x72, 0x70, 0x68, 0x61, 0x6e, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x22, 0x4c, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3c, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x9f, 0x01, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x29, 0x0a, 0x10,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x72,
	0x6b, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x61, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x4c, 0x0a, 0x12,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x22, 0x39, 0x0a, 0x0b, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x31, 0x0a, 0x09, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x44,
	0x61, 0x79, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x72, 0x65, 0x74, 0x65, 0x6e,
	0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x79, 0x73, 0x22, 0x27, 0x0a, 0x11, 0x43, 0x6c, 0x65, 0x61,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x72, 0x6f, 0x77,
	0x73, 0x22, 0x42, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x41, 0x6e, 0x64, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x48, 0x0a, 0x23, 0x44, 0x65, 0x6c, 0x55, 0x6e, 0x66, 0x69,
	0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x42, 0x79, 0x32,
	0xa8, 0x03, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x78, 0x41, 0x50, 0x49, 0x12, 0x3f,
	0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1b, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4d, 0x0a, 0x11, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1e, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47,
	0x0a, 0x0b, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x16, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1e, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x19, 0x43, 0x6c, 0x65, 0x61, 0x72,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x4d, 0x61, 0x70, 0x12, 0x16, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1e, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6a,
	0x0a, 0x1c, 0x44, 0x65, 0x6c, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x12, 0x30,
	0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c,
	0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x3b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74, 0x78, 0x5f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescOnce sync.Once
	file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescData = file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDesc
)

func file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescGZIP() []byte {
	file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescOnce.Do(func() {
		file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescData)
	})
	return file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDescData
}

var file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_internal_blocktx_blocktx_api_blocktx_api_proto_goTypes = []interface{}{
	(*HealthResponse)(nil),                      // 0: blocktx_api.HealthResponse
	(*Block)(nil),                               // 1: blocktx_api.Block
	(*Transactions)(nil),                        // 2: blocktx_api.Transactions
	(*TransactionBlock)(nil),                    // 3: blocktx_api.TransactionBlock
	(*TransactionBlocks)(nil),                   // 4: blocktx_api.TransactionBlocks
	(*Transaction)(nil),                         // 5: blocktx_api.Transaction
	(*ClearData)(nil),                           // 6: blocktx_api.ClearData
	(*ClearDataResponse)(nil),                   // 7: blocktx_api.ClearDataResponse
	(*TransactionAndSource)(nil),                // 8: blocktx_api.TransactionAndSource
	(*DelUnfinishedBlockProcessingRequest)(nil), // 9: blocktx_api.DelUnfinishedBlockProcessingRequest
	(*timestamppb.Timestamp)(nil),               // 10: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                       // 11: google.protobuf.Empty
}
var file_internal_blocktx_blocktx_api_blocktx_api_proto_depIdxs = []int32{
	10, // 0: blocktx_api.HealthResponse.timestamp:type_name -> google.protobuf.Timestamp
	5,  // 1: blocktx_api.Transactions.transactions:type_name -> blocktx_api.Transaction
	3,  // 2: blocktx_api.TransactionBlocks.transaction_blocks:type_name -> blocktx_api.TransactionBlock
	11, // 3: blocktx_api.BlockTxAPI.Health:input_type -> google.protobuf.Empty
	6,  // 4: blocktx_api.BlockTxAPI.ClearTransactions:input_type -> blocktx_api.ClearData
	6,  // 5: blocktx_api.BlockTxAPI.ClearBlocks:input_type -> blocktx_api.ClearData
	6,  // 6: blocktx_api.BlockTxAPI.ClearBlockTransactionsMap:input_type -> blocktx_api.ClearData
	9,  // 7: blocktx_api.BlockTxAPI.DelUnfinishedBlockProcessing:input_type -> blocktx_api.DelUnfinishedBlockProcessingRequest
	0,  // 8: blocktx_api.BlockTxAPI.Health:output_type -> blocktx_api.HealthResponse
	7,  // 9: blocktx_api.BlockTxAPI.ClearTransactions:output_type -> blocktx_api.ClearDataResponse
	7,  // 10: blocktx_api.BlockTxAPI.ClearBlocks:output_type -> blocktx_api.ClearDataResponse
	7,  // 11: blocktx_api.BlockTxAPI.ClearBlockTransactionsMap:output_type -> blocktx_api.ClearDataResponse
	11, // 12: blocktx_api.BlockTxAPI.DelUnfinishedBlockProcessing:output_type -> google.protobuf.Empty
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_internal_blocktx_blocktx_api_blocktx_api_proto_init() }
func file_internal_blocktx_blocktx_api_blocktx_api_proto_init() {
	if File_internal_blocktx_blocktx_api_blocktx_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transactions); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionBlock); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionBlocks); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClearData); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClearDataResponse); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionAndSource); i {
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
		file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelUnfinishedBlockProcessingRequest); i {
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
			RawDescriptor: file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_blocktx_blocktx_api_blocktx_api_proto_goTypes,
		DependencyIndexes: file_internal_blocktx_blocktx_api_blocktx_api_proto_depIdxs,
		MessageInfos:      file_internal_blocktx_blocktx_api_blocktx_api_proto_msgTypes,
	}.Build()
	File_internal_blocktx_blocktx_api_blocktx_api_proto = out.File
	file_internal_blocktx_blocktx_api_blocktx_api_proto_rawDesc = nil
	file_internal_blocktx_blocktx_api_blocktx_api_proto_goTypes = nil
	file_internal_blocktx_blocktx_api_blocktx_api_proto_depIdxs = nil
}
