// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package models

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,json=id,proto3" json:"id,omitempty"`
	Nickname             string   `protobuf:"bytes,2,opt,name=Nickname,json=nickname,proto3" json:"nickname,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=Password,json=password,proto3" json:"password,omitempty"`
	Email                string   `protobuf:"bytes,4,opt,name=Email,json=email,proto3" json:"email,omitempty"`
	Fullname             string   `protobuf:"bytes,5,opt,name=Fullname,json=fullname,proto3" json:"fullname,omitempty"`
	Avatar               string   `protobuf:"bytes,6,opt,name=Avatar,json=avatar,proto3" json:"avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_auth_b0e46575ded988fa, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetFullname() string {
	if m != nil {
		return m.Fullname
	}
	return ""
}

func (m *User) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

type SessionId struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,json=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionId) Reset()         { *m = SessionId{} }
func (m *SessionId) String() string { return proto.CompactTextString(m) }
func (*SessionId) ProtoMessage()    {}
func (*SessionId) Descriptor() ([]byte, []int) {
	return fileDescriptor_auth_b0e46575ded988fa, []int{1}
}
func (m *SessionId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionId.Unmarshal(m, b)
}
func (m *SessionId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionId.Marshal(b, m, deterministic)
}
func (dst *SessionId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionId.Merge(dst, src)
}
func (m *SessionId) XXX_Size() int {
	return xxx_messageInfo_SessionId.Size(m)
}
func (m *SessionId) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionId.DiscardUnknown(m)
}

var xxx_messageInfo_SessionId proto.InternalMessageInfo

func (m *SessionId) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Session struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,json=id,proto3" json:"id,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=Token,json=token,proto3" json:"token,omitempty"`
	TTL                  int64    `protobuf:"varint,3,opt,name=TTL,json=tTL,proto3" json:"ttl,omitempty"`
	User                 *User    `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_auth_b0e46575ded988fa, []int{2}
}
func (m *Session) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Session.Unmarshal(m, b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Session.Marshal(b, m, deterministic)
}
func (dst *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(dst, src)
}
func (m *Session) XXX_Size() int {
	return xxx_messageInfo_Session.Size(m)
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Session) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Session) GetTTL() int64 {
	if m != nil {
		return m.TTL
	}
	return 0
}

func (m *Session) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func init() {
	proto.RegisterType((*User)(nil), "models.User")
	proto.RegisterType((*SessionId)(nil), "models.SessionId")
	proto.RegisterType((*Session)(nil), "models.Session")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthClient interface {
	Auth(ctx context.Context, in *User, opts ...grpc.CallOption) (*SessionId, error)
	Check(ctx context.Context, in *SessionId, opts ...grpc.CallOption) (*Session, error)
	LogOut(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error)
}

type authClient struct {
	cc *grpc.ClientConn
}

func NewAuthClient(cc *grpc.ClientConn) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Auth(ctx context.Context, in *User, opts ...grpc.CallOption) (*SessionId, error) {
	out := new(SessionId)
	err := c.cc.Invoke(ctx, "/models.Auth/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Check(ctx context.Context, in *SessionId, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, "/models.Auth/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) LogOut(ctx context.Context, in *Session, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, "/models.Auth/LogOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
type AuthServer interface {
	Auth(context.Context, *User) (*SessionId, error)
	Check(context.Context, *SessionId) (*Session, error)
	LogOut(context.Context, *Session) (*Session, error)
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Auth/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Auth(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Auth/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Check(ctx, req.(*SessionId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_LogOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Session)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).LogOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Auth/LogOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).LogOut(ctx, req.(*Session))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "models.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _Auth_Auth_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _Auth_Check_Handler,
		},
		{
			MethodName: "LogOut",
			Handler:    _Auth_LogOut_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_auth_b0e46575ded988fa) }

var fileDescriptor_auth_b0e46575ded988fa = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xdd, 0x4a, 0x84, 0x40,
	0x14, 0x80, 0xd7, 0x55, 0xa7, 0xf5, 0x14, 0xfd, 0x1c, 0x22, 0xc4, 0x6e, 0xc4, 0xab, 0x0d, 0xca,
	0x8b, 0xed, 0x09, 0x96, 0x28, 0x58, 0x90, 0x0a, 0xb3, 0x07, 0x98, 0x74, 0x4a, 0xf1, 0x67, 0x16,
	0x67, 0xac, 0xb7, 0xe8, 0xba, 0xc7, 0x8d, 0x19, 0xa7, 0x28, 0xeb, 0x4a, 0x3e, 0xbf, 0x73, 0xe4,
	0xe3, 0x08, 0x40, 0x07, 0x59, 0xc6, 0xdb, 0x9e, 0x4b, 0x8e, 0xa4, 0xe5, 0x05, 0x6b, 0x44, 0xf4,
	0x61, 0x81, 0xf3, 0x28, 0x58, 0x8f, 0xfb, 0x30, 0xdf, 0x14, 0xbe, 0x15, 0x5a, 0x4b, 0x3b, 0x9d,
	0x57, 0x05, 0x06, 0xb0, 0xb8, 0xad, 0xf2, 0xba, 0xa3, 0x2d, 0xf3, 0xe7, 0xa1, 0xb5, 0xf4, 0xd2,
	0x45, 0x67, 0x58, 0xb9, 0x7b, 0x2a, 0xc4, 0x1b, 0xef, 0x0b, 0xdf, 0x1e, 0xdd, 0xd6, 0x30, 0x1e,
	0x83, 0x7b, 0xdd, 0xd2, 0xaa, 0xf1, 0x1d, 0x2d, 0x5c, 0xa6, 0x40, 0x6d, 0xdc, 0x0c, 0x4d, 0xa3,
	0xbf, 0xe6, 0x8e, 0x1b, 0xcf, 0x86, 0xf1, 0x04, 0xc8, 0xfa, 0x95, 0x4a, 0xda, 0xfb, 0x44, 0x1b,
	0x42, 0x35, 0x45, 0xa7, 0xe0, 0x3d, 0x30, 0x21, 0x2a, 0xde, 0x6d, 0x8a, 0x1f, 0x79, 0x9e, 0xca,
	0x8b, 0x72, 0xd8, 0x31, 0x72, 0xaa, 0x54, 0x41, 0xc6, 0x6b, 0xd6, 0x99, 0x6c, 0x57, 0x2a, 0xc0,
	0x43, 0xb0, 0xb3, 0x2c, 0xd1, 0xb9, 0x76, 0x6a, 0xcb, 0x2c, 0xc1, 0x10, 0x9c, 0x41, 0xb0, 0x5e,
	0x87, 0xee, 0xae, 0xf6, 0xe2, 0xf1, 0x22, 0xb1, 0xba, 0x46, 0xaa, 0xcd, 0xea, 0xdd, 0x02, 0x67,
	0x3d, 0xc8, 0x12, 0xcf, 0xcc, 0xf3, 0xd7, 0x50, 0x70, 0xf4, 0x45, 0xdf, 0x99, 0xd1, 0x0c, 0x2f,
	0xc0, 0xbd, 0x2a, 0x59, 0x5e, 0xe3, 0x5f, 0x1b, 0x1c, 0x4c, 0x5e, 0x45, 0x33, 0x3c, 0x07, 0x92,
	0xf0, 0x97, 0xbb, 0x41, 0xe2, 0x54, 0xfe, 0x33, 0xfd, 0x44, 0xf4, 0xcf, 0xbb, 0xfc, 0x0c, 0x00,
	0x00, 0xff, 0xff, 0x62, 0x2d, 0x70, 0x7e, 0xca, 0x01, 0x00, 0x00,
}
