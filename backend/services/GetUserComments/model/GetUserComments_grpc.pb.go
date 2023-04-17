// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: model/GetUserComments.proto

package model

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SocialGrpcClient is the client API for SocialGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SocialGrpcClient interface {
	// rpc Chkauth(JwtdataInput) returns (Authd) {}
	// rpc GetAllComments(GetComments) returns (MongoFields) {}
	GetUserComments(ctx context.Context, in *GetComments, opts ...grpc.CallOption) (*MongoFields, error)
}

type socialGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewSocialGrpcClient(cc grpc.ClientConnInterface) SocialGrpcClient {
	return &socialGrpcClient{cc}
}

func (c *socialGrpcClient) GetUserComments(ctx context.Context, in *GetComments, opts ...grpc.CallOption) (*MongoFields, error) {
	out := new(MongoFields)
	err := c.cc.Invoke(ctx, "/SocialGrpc/GetUserComments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SocialGrpcServer is the server API for SocialGrpc service.
// All implementations must embed UnimplementedSocialGrpcServer
// for forward compatibility
type SocialGrpcServer interface {
	// rpc Chkauth(JwtdataInput) returns (Authd) {}
	// rpc GetAllComments(GetComments) returns (MongoFields) {}
	GetUserComments(context.Context, *GetComments) (*MongoFields, error)
	mustEmbedUnimplementedSocialGrpcServer()
}

// UnimplementedSocialGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedSocialGrpcServer struct {
}

func (UnimplementedSocialGrpcServer) GetUserComments(context.Context, *GetComments) (*MongoFields, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserComments not implemented")
}
func (UnimplementedSocialGrpcServer) mustEmbedUnimplementedSocialGrpcServer() {}

// UnsafeSocialGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SocialGrpcServer will
// result in compilation errors.
type UnsafeSocialGrpcServer interface {
	mustEmbedUnimplementedSocialGrpcServer()
}

func RegisterSocialGrpcServer(s grpc.ServiceRegistrar, srv SocialGrpcServer) {
	s.RegisterService(&SocialGrpc_ServiceDesc, srv)
}

func _SocialGrpc_GetUserComments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetComments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialGrpcServer).GetUserComments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SocialGrpc/GetUserComments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialGrpcServer).GetUserComments(ctx, req.(*GetComments))
	}
	return interceptor(ctx, in, info, handler)
}

// SocialGrpc_ServiceDesc is the grpc.ServiceDesc for SocialGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SocialGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SocialGrpc",
	HandlerType: (*SocialGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserComments",
			Handler:    _SocialGrpc_GetUserComments_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model/GetUserComments.proto",
}
