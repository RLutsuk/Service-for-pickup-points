// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: pvz.proto

package pickuppoint

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PPService_GetPickupPointList_FullMethodName = "/pickuppoint.PPService/GetPickupPointList"
)

// PPServiceClient is the client API for PPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PPServiceClient interface {
	GetPickupPointList(ctx context.Context, in *GetPickupPointRequest, opts ...grpc.CallOption) (*GetPickupPointResponse, error)
}

type pPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPPServiceClient(cc grpc.ClientConnInterface) PPServiceClient {
	return &pPServiceClient{cc}
}

func (c *pPServiceClient) GetPickupPointList(ctx context.Context, in *GetPickupPointRequest, opts ...grpc.CallOption) (*GetPickupPointResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPickupPointResponse)
	err := c.cc.Invoke(ctx, PPService_GetPickupPointList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PPServiceServer is the server API for PPService service.
// All implementations must embed UnimplementedPPServiceServer
// for forward compatibility.
type PPServiceServer interface {
	GetPickupPointList(context.Context, *GetPickupPointRequest) (*GetPickupPointResponse, error)
	mustEmbedUnimplementedPPServiceServer()
}

// UnimplementedPPServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPPServiceServer struct{}

func (UnimplementedPPServiceServer) GetPickupPointList(context.Context, *GetPickupPointRequest) (*GetPickupPointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPickupPointList not implemented")
}
func (UnimplementedPPServiceServer) mustEmbedUnimplementedPPServiceServer() {}
func (UnimplementedPPServiceServer) testEmbeddedByValue()                   {}

// UnsafePPServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PPServiceServer will
// result in compilation errors.
type UnsafePPServiceServer interface {
	mustEmbedUnimplementedPPServiceServer()
}

func RegisterPPServiceServer(s grpc.ServiceRegistrar, srv PPServiceServer) {
	// If the following call pancis, it indicates UnimplementedPPServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PPService_ServiceDesc, srv)
}

func _PPService_GetPickupPointList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPickupPointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PPServiceServer).GetPickupPointList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PPService_GetPickupPointList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PPServiceServer).GetPickupPointList(ctx, req.(*GetPickupPointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PPService_ServiceDesc is the grpc.ServiceDesc for PPService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pickuppoint.PPService",
	HandlerType: (*PPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPickupPointList",
			Handler:    _PPService_GetPickupPointList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pvz.proto",
}
