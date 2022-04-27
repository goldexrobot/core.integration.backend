// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: bot.proto

package bot

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

// BotClient is the client API for Bot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BotClient interface {
	// List bots related to the project
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResponse, error)
	// Send bot command
	SendCommand(ctx context.Context, in *SendCommandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type botClient struct {
	cc grpc.ClientConnInterface
}

func NewBotClient(cc grpc.ClientConnInterface) BotClient {
	return &botClient{cc}
}

func (c *botClient) List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/core.backend.integration.api.v1.bot.Bot/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botClient) SendCommand(ctx context.Context, in *SendCommandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/core.backend.integration.api.v1.bot.Bot/SendCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BotServer is the server API for Bot service.
// All implementations must embed UnimplementedBotServer
// for forward compatibility
type BotServer interface {
	// List bots related to the project
	List(context.Context, *emptypb.Empty) (*ListResponse, error)
	// Send bot command
	SendCommand(context.Context, *SendCommandRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedBotServer()
}

// UnimplementedBotServer must be embedded to have forward compatible implementations.
type UnimplementedBotServer struct {
}

func (UnimplementedBotServer) List(context.Context, *emptypb.Empty) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedBotServer) SendCommand(context.Context, *SendCommandRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCommand not implemented")
}
func (UnimplementedBotServer) mustEmbedUnimplementedBotServer() {}

// UnsafeBotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BotServer will
// result in compilation errors.
type UnsafeBotServer interface {
	mustEmbedUnimplementedBotServer()
}

func RegisterBotServer(s grpc.ServiceRegistrar, srv BotServer) {
	s.RegisterService(&Bot_ServiceDesc, srv)
}

func _Bot_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.backend.integration.api.v1.bot.Bot/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServer).List(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bot_SendCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServer).SendCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.backend.integration.api.v1.bot.Bot/SendCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServer).SendCommand(ctx, req.(*SendCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bot_ServiceDesc is the grpc.ServiceDesc for Bot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.backend.integration.api.v1.bot.Bot",
	HandlerType: (*BotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Bot_List_Handler,
		},
		{
			MethodName: "SendCommand",
			Handler:    _Bot_SendCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bot.proto",
}
