// Code generated by protoc-gen-go-handlers_grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-handlers_grpc v1.3.0
// - protoc             v4.24.3
// source: api/shorten.proto

package shorten

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

const (
	Shorten_CreateShortURL_FullMethodName      = "/cart.Shorten/CreateShortURL"
	Shorten_BatchCreateShortURL_FullMethodName = "/cart.Shorten/BatchCreateShortURL"
	Shorten_GetByShort_FullMethodName          = "/cart.Shorten/GetByShort"
	Shorten_GetUserURLs_FullMethodName         = "/cart.Shorten/GetUserURLs"
	Shorten_DeleteUserURLsBatch_FullMethodName = "/cart.Shorten/DeleteUserURLsBatch"
)

// ShortenClient is the client API for Shorten service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenClient interface {
	CreateShortURL(ctx context.Context, in *CreateShortURLRequest, opts ...grpc.CallOption) (*CreateShortURLResponse, error)
	BatchCreateShortURL(ctx context.Context, in *BatchCreateShortURLRequest, opts ...grpc.CallOption) (*BatchCreateShortURLResponse, error)
	GetByShort(ctx context.Context, in *GetByShortRequest, opts ...grpc.CallOption) (*GetByShortResponse, error)
	GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error)
	DeleteUserURLsBatch(ctx context.Context, in *DeleteUserURLsBatchRequest, opts ...grpc.CallOption) (*DeleteUserURLsBatchResponse, error)
}

type shortenClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenClient(cc grpc.ClientConnInterface) ShortenClient {
	return &shortenClient{cc}
}

func (c *shortenClient) CreateShortURL(ctx context.Context, in *CreateShortURLRequest, opts ...grpc.CallOption) (*CreateShortURLResponse, error) {
	out := new(CreateShortURLResponse)
	err := c.cc.Invoke(ctx, Shorten_CreateShortURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenClient) BatchCreateShortURL(ctx context.Context, in *BatchCreateShortURLRequest, opts ...grpc.CallOption) (*BatchCreateShortURLResponse, error) {
	out := new(BatchCreateShortURLResponse)
	err := c.cc.Invoke(ctx, Shorten_BatchCreateShortURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenClient) GetByShort(ctx context.Context, in *GetByShortRequest, opts ...grpc.CallOption) (*GetByShortResponse, error) {
	out := new(GetByShortResponse)
	err := c.cc.Invoke(ctx, Shorten_GetByShort_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenClient) GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error) {
	out := new(GetUserURLsResponse)
	err := c.cc.Invoke(ctx, Shorten_GetUserURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenClient) DeleteUserURLsBatch(ctx context.Context, in *DeleteUserURLsBatchRequest, opts ...grpc.CallOption) (*DeleteUserURLsBatchResponse, error) {
	out := new(DeleteUserURLsBatchResponse)
	err := c.cc.Invoke(ctx, Shorten_DeleteUserURLsBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenServer is the server API for Shorten service.
// All implementations must embed UnimplementedShortenServer
// for forward compatibility
type ShortenServer interface {
	CreateShortURL(context.Context, *CreateShortURLRequest) (*CreateShortURLResponse, error)
	BatchCreateShortURL(context.Context, *BatchCreateShortURLRequest) (*BatchCreateShortURLResponse, error)
	GetByShort(context.Context, *GetByShortRequest) (*GetByShortResponse, error)
	GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error)
	DeleteUserURLsBatch(context.Context, *DeleteUserURLsBatchRequest) (*DeleteUserURLsBatchResponse, error)
	mustEmbedUnimplementedShortenServer()
}

// UnimplementedShortenServer must be embedded to have forward compatible implementations.
type UnimplementedShortenServer struct {
}

func (UnimplementedShortenServer) CreateShortURL(context.Context, *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortURL not implemented")
}
func (UnimplementedShortenServer) BatchCreateShortURL(context.Context, *BatchCreateShortURLRequest) (*BatchCreateShortURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchCreateShortURL not implemented")
}
func (UnimplementedShortenServer) GetByShort(context.Context, *GetByShortRequest) (*GetByShortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByShort not implemented")
}
func (UnimplementedShortenServer) GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedShortenServer) DeleteUserURLsBatch(context.Context, *DeleteUserURLsBatchRequest) (*DeleteUserURLsBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserURLsBatch not implemented")
}
func (UnimplementedShortenServer) mustEmbedUnimplementedShortenServer() {}

// UnsafeShortenServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenServer will
// result in compilation errors.
type UnsafeShortenServer interface {
	mustEmbedUnimplementedShortenServer()
}

func RegisterShortenServer(s grpc.ServiceRegistrar, srv ShortenServer) {
	s.RegisterService(&Shorten_ServiceDesc, srv)
}

func _Shorten_CreateShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShortURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenServer).CreateShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorten_CreateShortURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenServer).CreateShortURL(ctx, req.(*CreateShortURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorten_BatchCreateShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchCreateShortURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenServer).BatchCreateShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorten_BatchCreateShortURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenServer).BatchCreateShortURL(ctx, req.(*BatchCreateShortURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorten_GetByShort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByShortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenServer).GetByShort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorten_GetByShort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenServer).GetByShort(ctx, req.(*GetByShortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorten_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorten_GetUserURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenServer).GetUserURLs(ctx, req.(*GetUserURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorten_DeleteUserURLsBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserURLsBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenServer).DeleteUserURLsBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorten_DeleteUserURLsBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenServer).DeleteUserURLsBatch(ctx, req.(*DeleteUserURLsBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shorten_ServiceDesc is the grpc.ServiceDesc for Shorten service.
// It's only intended for direct use with handlers_grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shorten_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cart.Shorten",
	HandlerType: (*ShortenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortURL",
			Handler:    _Shorten_CreateShortURL_Handler,
		},
		{
			MethodName: "BatchCreateShortURL",
			Handler:    _Shorten_BatchCreateShortURL_Handler,
		},
		{
			MethodName: "GetByShort",
			Handler:    _Shorten_GetByShort_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _Shorten_GetUserURLs_Handler,
		},
		{
			MethodName: "DeleteUserURLsBatch",
			Handler:    _Shorten_DeleteUserURLsBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/shorten.proto",
}
