// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: steaming/streaming.proto

package streaming

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

// DataStreamingClient is the client API for DataStreaming service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataStreamingClient interface {
	UploadData(ctx context.Context, opts ...grpc.CallOption) (DataStreaming_UploadDataClient, error)
}

type dataStreamingClient struct {
	cc grpc.ClientConnInterface
}

func NewDataStreamingClient(cc grpc.ClientConnInterface) DataStreamingClient {
	return &dataStreamingClient{cc}
}

func (c *dataStreamingClient) UploadData(ctx context.Context, opts ...grpc.CallOption) (DataStreaming_UploadDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &DataStreaming_ServiceDesc.Streams[0], "/streaming.DataStreaming/UploadData", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataStreamingUploadDataClient{stream}
	return x, nil
}

type DataStreaming_UploadDataClient interface {
	Send(*UploadDataRequest) error
	CloseAndRecv() (*UploadDataResponse, error)
	grpc.ClientStream
}

type dataStreamingUploadDataClient struct {
	grpc.ClientStream
}

func (x *dataStreamingUploadDataClient) Send(m *UploadDataRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dataStreamingUploadDataClient) CloseAndRecv() (*UploadDataResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DataStreamingServer is the server API for DataStreaming service.
// All implementations must embed UnimplementedDataStreamingServer
// for forward compatibility
type DataStreamingServer interface {
	UploadData(DataStreaming_UploadDataServer) error
	mustEmbedUnimplementedDataStreamingServer()
}

// UnimplementedDataStreamingServer must be embedded to have forward compatible implementations.
type UnimplementedDataStreamingServer struct {
}

func (UnimplementedDataStreamingServer) UploadData(DataStreaming_UploadDataServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadData not implemented")
}
func (UnimplementedDataStreamingServer) mustEmbedUnimplementedDataStreamingServer() {}

// UnsafeDataStreamingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataStreamingServer will
// result in compilation errors.
type UnsafeDataStreamingServer interface {
	mustEmbedUnimplementedDataStreamingServer()
}

func RegisterDataStreamingServer(s grpc.ServiceRegistrar, srv DataStreamingServer) {
	s.RegisterService(&DataStreaming_ServiceDesc, srv)
}

func _DataStreaming_UploadData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataStreamingServer).UploadData(&dataStreamingUploadDataServer{stream})
}

type DataStreaming_UploadDataServer interface {
	SendAndClose(*UploadDataResponse) error
	Recv() (*UploadDataRequest, error)
	grpc.ServerStream
}

type dataStreamingUploadDataServer struct {
	grpc.ServerStream
}

func (x *dataStreamingUploadDataServer) SendAndClose(m *UploadDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dataStreamingUploadDataServer) Recv() (*UploadDataRequest, error) {
	m := new(UploadDataRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DataStreaming_ServiceDesc is the grpc.ServiceDesc for DataStreaming service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataStreaming_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "streaming.DataStreaming",
	HandlerType: (*DataStreamingServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadData",
			Handler:       _DataStreaming_UploadData_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "steaming/streaming.proto",
}
