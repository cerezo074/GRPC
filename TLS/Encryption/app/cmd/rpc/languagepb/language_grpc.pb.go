// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package languagepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// LanguageServiceClient is the client API for LanguageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LanguageServiceClient interface {
	//Unary
	Detect(ctx context.Context, in *LanguageDetectorRequest, opts ...grpc.CallOption) (*LanguageDetectorResponse, error)
}

type languageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLanguageServiceClient(cc grpc.ClientConnInterface) LanguageServiceClient {
	return &languageServiceClient{cc}
}

func (c *languageServiceClient) Detect(ctx context.Context, in *LanguageDetectorRequest, opts ...grpc.CallOption) (*LanguageDetectorResponse, error) {
	out := new(LanguageDetectorResponse)
	err := c.cc.Invoke(ctx, "/LanguageService/Detect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LanguageServiceServer is the server API for LanguageService service.
// All implementations must embed UnimplementedLanguageServiceServer
// for forward compatibility
type LanguageServiceServer interface {
	//Unary
	Detect(context.Context, *LanguageDetectorRequest) (*LanguageDetectorResponse, error)
	mustEmbedUnimplementedLanguageServiceServer()
}

// UnimplementedLanguageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLanguageServiceServer struct {
}

func (UnimplementedLanguageServiceServer) Detect(context.Context, *LanguageDetectorRequest) (*LanguageDetectorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detect not implemented")
}
func (UnimplementedLanguageServiceServer) mustEmbedUnimplementedLanguageServiceServer() {}

// UnsafeLanguageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LanguageServiceServer will
// result in compilation errors.
type UnsafeLanguageServiceServer interface {
	mustEmbedUnimplementedLanguageServiceServer()
}

func RegisterLanguageServiceServer(s grpc.ServiceRegistrar, srv LanguageServiceServer) {
	s.RegisterService(&_LanguageService_serviceDesc, srv)
}

func _LanguageService_Detect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LanguageDetectorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LanguageServiceServer).Detect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/LanguageService/Detect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LanguageServiceServer).Detect(ctx, req.(*LanguageDetectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LanguageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "LanguageService",
	HandlerType: (*LanguageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Detect",
			Handler:    _LanguageService_Detect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/language.proto",
}