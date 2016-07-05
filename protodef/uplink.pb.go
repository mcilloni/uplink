// Code generated by protoc-gen-go.
// source: uplink.proto
// DO NOT EDIT!

/*
Package protodef is a generated protocol buffer package.

It is generated from these files:
	uplink.proto

It has these top-level messages:
	Empty
	BoolResp
	NewUserReq
	NewUserResp
	LoginReq
	Username
	AuthInfo
	Challenge
	LoginResp
	LoginAccepted
	UserInfo
	SessInfo
	Notification
*/
package protodef

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

type Notification_Type int32

const (
	Notification_MESSAGE Notification_Type = 0
	Notification_REQUEST Notification_Type = 1
	Notification_INVITE  Notification_Type = 2
)

var Notification_Type_name = map[int32]string{
	0: "MESSAGE",
	1: "REQUEST",
	2: "INVITE",
}
var Notification_Type_value = map[string]int32{
	"MESSAGE": 0,
	"REQUEST": 1,
	"INVITE":  2,
}

func (x Notification_Type) String() string {
	return proto.EnumName(Notification_Type_name, int32(x))
}
func (Notification_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{12, 0} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type BoolResp struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *BoolResp) Reset()                    { *m = BoolResp{} }
func (m *BoolResp) String() string            { return proto.CompactTextString(m) }
func (*BoolResp) ProtoMessage()               {}
func (*BoolResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type NewUserReq struct {
	Name          string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Pass          string `protobuf:"bytes,2,opt,name=pass" json:"pass,omitempty"`
	PublicKey     []byte `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	EncPrivateKey []byte `protobuf:"bytes,4,opt,name=enc_private_key,json=encPrivateKey,proto3" json:"enc_private_key,omitempty"`
	KeyIv         []byte `protobuf:"bytes,5,opt,name=key_iv,json=keyIv,proto3" json:"key_iv,omitempty"`
	KeySalt       []byte `protobuf:"bytes,6,opt,name=key_salt,json=keySalt,proto3" json:"key_salt,omitempty"`
}

func (m *NewUserReq) Reset()                    { *m = NewUserReq{} }
func (m *NewUserReq) String() string            { return proto.CompactTextString(m) }
func (*NewUserReq) ProtoMessage()               {}
func (*NewUserReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type NewUserResp struct {
	SessionInfo *SessInfo `protobuf:"bytes,1,opt,name=session_info,json=sessionInfo" json:"session_info,omitempty"`
}

func (m *NewUserResp) Reset()                    { *m = NewUserResp{} }
func (m *NewUserResp) String() string            { return proto.CompactTextString(m) }
func (*NewUserResp) ProtoMessage()               {}
func (*NewUserResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NewUserResp) GetSessionInfo() *SessInfo {
	if m != nil {
		return m.SessionInfo
	}
	return nil
}

type LoginReq struct {
	// Types that are valid to be assigned to LoginSteps:
	//	*LoginReq_Step1
	//	*LoginReq_Step2
	LoginSteps isLoginReq_LoginSteps `protobuf_oneof:"login_steps"`
}

func (m *LoginReq) Reset()                    { *m = LoginReq{} }
func (m *LoginReq) String() string            { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()               {}
func (*LoginReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isLoginReq_LoginSteps interface {
	isLoginReq_LoginSteps()
}

type LoginReq_Step1 struct {
	Step1 *AuthInfo `protobuf:"bytes,1,opt,name=step1,oneof"`
}
type LoginReq_Step2 struct {
	Step2 *Challenge `protobuf:"bytes,2,opt,name=step2,oneof"`
}

func (*LoginReq_Step1) isLoginReq_LoginSteps() {}
func (*LoginReq_Step2) isLoginReq_LoginSteps() {}

func (m *LoginReq) GetLoginSteps() isLoginReq_LoginSteps {
	if m != nil {
		return m.LoginSteps
	}
	return nil
}

func (m *LoginReq) GetStep1() *AuthInfo {
	if x, ok := m.GetLoginSteps().(*LoginReq_Step1); ok {
		return x.Step1
	}
	return nil
}

func (m *LoginReq) GetStep2() *Challenge {
	if x, ok := m.GetLoginSteps().(*LoginReq_Step2); ok {
		return x.Step2
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoginReq) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoginReq_OneofMarshaler, _LoginReq_OneofUnmarshaler, _LoginReq_OneofSizer, []interface{}{
		(*LoginReq_Step1)(nil),
		(*LoginReq_Step2)(nil),
	}
}

func _LoginReq_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoginReq)
	// login_steps
	switch x := m.LoginSteps.(type) {
	case *LoginReq_Step1:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Step1); err != nil {
			return err
		}
	case *LoginReq_Step2:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Step2); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoginReq.LoginSteps has unexpected type %T", x)
	}
	return nil
}

func _LoginReq_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoginReq)
	switch tag {
	case 1: // login_steps.step1
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(AuthInfo)
		err := b.DecodeMessage(msg)
		m.LoginSteps = &LoginReq_Step1{msg}
		return true, err
	case 2: // login_steps.step2
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Challenge)
		err := b.DecodeMessage(msg)
		m.LoginSteps = &LoginReq_Step2{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoginReq_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoginReq)
	// login_steps
	switch x := m.LoginSteps.(type) {
	case *LoginReq_Step1:
		s := proto.Size(x.Step1)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoginReq_Step2:
		s := proto.Size(x.Step2)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Username struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *Username) Reset()                    { *m = Username{} }
func (m *Username) String() string            { return proto.CompactTextString(m) }
func (*Username) ProtoMessage()               {}
func (*Username) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type AuthInfo struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Pass string `protobuf:"bytes,2,opt,name=pass" json:"pass,omitempty"`
}

func (m *AuthInfo) Reset()                    { *m = AuthInfo{} }
func (m *AuthInfo) String() string            { return proto.CompactTextString(m) }
func (*AuthInfo) ProtoMessage()               {}
func (*AuthInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type Challenge struct {
	Token []byte `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *Challenge) Reset()                    { *m = Challenge{} }
func (m *Challenge) String() string            { return proto.CompactTextString(m) }
func (*Challenge) ProtoMessage()               {}
func (*Challenge) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type LoginResp struct {
	// Types that are valid to be assigned to LoginSteps:
	//	*LoginResp_Step1
	//	*LoginResp_Step2
	LoginSteps isLoginResp_LoginSteps `protobuf_oneof:"login_steps"`
}

func (m *LoginResp) Reset()                    { *m = LoginResp{} }
func (m *LoginResp) String() string            { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()               {}
func (*LoginResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type isLoginResp_LoginSteps interface {
	isLoginResp_LoginSteps()
}

type LoginResp_Step1 struct {
	Step1 *LoginAccepted `protobuf:"bytes,1,opt,name=step1,oneof"`
}
type LoginResp_Step2 struct {
	Step2 *SessInfo `protobuf:"bytes,2,opt,name=step2,oneof"`
}

func (*LoginResp_Step1) isLoginResp_LoginSteps() {}
func (*LoginResp_Step2) isLoginResp_LoginSteps() {}

func (m *LoginResp) GetLoginSteps() isLoginResp_LoginSteps {
	if m != nil {
		return m.LoginSteps
	}
	return nil
}

func (m *LoginResp) GetStep1() *LoginAccepted {
	if x, ok := m.GetLoginSteps().(*LoginResp_Step1); ok {
		return x.Step1
	}
	return nil
}

func (m *LoginResp) GetStep2() *SessInfo {
	if x, ok := m.GetLoginSteps().(*LoginResp_Step2); ok {
		return x.Step2
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoginResp) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoginResp_OneofMarshaler, _LoginResp_OneofUnmarshaler, _LoginResp_OneofSizer, []interface{}{
		(*LoginResp_Step1)(nil),
		(*LoginResp_Step2)(nil),
	}
}

func _LoginResp_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoginResp)
	// login_steps
	switch x := m.LoginSteps.(type) {
	case *LoginResp_Step1:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Step1); err != nil {
			return err
		}
	case *LoginResp_Step2:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Step2); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoginResp.LoginSteps has unexpected type %T", x)
	}
	return nil
}

func _LoginResp_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoginResp)
	switch tag {
	case 1: // login_steps.step1
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LoginAccepted)
		err := b.DecodeMessage(msg)
		m.LoginSteps = &LoginResp_Step1{msg}
		return true, err
	case 2: // login_steps.step2
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SessInfo)
		err := b.DecodeMessage(msg)
		m.LoginSteps = &LoginResp_Step2{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoginResp_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoginResp)
	// login_steps
	switch x := m.LoginSteps.(type) {
	case *LoginResp_Step1:
		s := proto.Size(x.Step1)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoginResp_Step2:
		s := proto.Size(x.Step2)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type LoginAccepted struct {
	UserInfo  *UserInfo  `protobuf:"bytes,1,opt,name=user_info,json=userInfo" json:"user_info,omitempty"`
	Challenge *Challenge `protobuf:"bytes,2,opt,name=challenge" json:"challenge,omitempty"`
}

func (m *LoginAccepted) Reset()                    { *m = LoginAccepted{} }
func (m *LoginAccepted) String() string            { return proto.CompactTextString(m) }
func (*LoginAccepted) ProtoMessage()               {}
func (*LoginAccepted) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *LoginAccepted) GetUserInfo() *UserInfo {
	if m != nil {
		return m.UserInfo
	}
	return nil
}

func (m *LoginAccepted) GetChallenge() *Challenge {
	if m != nil {
		return m.Challenge
	}
	return nil
}

type UserInfo struct {
	PublicKey     []byte `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	EncPrivateKey []byte `protobuf:"bytes,2,opt,name=enc_private_key,json=encPrivateKey,proto3" json:"enc_private_key,omitempty"`
	KeyIv         []byte `protobuf:"bytes,5,opt,name=key_iv,json=keyIv,proto3" json:"key_iv,omitempty"`
	KeySalt       []byte `protobuf:"bytes,6,opt,name=key_salt,json=keySalt,proto3" json:"key_salt,omitempty"`
}

func (m *UserInfo) Reset()                    { *m = UserInfo{} }
func (m *UserInfo) String() string            { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()               {}
func (*UserInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type SessInfo struct {
	Uid       int64  `protobuf:"varint,1,opt,name=uid" json:"uid,omitempty"`
	SessionId string `protobuf:"bytes,2,opt,name=session_id,json=sessionId" json:"session_id,omitempty"`
}

func (m *SessInfo) Reset()                    { *m = SessInfo{} }
func (m *SessInfo) String() string            { return proto.CompactTextString(m) }
func (*SessInfo) ProtoMessage()               {}
func (*SessInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type Notification struct {
	Type       Notification_Type `protobuf:"varint,1,opt,name=type,enum=protodef.Notification_Type" json:"type,omitempty"`
	SenderName string            `protobuf:"bytes,2,opt,name=sender_name,json=senderName" json:"sender_name,omitempty"`
	Body       []byte            `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *Notification) Reset()                    { *m = Notification{} }
func (m *Notification) String() string            { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()               {}
func (*Notification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func init() {
	proto.RegisterType((*Empty)(nil), "protodef.Empty")
	proto.RegisterType((*BoolResp)(nil), "protodef.BoolResp")
	proto.RegisterType((*NewUserReq)(nil), "protodef.NewUserReq")
	proto.RegisterType((*NewUserResp)(nil), "protodef.NewUserResp")
	proto.RegisterType((*LoginReq)(nil), "protodef.LoginReq")
	proto.RegisterType((*Username)(nil), "protodef.Username")
	proto.RegisterType((*AuthInfo)(nil), "protodef.AuthInfo")
	proto.RegisterType((*Challenge)(nil), "protodef.Challenge")
	proto.RegisterType((*LoginResp)(nil), "protodef.LoginResp")
	proto.RegisterType((*LoginAccepted)(nil), "protodef.LoginAccepted")
	proto.RegisterType((*UserInfo)(nil), "protodef.UserInfo")
	proto.RegisterType((*SessInfo)(nil), "protodef.SessInfo")
	proto.RegisterType((*Notification)(nil), "protodef.Notification")
	proto.RegisterEnum("protodef.Notification_Type", Notification_Type_name, Notification_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Uplink service

type UplinkClient interface {
	Exists(ctx context.Context, in *Username, opts ...grpc.CallOption) (*BoolResp, error)
	LoginExchange(ctx context.Context, opts ...grpc.CallOption) (Uplink_LoginExchangeClient, error)
	NewUser(ctx context.Context, in *NewUserReq, opts ...grpc.CallOption) (*NewUserResp, error)
	Notifications(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Uplink_NotificationsClient, error)
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoolResp, error)
}

type uplinkClient struct {
	cc *grpc.ClientConn
}

func NewUplinkClient(cc *grpc.ClientConn) UplinkClient {
	return &uplinkClient{cc}
}

func (c *uplinkClient) Exists(ctx context.Context, in *Username, opts ...grpc.CallOption) (*BoolResp, error) {
	out := new(BoolResp)
	err := grpc.Invoke(ctx, "/protodef.Uplink/Exists", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uplinkClient) LoginExchange(ctx context.Context, opts ...grpc.CallOption) (Uplink_LoginExchangeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Uplink_serviceDesc.Streams[0], c.cc, "/protodef.Uplink/LoginExchange", opts...)
	if err != nil {
		return nil, err
	}
	x := &uplinkLoginExchangeClient{stream}
	return x, nil
}

type Uplink_LoginExchangeClient interface {
	Send(*LoginReq) error
	Recv() (*LoginResp, error)
	grpc.ClientStream
}

type uplinkLoginExchangeClient struct {
	grpc.ClientStream
}

func (x *uplinkLoginExchangeClient) Send(m *LoginReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *uplinkLoginExchangeClient) Recv() (*LoginResp, error) {
	m := new(LoginResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *uplinkClient) NewUser(ctx context.Context, in *NewUserReq, opts ...grpc.CallOption) (*NewUserResp, error) {
	out := new(NewUserResp)
	err := grpc.Invoke(ctx, "/protodef.Uplink/NewUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uplinkClient) Notifications(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Uplink_NotificationsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Uplink_serviceDesc.Streams[1], c.cc, "/protodef.Uplink/Notifications", opts...)
	if err != nil {
		return nil, err
	}
	x := &uplinkNotificationsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Uplink_NotificationsClient interface {
	Recv() (*Notification, error)
	grpc.ClientStream
}

type uplinkNotificationsClient struct {
	grpc.ClientStream
}

func (x *uplinkNotificationsClient) Recv() (*Notification, error) {
	m := new(Notification)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *uplinkClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoolResp, error) {
	out := new(BoolResp)
	err := grpc.Invoke(ctx, "/protodef.Uplink/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Uplink service

type UplinkServer interface {
	Exists(context.Context, *Username) (*BoolResp, error)
	LoginExchange(Uplink_LoginExchangeServer) error
	NewUser(context.Context, *NewUserReq) (*NewUserResp, error)
	Notifications(*Empty, Uplink_NotificationsServer) error
	Ping(context.Context, *Empty) (*BoolResp, error)
}

func RegisterUplinkServer(s *grpc.Server, srv UplinkServer) {
	s.RegisterService(&_Uplink_serviceDesc, srv)
}

func _Uplink_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Username)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UplinkServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protodef.Uplink/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UplinkServer).Exists(ctx, req.(*Username))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uplink_LoginExchange_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UplinkServer).LoginExchange(&uplinkLoginExchangeServer{stream})
}

type Uplink_LoginExchangeServer interface {
	Send(*LoginResp) error
	Recv() (*LoginReq, error)
	grpc.ServerStream
}

type uplinkLoginExchangeServer struct {
	grpc.ServerStream
}

func (x *uplinkLoginExchangeServer) Send(m *LoginResp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *uplinkLoginExchangeServer) Recv() (*LoginReq, error) {
	m := new(LoginReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Uplink_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UplinkServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protodef.Uplink/NewUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UplinkServer).NewUser(ctx, req.(*NewUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uplink_Notifications_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UplinkServer).Notifications(m, &uplinkNotificationsServer{stream})
}

type Uplink_NotificationsServer interface {
	Send(*Notification) error
	grpc.ServerStream
}

type uplinkNotificationsServer struct {
	grpc.ServerStream
}

func (x *uplinkNotificationsServer) Send(m *Notification) error {
	return x.ServerStream.SendMsg(m)
}

func _Uplink_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UplinkServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protodef.Uplink/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UplinkServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Uplink_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protodef.Uplink",
	HandlerType: (*UplinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exists",
			Handler:    _Uplink_Exists_Handler,
		},
		{
			MethodName: "NewUser",
			Handler:    _Uplink_NewUser_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Uplink_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LoginExchange",
			Handler:       _Uplink_LoginExchange_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Notifications",
			Handler:       _Uplink_Notifications_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("uplink.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 682 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0xd3, 0x7c, 0xd8, 0xe3, 0x84, 0x46, 0xdb, 0x16, 0x42, 0x10, 0x5f, 0x16, 0x42, 0xa5,
	0x45, 0xa6, 0x35, 0x82, 0x03, 0x70, 0x69, 0xc0, 0x82, 0x08, 0x88, 0x8a, 0xd3, 0x72, 0x8d, 0x12,
	0x7b, 0x9b, 0xae, 0xea, 0xd8, 0x26, 0x6b, 0xb7, 0xcd, 0x99, 0x5f, 0xc1, 0x95, 0x3b, 0xff, 0x91,
	0xd9, 0xb5, 0x9d, 0x8f, 0x36, 0x45, 0x3d, 0x70, 0xca, 0xce, 0xdb, 0x37, 0x99, 0x9d, 0x37, 0xcf,
	0x03, 0xd5, 0x24, 0xf2, 0x59, 0x70, 0x6a, 0x46, 0xe3, 0x30, 0x0e, 0x89, 0x2a, 0x7f, 0x3c, 0x7a,
	0x6c, 0x54, 0xa0, 0x64, 0x8f, 0xa2, 0x78, 0x62, 0x3c, 0x01, 0xb5, 0x15, 0x86, 0xbe, 0x43, 0x79,
	0x44, 0x1a, 0x50, 0xe1, 0x89, 0xeb, 0x52, 0xce, 0x1b, 0xca, 0x23, 0x65, 0x4b, 0x75, 0xf2, 0xd0,
	0xf8, 0xa3, 0x00, 0x74, 0xe8, 0xf9, 0x11, 0xa7, 0x63, 0x87, 0xfe, 0x20, 0x04, 0x8a, 0x41, 0x7f,
	0x44, 0x25, 0x4b, 0x73, 0xe4, 0x59, 0x60, 0x51, 0x1f, 0x33, 0x0b, 0x29, 0x26, 0xce, 0xe4, 0x3e,
	0x40, 0x94, 0x0c, 0x7c, 0xe6, 0xf6, 0x4e, 0xe9, 0xa4, 0xb1, 0x8a, 0x37, 0x55, 0x47, 0x4b, 0x91,
	0xcf, 0x74, 0x42, 0x9e, 0xc2, 0x1a, 0x0d, 0xdc, 0x5e, 0x34, 0x66, 0x67, 0xfd, 0x98, 0x4a, 0x4e,
	0x51, 0x72, 0x6a, 0x08, 0x1f, 0xa4, 0xa8, 0xe0, 0x6d, 0x42, 0x19, 0xef, 0x7a, 0xec, 0xac, 0x51,
	0x92, 0xd7, 0x25, 0x8c, 0xda, 0x67, 0xe4, 0x2e, 0xa8, 0x02, 0xe6, 0x7d, 0x3f, 0x6e, 0x94, 0xe5,
	0x45, 0x05, 0xe3, 0x2e, 0x86, 0xc6, 0x07, 0xd0, 0xa7, 0xcf, 0xc5, 0xc6, 0x5e, 0x41, 0x95, 0x63,
	0x1b, 0x2c, 0x0c, 0x7a, 0x2c, 0x38, 0x0e, 0xe5, 0xbb, 0x75, 0x8b, 0x98, 0xb9, 0x1c, 0x66, 0x17,
	0x6f, 0xdb, 0x78, 0xe3, 0xe8, 0x19, 0x4f, 0x04, 0xc6, 0x18, 0xd4, 0x2f, 0xe1, 0x90, 0x05, 0xa2,
	0xe5, 0x6d, 0x28, 0xf1, 0x98, 0x46, 0x7b, 0x57, 0x73, 0xf7, 0x93, 0xf8, 0x44, 0xd0, 0x3f, 0xad,
	0x38, 0x29, 0x85, 0xec, 0xa4, 0x5c, 0x4b, 0x6a, 0xa1, 0x5b, 0xeb, 0x33, 0xee, 0xfb, 0x93, 0xbe,
	0xef, 0xd3, 0x60, 0x48, 0x73, 0xb2, 0xd5, 0xaa, 0x81, 0xee, 0x8b, 0x22, 0x3d, 0x11, 0x72, 0xe3,
	0x01, 0xa8, 0xe2, 0xd9, 0xb9, 0xa4, 0x97, 0x65, 0x36, 0x2c, 0x50, 0xf3, 0x82, 0x37, 0x1d, 0x83,
	0xf1, 0x18, 0xb4, 0x69, 0x61, 0xb2, 0x01, 0xa5, 0x38, 0x3c, 0xa5, 0x81, 0xcc, 0x42, 0x2d, 0x65,
	0x60, 0x9c, 0x83, 0x96, 0xb5, 0x8a, 0x72, 0xbd, 0x58, 0xec, 0xf5, 0xce, 0xec, 0xfd, 0x92, 0xb3,
	0x8f, 0x9e, 0x88, 0x62, 0xea, 0xcd, 0x1a, 0xde, 0x5e, 0x6c, 0x78, 0x89, 0xb0, 0xd7, 0xf6, 0xcb,
	0xa1, 0xb6, 0xf0, 0xa7, 0x58, 0x5c, 0x4b, 0x50, 0x80, 0x6b, 0x06, 0x25, 0xb4, 0x91, 0x83, 0x52,
	0x93, 0xec, 0x44, 0xf6, 0x40, 0x73, 0xf3, 0xee, 0xfe, 0xa1, 0xb8, 0x33, 0x63, 0x19, 0x3f, 0x95,
	0x54, 0x65, 0x99, 0xbf, 0x68, 0x52, 0xe5, 0x06, 0x26, 0x2d, 0xfc, 0x1f, 0x93, 0xbe, 0x05, 0x35,
	0x97, 0x87, 0xd4, 0x61, 0x35, 0x61, 0x9e, 0xac, 0xbe, 0xea, 0x88, 0xa3, 0x78, 0xd6, 0xd4, 0xb3,
	0x5e, 0x36, 0x4e, 0x2d, 0x77, 0xa7, 0x67, 0xfc, 0x56, 0xa0, 0xda, 0x09, 0x63, 0x76, 0xcc, 0xdc,
	0x7e, 0x8c, 0x10, 0xea, 0x56, 0x8c, 0x27, 0x51, 0x6a, 0x86, 0x5b, 0xd6, 0xbd, 0x99, 0x02, 0xf3,
	0x2c, 0xf3, 0x10, 0x29, 0x8e, 0x24, 0x92, 0x87, 0x80, 0x66, 0x0f, 0x3c, 0x94, 0x5a, 0x9a, 0x28,
	0xad, 0x00, 0x29, 0xd4, 0xc9, 0xac, 0x34, 0x08, 0xbd, 0xfc, 0xbb, 0x95, 0x67, 0xe3, 0x39, 0x14,
	0xc5, 0x5f, 0x10, 0x1d, 0x2a, 0x5f, 0xed, 0x6e, 0x77, 0xff, 0xa3, 0x5d, 0x5f, 0x11, 0x81, 0x63,
	0x7f, 0x3b, 0xb2, 0xbb, 0x87, 0x75, 0x85, 0x00, 0x94, 0xdb, 0x9d, 0xef, 0xed, 0x43, 0xbb, 0x5e,
	0xb0, 0x7e, 0x15, 0xa0, 0x7c, 0x24, 0x17, 0x10, 0xd9, 0x85, 0xb2, 0x7d, 0xc1, 0x78, 0xcc, 0xc9,
	0xa5, 0x69, 0x8a, 0xe2, 0xcd, 0x39, 0x6c, 0xba, 0x8d, 0xde, 0x65, 0xce, 0xb0, 0x2f, 0x70, 0x70,
	0xc2, 0xb9, 0xe4, 0x92, 0x0f, 0xf1, 0xb3, 0x6c, 0xae, 0x5f, 0xc1, 0x78, 0xb4, 0xa5, 0xec, 0x2a,
	0xe4, 0x35, 0x54, 0xb2, 0x0d, 0x40, 0x36, 0xe6, 0xb4, 0x98, 0xee, 0xb0, 0xe6, 0xe6, 0x12, 0x14,
	0xab, 0xbe, 0x81, 0xda, 0xbc, 0x60, 0x9c, 0xac, 0xcd, 0x78, 0x72, 0x63, 0x36, 0x6f, 0x2f, 0x97,
	0x16, 0x6b, 0xee, 0x40, 0xf1, 0x80, 0x05, 0xc3, 0xab, 0x29, 0x4b, 0xda, 0x6b, 0x3d, 0x83, 0xa6,
	0x1b, 0x8e, 0xcc, 0x21, 0x8b, 0x4f, 0x92, 0x81, 0x39, 0x72, 0x99, 0xef, 0x87, 0x01, 0x33, 0xd3,
	0x7d, 0xdd, 0xd2, 0x53, 0xd9, 0x0e, 0x44, 0xda, 0xa0, 0x2c, 0xb3, 0x5f, 0xfe, 0x0d, 0x00, 0x00,
	0xff, 0xff, 0x68, 0x43, 0xd2, 0x14, 0xcd, 0x05, 0x00, 0x00,
}
