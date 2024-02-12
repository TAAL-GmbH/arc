// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: blocktx/blocktx_api/blocktx_api.proto

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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BlockTxAPI_Health_FullMethodName = "/blocktx_api.BlockTxAPI/Health"
)

// BlockTxAPIClient is the client API for BlockTxAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockTxAPIClient interface {
	// Health returns the health of the API.
	Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error)
}

type blockTxAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockTxAPIClient(cc grpc.ClientConnInterface) BlockTxAPIClient {
	return &blockTxAPIClient{cc}
}

func (c *blockTxAPIClient) Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, BlockTxAPI_Health_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockTxAPIServer is the server API for BlockTxAPI service.
// All implementations must embed UnimplementedBlockTxAPIServer
// for forward compatibility
type BlockTxAPIServer interface {
	// Health returns the health of the API.
	Health(context.Context, *emptypb.Empty) (*HealthResponse, error)
	mustEmbedUnimplementedBlockTxAPIServer()
}

// UnimplementedBlockTxAPIServer must be embedded to have forward compatible implementations.
type UnimplementedBlockTxAPIServer struct {
}

func (UnimplementedBlockTxAPIServer) Health(context.Context, *emptypb.Empty) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedBlockTxAPIServer) mustEmbedUnimplementedBlockTxAPIServer() {}

// UnsafeBlockTxAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockTxAPIServer will
// result in compilation errors.
type UnsafeBlockTxAPIServer interface {
	mustEmbedUnimplementedBlockTxAPIServer()
}

func RegisterBlockTxAPIServer(s grpc.ServiceRegistrar, srv BlockTxAPIServer) {
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blocktx/blocktx_api/blocktx_api.proto",
}
