// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: notes.proto

package notes

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

// NoteServiceClient is the client API for NoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NoteServiceClient interface {
	GetNotes(ctx context.Context, in *NoteRequest, opts ...grpc.CallOption) (*NoteResponse, error)
}

type noteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNoteServiceClient(cc grpc.ClientConnInterface) NoteServiceClient {
	return &noteServiceClient{cc}
}

func (c *noteServiceClient) GetNotes(ctx context.Context, in *NoteRequest, opts ...grpc.CallOption) (*NoteResponse, error) {
	out := new(NoteResponse)
	err := c.cc.Invoke(ctx, "/notes.NoteService/GetNotes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NoteServiceServer is the server API for NoteService service.
// All implementations must embed UnimplementedNoteServiceServer
// for forward compatibility
type NoteServiceServer interface {
	GetNotes(context.Context, *NoteRequest) (*NoteResponse, error)
	mustEmbedUnimplementedNoteServiceServer()
}

// UnimplementedNoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNoteServiceServer struct {
}

func (UnimplementedNoteServiceServer) GetNotes(context.Context, *NoteRequest) (*NoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotes not implemented")
}
func (UnimplementedNoteServiceServer) mustEmbedUnimplementedNoteServiceServer() {}

// UnsafeNoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NoteServiceServer will
// result in compilation errors.
type UnsafeNoteServiceServer interface {
	mustEmbedUnimplementedNoteServiceServer()
}

func RegisterNoteServiceServer(s grpc.ServiceRegistrar, srv NoteServiceServer) {
	s.RegisterService(&NoteService_ServiceDesc, srv)
}

func _NoteService_GetNotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServiceServer).GetNotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notes.NoteService/GetNotes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServiceServer).GetNotes(ctx, req.(*NoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NoteService_ServiceDesc is the grpc.ServiceDesc for NoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notes.NoteService",
	HandlerType: (*NoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNotes",
			Handler:    _NoteService_GetNotes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notes.proto",
}