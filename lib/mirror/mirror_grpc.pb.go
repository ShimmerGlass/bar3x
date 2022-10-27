// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mirror

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// MirrorClient is the client API for Mirror service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MirrorClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Mirror_SubscribeClient, error)
}

type mirrorClient struct {
	cc grpc.ClientConnInterface
}

func NewMirrorClient(cc grpc.ClientConnInterface) MirrorClient {
	return &mirrorClient{cc}
}

func (c *mirrorClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Mirror_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Mirror_serviceDesc.Streams[0], "/mirror.Mirror/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &mirrorSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Mirror_SubscribeClient interface {
	Recv() (*Image, error)
	grpc.ClientStream
}

type mirrorSubscribeClient struct {
	grpc.ClientStream
}

func (x *mirrorSubscribeClient) Recv() (*Image, error) {
	m := new(Image)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MirrorServer is the server API for Mirror service.
// All implementations must embed UnimplementedMirrorServer
// for forward compatibility
type MirrorServer interface {
	Subscribe(*SubscribeRequest, Mirror_SubscribeServer) error
	mustEmbedUnimplementedMirrorServer()
}

// UnimplementedMirrorServer must be embedded to have forward compatible implementations.
type UnimplementedMirrorServer struct {
}

func (UnimplementedMirrorServer) Subscribe(*SubscribeRequest, Mirror_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedMirrorServer) mustEmbedUnimplementedMirrorServer() {}

// UnsafeMirrorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MirrorServer will
// result in compilation errors.
type UnsafeMirrorServer interface {
	mustEmbedUnimplementedMirrorServer()
}

func RegisterMirrorServer(s grpc.ServiceRegistrar, srv MirrorServer) {
	s.RegisterService(&_Mirror_serviceDesc, srv)
}

func _Mirror_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MirrorServer).Subscribe(m, &mirrorSubscribeServer{stream})
}

type Mirror_SubscribeServer interface {
	Send(*Image) error
	grpc.ServerStream
}

type mirrorSubscribeServer struct {
	grpc.ServerStream
}

func (x *mirrorSubscribeServer) Send(m *Image) error {
	return x.ServerStream.SendMsg(m)
}

var _Mirror_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mirror.Mirror",
	HandlerType: (*MirrorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Mirror_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mirror.proto",
}