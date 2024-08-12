// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: internal/blocktx/blocktx_api/blocktx_api.proto

package blocktx_api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BlockTxAPI_Health_FullMethodName                       = "/blocktx_api.BlockTxAPI/Health"
	BlockTxAPI_ClearTransactions_FullMethodName            = "/blocktx_api.BlockTxAPI/ClearTransactions"
	BlockTxAPI_ClearBlocks_FullMethodName                  = "/blocktx_api.BlockTxAPI/ClearBlocks"
	BlockTxAPI_ClearBlockTransactionsMap_FullMethodName    = "/blocktx_api.BlockTxAPI/ClearBlockTransactionsMap"
	BlockTxAPI_DelUnfinishedBlockProcessing_FullMethodName = "/blocktx_api.BlockTxAPI/DelUnfinishedBlockProcessing"
	BlockTxAPI_VerifyMerkleRoots_FullMethodName            = "/blocktx_api.BlockTxAPI/VerifyMerkleRoots"
)

// BlockTxAPIClient is the client API for BlockTxAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockTxAPIClient interface {
	// Health returns the health of the API.
	Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error)
	// ClearTransactions clears transaction data
	ClearTransactions(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error)
	// ClearBlocks clears block data
	ClearBlocks(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error)
	// ClearBlockTransactionsMap clears block-transaction-map data
	ClearBlockTransactionsMap(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error)
	// DelUnfinishedBlockProcessing deletes unfinished block processing
	DelUnfinishedBlockProcessing(ctx context.Context, in *DelUnfinishedBlockProcessingRequest, opts ...grpc.CallOption) (*RowsAffectedResponse, error)
	// VerifyMerkleRoots verifies the merkle roots existance in blocktx db and returns unverified block heights
	VerifyMerkleRoots(ctx context.Context, in *MerkleRootsVerificationRequest, opts ...grpc.CallOption) (*MerkleRootVerificationResponse, error)
}

type blockTxAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockTxAPIClient(cc grpc.ClientConnInterface) BlockTxAPIClient {
	return &blockTxAPIClient{cc}
}

func (c *blockTxAPIClient) Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockTxAPIClient) ClearTransactions(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RowsAffectedResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_ClearTransactions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockTxAPIClient) ClearBlocks(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RowsAffectedResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_ClearBlocks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockTxAPIClient) ClearBlockTransactionsMap(ctx context.Context, in *ClearData, opts ...grpc.CallOption) (*RowsAffectedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RowsAffectedResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_ClearBlockTransactionsMap_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockTxAPIClient) DelUnfinishedBlockProcessing(ctx context.Context, in *DelUnfinishedBlockProcessingRequest, opts ...grpc.CallOption) (*RowsAffectedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RowsAffectedResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_DelUnfinishedBlockProcessing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockTxAPIClient) VerifyMerkleRoots(ctx context.Context, in *MerkleRootsVerificationRequest, opts ...grpc.CallOption) (*MerkleRootVerificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MerkleRootVerificationResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_VerifyMerkleRoots_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockTxAPIServer is the server API for BlockTxAPI service.
// All implementations must embed UnimplementedBlockTxAPIServer
// for forward compatibility.
type BlockTxAPIServer interface {
	// Health returns the health of the API.
	Health(context.Context, *emptypb.Empty) (*HealthResponse, error)
	// ClearTransactions clears transaction data
	ClearTransactions(context.Context, *ClearData) (*RowsAffectedResponse, error)
	// ClearBlocks clears block data
	ClearBlocks(context.Context, *ClearData) (*RowsAffectedResponse, error)
	// ClearBlockTransactionsMap clears block-transaction-map data
	ClearBlockTransactionsMap(context.Context, *ClearData) (*RowsAffectedResponse, error)
	// DelUnfinishedBlockProcessing deletes unfinished block processing
	DelUnfinishedBlockProcessing(context.Context, *DelUnfinishedBlockProcessingRequest) (*RowsAffectedResponse, error)
	// VerifyMerkleRoots verifies the merkle roots existance in blocktx db and returns unverified block heights
	VerifyMerkleRoots(context.Context, *MerkleRootsVerificationRequest) (*MerkleRootVerificationResponse, error)
	mustEmbedUnimplementedBlockTxAPIServer()
}

// UnimplementedBlockTxAPIServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBlockTxAPIServer struct{}

func (UnimplementedBlockTxAPIServer) Health(context.Context, *emptypb.Empty) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedBlockTxAPIServer) ClearTransactions(context.Context, *ClearData) (*RowsAffectedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearTransactions not implemented")
}
func (UnimplementedBlockTxAPIServer) ClearBlocks(context.Context, *ClearData) (*RowsAffectedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearBlocks not implemented")
}
func (UnimplementedBlockTxAPIServer) ClearBlockTransactionsMap(context.Context, *ClearData) (*RowsAffectedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearBlockTransactionsMap not implemented")
}
func (UnimplementedBlockTxAPIServer) DelUnfinishedBlockProcessing(context.Context, *DelUnfinishedBlockProcessingRequest) (*RowsAffectedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelUnfinishedBlockProcessing not implemented")
}
func (UnimplementedBlockTxAPIServer) VerifyMerkleRoots(context.Context, *MerkleRootsVerificationRequest) (*MerkleRootVerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyMerkleRoots not implemented")
}
func (UnimplementedBlockTxAPIServer) mustEmbedUnimplementedBlockTxAPIServer() {}
func (UnimplementedBlockTxAPIServer) testEmbeddedByValue()                    {}

// UnsafeBlockTxAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockTxAPIServer will
// result in compilation errors.
type UnsafeBlockTxAPIServer interface {
	mustEmbedUnimplementedBlockTxAPIServer()
}

func RegisterBlockTxAPIServer(s grpc.ServiceRegistrar, srv BlockTxAPIServer) {
	// If the following call pancis, it indicates UnimplementedBlockTxAPIServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BlockTxAPI_ServiceDesc, srv)
}

func _BlockTxAPI_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).Health(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockTxAPI_ClearTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).ClearTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_ClearTransactions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).ClearTransactions(ctx, req.(*ClearData))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockTxAPI_ClearBlocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).ClearBlocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_ClearBlocks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).ClearBlocks(ctx, req.(*ClearData))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockTxAPI_ClearBlockTransactionsMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).ClearBlockTransactionsMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_ClearBlockTransactionsMap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).ClearBlockTransactionsMap(ctx, req.(*ClearData))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockTxAPI_DelUnfinishedBlockProcessing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelUnfinishedBlockProcessingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).DelUnfinishedBlockProcessing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_DelUnfinishedBlockProcessing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).DelUnfinishedBlockProcessing(ctx, req.(*DelUnfinishedBlockProcessingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockTxAPI_VerifyMerkleRoots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MerkleRootsVerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockTxAPIServer).VerifyMerkleRoots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockTxAPI_VerifyMerkleRoots_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockTxAPIServer).VerifyMerkleRoots(ctx, req.(*MerkleRootsVerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockTxAPI_ServiceDesc is the grpc.ServiceDesc for BlockTxAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockTxAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blocktx_api.BlockTxAPI",
	HandlerType: (*BlockTxAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _BlockTxAPI_Health_Handler,
		},
		{
			MethodName: "ClearTransactions",
			Handler:    _BlockTxAPI_ClearTransactions_Handler,
		},
		{
			MethodName: "ClearBlocks",
			Handler:    _BlockTxAPI_ClearBlocks_Handler,
		},
		{
			MethodName: "ClearBlockTransactionsMap",
			Handler:    _BlockTxAPI_ClearBlockTransactionsMap_Handler,
		},
		{
			MethodName: "DelUnfinishedBlockProcessing",
			Handler:    _BlockTxAPI_DelUnfinishedBlockProcessing_Handler,
		},
		{
			MethodName: "VerifyMerkleRoots",
			Handler:    _BlockTxAPI_VerifyMerkleRoots_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/blocktx/blocktx_api/blocktx_api.proto",
}
