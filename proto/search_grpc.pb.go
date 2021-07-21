// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.14.0
// source: proto/search.proto

package proto

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

// SearcherClient is the client API for Searcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearcherClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	Details(ctx context.Context, in *DetailsRequest, opts ...grpc.CallOption) (*DetailsResponse, error)
	Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error)
	SearchKeyword(ctx context.Context, in *SearchKeywordRequest, opts ...grpc.CallOption) (*SearchKeywordResponse, error)
	Keywords(ctx context.Context, in *KeywordsRequest, opts ...grpc.CallOption) (*KeywordsResponse, error)
}

type searcherClient struct {
	cc grpc.ClientConnInterface
}

func NewSearcherClient(cc grpc.ClientConnInterface) SearcherClient {
	return &searcherClient{cc}
}

func (c *searcherClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/search.Searcher/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searcherClient) Details(ctx context.Context, in *DetailsRequest, opts ...grpc.CallOption) (*DetailsResponse, error) {
	out := new(DetailsResponse)
	err := c.cc.Invoke(ctx, "/search.Searcher/Details", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searcherClient) Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error) {
	out := new(DetailResponse)
	err := c.cc.Invoke(ctx, "/search.Searcher/Detail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searcherClient) SearchKeyword(ctx context.Context, in *SearchKeywordRequest, opts ...grpc.CallOption) (*SearchKeywordResponse, error) {
	out := new(SearchKeywordResponse)
	err := c.cc.Invoke(ctx, "/search.Searcher/SearchKeyword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searcherClient) Keywords(ctx context.Context, in *KeywordsRequest, opts ...grpc.CallOption) (*KeywordsResponse, error) {
	out := new(KeywordsResponse)
	err := c.cc.Invoke(ctx, "/search.Searcher/Keywords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearcherServer is the server API for Searcher service.
// All implementations must embed UnimplementedSearcherServer
// for forward compatibility
type SearcherServer interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	Details(context.Context, *DetailsRequest) (*DetailsResponse, error)
	Detail(context.Context, *DetailRequest) (*DetailResponse, error)
	SearchKeyword(context.Context, *SearchKeywordRequest) (*SearchKeywordResponse, error)
	Keywords(context.Context, *KeywordsRequest) (*KeywordsResponse, error)
}

// UnimplementedSearcherServer must be embedded to have forward compatible implementations.
type UnimplementedSearcherServer struct {
}

func (UnimplementedSearcherServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearcherServer) Details(context.Context, *DetailsRequest) (*DetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Details not implemented")
}
func (UnimplementedSearcherServer) Detail(context.Context, *DetailRequest) (*DetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detail not implemented")
}
func (UnimplementedSearcherServer) SearchKeyword(context.Context, *SearchKeywordRequest) (*SearchKeywordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchKeyword not implemented")
}
func (UnimplementedSearcherServer) Keywords(context.Context, *KeywordsRequest) (*KeywordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keywords not implemented")
}
func (UnimplementedSearcherServer) mustEmbedUnimplementedSearcherServer() {}

// UnsafeSearcherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearcherServer will
// result in compilation errors.
type UnsafeSearcherServer interface {
	mustEmbedUnimplementedSearcherServer()
}

func RegisterSearcherServer(s grpc.ServiceRegistrar, srv SearcherServer) {
	s.RegisterService(&Searcher_ServiceDesc, srv)
}

func _Searcher_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearcherServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Searcher/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearcherServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Searcher_Details_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearcherServer).Details(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Searcher/Details",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearcherServer).Details(ctx, req.(*DetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Searcher_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearcherServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Searcher/Detail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearcherServer).Detail(ctx, req.(*DetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Searcher_SearchKeyword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchKeywordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearcherServer).SearchKeyword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Searcher/SearchKeyword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearcherServer).SearchKeyword(ctx, req.(*SearchKeywordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Searcher_Keywords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeywordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearcherServer).Keywords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Searcher/Keywords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearcherServer).Keywords(ctx, req.(*KeywordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Searcher_ServiceDesc is the grpc.ServiceDesc for Searcher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Searcher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.Searcher",
	HandlerType: (*SearcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _Searcher_Search_Handler,
		},
		{
			MethodName: "Details",
			Handler:    _Searcher_Details_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _Searcher_Detail_Handler,
		},
		{
			MethodName: "SearchKeyword",
			Handler:    _Searcher_SearchKeyword_Handler,
		},
		{
			MethodName: "Keywords",
			Handler:    _Searcher_Keywords_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/search.proto",
}