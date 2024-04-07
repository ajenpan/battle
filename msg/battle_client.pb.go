// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: battle_client.proto

package msg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReqJoinBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReadyState int32 `protobuf:"varint,1,opt,name=ready_state,json=readyState,proto3" json:"ready_state,omitempty"`
}

func (x *ReqJoinBattle) Reset() {
	*x = ReqJoinBattle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqJoinBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqJoinBattle) ProtoMessage() {}

func (x *ReqJoinBattle) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqJoinBattle.ProtoReflect.Descriptor instead.
func (*ReqJoinBattle) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{0}
}

func (x *ReqJoinBattle) GetReadyState() int32 {
	if x != nil {
		return x.ReadyState
	}
	return 0
}

type RespJoinBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RespJoinBattle) Reset() {
	*x = RespJoinBattle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespJoinBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespJoinBattle) ProtoMessage() {}

func (x *RespJoinBattle) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespJoinBattle.ProtoReflect.Descriptor instead.
func (*RespJoinBattle) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{1}
}

type ReqQuitBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReqQuitBattle) Reset() {
	*x = ReqQuitBattle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqQuitBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqQuitBattle) ProtoMessage() {}

func (x *ReqQuitBattle) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqQuitBattle.ProtoReflect.Descriptor instead.
func (*ReqQuitBattle) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{2}
}

type RespQuitBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RespQuitBattle) Reset() {
	*x = RespQuitBattle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RespQuitBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespQuitBattle) ProtoMessage() {}

func (x *RespQuitBattle) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespQuitBattle.ProtoReflect.Descriptor instead.
func (*RespQuitBattle) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{3}
}

type MsgToClient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgid int32  `protobuf:"varint,1,opt,name=msgid,proto3" json:"msgid,omitempty"`
	Body  []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *MsgToClient) Reset() {
	*x = MsgToClient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgToClient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgToClient) ProtoMessage() {}

func (x *MsgToClient) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgToClient.ProtoReflect.Descriptor instead.
func (*MsgToClient) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{4}
}

func (x *MsgToClient) GetMsgid() int32 {
	if x != nil {
		return x.Msgid
	}
	return 0
}

func (x *MsgToClient) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type MsgToLogic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgid int32  `protobuf:"varint,1,opt,name=msgid,proto3" json:"msgid,omitempty"`
	Body  []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *MsgToLogic) Reset() {
	*x = MsgToLogic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgToLogic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgToLogic) ProtoMessage() {}

func (x *MsgToLogic) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgToLogic.ProtoReflect.Descriptor instead.
func (*MsgToLogic) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{5}
}

func (x *MsgToLogic) GetMsgid() int32 {
	if x != nil {
		return x.Msgid
	}
	return 0
}

func (x *MsgToLogic) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type BattleMessageWrap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Battleid uint64 `protobuf:"varint,1,opt,name=battleid,proto3" json:"battleid,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*BattleMessageWrap_ReqJoin
	//	*BattleMessageWrap_RespJoin
	//	*BattleMessageWrap_ReqQuit
	//	*BattleMessageWrap_RespQuit
	//	*BattleMessageWrap_ToClient
	//	*BattleMessageWrap_ToLogic
	Payload isBattleMessageWrap_Payload `protobuf_oneof:"payload"`
}

func (x *BattleMessageWrap) Reset() {
	*x = BattleMessageWrap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_battle_client_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BattleMessageWrap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BattleMessageWrap) ProtoMessage() {}

func (x *BattleMessageWrap) ProtoReflect() protoreflect.Message {
	mi := &file_battle_client_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BattleMessageWrap.ProtoReflect.Descriptor instead.
func (*BattleMessageWrap) Descriptor() ([]byte, []int) {
	return file_battle_client_proto_rawDescGZIP(), []int{6}
}

func (x *BattleMessageWrap) GetBattleid() uint64 {
	if x != nil {
		return x.Battleid
	}
	return 0
}

func (m *BattleMessageWrap) GetPayload() isBattleMessageWrap_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *BattleMessageWrap) GetReqJoin() *ReqJoinBattle {
	if x, ok := x.GetPayload().(*BattleMessageWrap_ReqJoin); ok {
		return x.ReqJoin
	}
	return nil
}

func (x *BattleMessageWrap) GetRespJoin() *RespJoinBattle {
	if x, ok := x.GetPayload().(*BattleMessageWrap_RespJoin); ok {
		return x.RespJoin
	}
	return nil
}

func (x *BattleMessageWrap) GetReqQuit() *ReqQuitBattle {
	if x, ok := x.GetPayload().(*BattleMessageWrap_ReqQuit); ok {
		return x.ReqQuit
	}
	return nil
}

func (x *BattleMessageWrap) GetRespQuit() *RespQuitBattle {
	if x, ok := x.GetPayload().(*BattleMessageWrap_RespQuit); ok {
		return x.RespQuit
	}
	return nil
}

func (x *BattleMessageWrap) GetToClient() *MsgToClient {
	if x, ok := x.GetPayload().(*BattleMessageWrap_ToClient); ok {
		return x.ToClient
	}
	return nil
}

func (x *BattleMessageWrap) GetToLogic() *MsgToLogic {
	if x, ok := x.GetPayload().(*BattleMessageWrap_ToLogic); ok {
		return x.ToLogic
	}
	return nil
}

type isBattleMessageWrap_Payload interface {
	isBattleMessageWrap_Payload()
}

type BattleMessageWrap_ReqJoin struct {
	ReqJoin *ReqJoinBattle `protobuf:"bytes,2,opt,name=req_join,json=reqJoin,proto3,oneof"`
}

type BattleMessageWrap_RespJoin struct {
	RespJoin *RespJoinBattle `protobuf:"bytes,3,opt,name=resp_join,json=respJoin,proto3,oneof"`
}

type BattleMessageWrap_ReqQuit struct {
	ReqQuit *ReqQuitBattle `protobuf:"bytes,4,opt,name=req_quit,json=reqQuit,proto3,oneof"`
}

type BattleMessageWrap_RespQuit struct {
	RespQuit *RespQuitBattle `protobuf:"bytes,5,opt,name=resp_quit,json=respQuit,proto3,oneof"`
}

type BattleMessageWrap_ToClient struct {
	ToClient *MsgToClient `protobuf:"bytes,6,opt,name=to_client,json=toClient,proto3,oneof"`
}

type BattleMessageWrap_ToLogic struct {
	ToLogic *MsgToLogic `protobuf:"bytes,7,opt,name=to_logic,json=toLogic,proto3,oneof"`
}

func (*BattleMessageWrap_ReqJoin) isBattleMessageWrap_Payload() {}

func (*BattleMessageWrap_RespJoin) isBattleMessageWrap_Payload() {}

func (*BattleMessageWrap_ReqQuit) isBattleMessageWrap_Payload() {}

func (*BattleMessageWrap_RespQuit) isBattleMessageWrap_Payload() {}

func (*BattleMessageWrap_ToClient) isBattleMessageWrap_Payload() {}

func (*BattleMessageWrap_ToLogic) isBattleMessageWrap_Payload() {}

var File_battle_client_proto protoreflect.FileDescriptor

var file_battle_client_proto_rawDesc = []byte{
	0x0a, 0x13, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x22, 0x30, 0x0a,
	0x0d, 0x52, 0x65, 0x71, 0x4a, 0x6f, 0x69, 0x6e, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x61, 0x64, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x10, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x42, 0x61, 0x74, 0x74, 0x6c,
	0x65, 0x22, 0x0f, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x51, 0x75, 0x69, 0x74, 0x42, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x51, 0x75, 0x69, 0x74, 0x42, 0x61,
	0x74, 0x74, 0x6c, 0x65, 0x22, 0x37, 0x0a, 0x0b, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x73, 0x67, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x36, 0x0a,
	0x0a, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x73, 0x67, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0xf5, 0x02, 0x0a, 0x11, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x62,
	0x61, 0x74, 0x74, 0x6c, 0x65, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62,
	0x61, 0x74, 0x74, 0x6c, 0x65, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x5f, 0x6a,
	0x6f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x4a, 0x6f, 0x69, 0x6e, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x71, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x35, 0x0a, 0x09, 0x72,
	0x65, 0x73, 0x70, 0x5f, 0x6a, 0x6f, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x4a, 0x6f, 0x69, 0x6e,
	0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x4a, 0x6f,
	0x69, 0x6e, 0x12, 0x32, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x5f, 0x71, 0x75, 0x69, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2e, 0x52, 0x65,
	0x71, 0x51, 0x75, 0x69, 0x74, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x07, 0x72,
	0x65, 0x71, 0x51, 0x75, 0x69, 0x74, 0x12, 0x35, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x70, 0x5f, 0x71,
	0x75, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x51, 0x75, 0x69, 0x74, 0x42, 0x61, 0x74, 0x74, 0x6c,
	0x65, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x51, 0x75, 0x69, 0x74, 0x12, 0x32, 0x0a,
	0x09, 0x74, 0x6f, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x08, 0x74, 0x6f, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x12, 0x2f, 0x0a, 0x08, 0x74, 0x6f, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2e, 0x4d, 0x73, 0x67,
	0x54, 0x6f, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x48, 0x00, 0x52, 0x07, 0x74, 0x6f, 0x4c, 0x6f, 0x67,
	0x69, 0x63, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x10, 0x5a,
	0x05, 0x2e, 0x3b, 0x6d, 0x73, 0x67, 0xaa, 0x02, 0x06, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_battle_client_proto_rawDescOnce sync.Once
	file_battle_client_proto_rawDescData = file_battle_client_proto_rawDesc
)

func file_battle_client_proto_rawDescGZIP() []byte {
	file_battle_client_proto_rawDescOnce.Do(func() {
		file_battle_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_battle_client_proto_rawDescData)
	})
	return file_battle_client_proto_rawDescData
}

var file_battle_client_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_battle_client_proto_goTypes = []interface{}{
	(*ReqJoinBattle)(nil),     // 0: battle.ReqJoinBattle
	(*RespJoinBattle)(nil),    // 1: battle.RespJoinBattle
	(*ReqQuitBattle)(nil),     // 2: battle.ReqQuitBattle
	(*RespQuitBattle)(nil),    // 3: battle.RespQuitBattle
	(*MsgToClient)(nil),       // 4: battle.MsgToClient
	(*MsgToLogic)(nil),        // 5: battle.MsgToLogic
	(*BattleMessageWrap)(nil), // 6: battle.BattleMessageWrap
}
var file_battle_client_proto_depIdxs = []int32{
	0, // 0: battle.BattleMessageWrap.req_join:type_name -> battle.ReqJoinBattle
	1, // 1: battle.BattleMessageWrap.resp_join:type_name -> battle.RespJoinBattle
	2, // 2: battle.BattleMessageWrap.req_quit:type_name -> battle.ReqQuitBattle
	3, // 3: battle.BattleMessageWrap.resp_quit:type_name -> battle.RespQuitBattle
	4, // 4: battle.BattleMessageWrap.to_client:type_name -> battle.MsgToClient
	5, // 5: battle.BattleMessageWrap.to_logic:type_name -> battle.MsgToLogic
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_battle_client_proto_init() }
func file_battle_client_proto_init() {
	if File_battle_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_battle_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqJoinBattle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespJoinBattle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqQuitBattle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RespQuitBattle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgToClient); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgToLogic); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_battle_client_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BattleMessageWrap); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_battle_client_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*BattleMessageWrap_ReqJoin)(nil),
		(*BattleMessageWrap_RespJoin)(nil),
		(*BattleMessageWrap_ReqQuit)(nil),
		(*BattleMessageWrap_RespQuit)(nil),
		(*BattleMessageWrap_ToClient)(nil),
		(*BattleMessageWrap_ToLogic)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_battle_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_battle_client_proto_goTypes,
		DependencyIndexes: file_battle_client_proto_depIdxs,
		MessageInfos:      file_battle_client_proto_msgTypes,
	}.Build()
	File_battle_client_proto = out.File
	file_battle_client_proto_rawDesc = nil
	file_battle_client_proto_goTypes = nil
	file_battle_client_proto_depIdxs = nil
}
