package mapo

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

// ClientState from MAPO tracks the current validator set, latest height,
// and a possible frozen height.
type ClientState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frozen           bool        `protobuf:"varint,1,opt,name=frozen,proto3" json:"frozen,omitempty"`
	LatestEpoch      uint64      `protobuf:"varint,2,opt,name=latestEpoch,proto3" json:"latestEpoch,omitempty"`
	EpochSize        uint64      `protobuf:"varint,3,opt,name=epochSize,proto3" json:"epochSize,omitempty"`
	ClientIdentifier *Identifier `protobuf:"bytes,4,opt,name=ClientIdentifier,proto3" json:"ClientIdentifier,omitempty"`
}

func (x *ClientState) Reset() {
	*x = ClientState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientState) ProtoMessage() {}

func (x *ClientState) ProtoReflect() protoreflect.Message {
	mi := &file_mapo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientState.ProtoReflect.Descriptor instead.
func (*ClientState) Descriptor() ([]byte, []int) {
	return file_mapo_proto_rawDescGZIP(), []int{0}
}

func (x *ClientState) GetFrozen() bool {
	if x != nil {
		return x.Frozen
	}
	return false
}

func (x *ClientState) GetLatestEpoch() uint64 {
	if x != nil {
		return x.LatestEpoch
	}
	return 0
}

func (x *ClientState) GetEpochSize() uint64 {
	if x != nil {
		return x.EpochSize
	}
	return 0
}

func (x *ClientState) GetClientIdentifier() *Identifier {
	if x != nil {
		return x.ClientIdentifier
	}
	return nil
}

type Identifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Identifier) Reset() {
	*x = Identifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Identifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Identifier) ProtoMessage() {}

func (x *Identifier) ProtoReflect() protoreflect.Message {
	mi := &file_mapo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Identifier.ProtoReflect.Descriptor instead.
func (*Identifier) Descriptor() ([]byte, []int) {
	return file_mapo_proto_rawDescGZIP(), []int{1}
}

var File_mapo_proto protoreflect.FileDescriptor

var file_mapo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x70, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x69, 0x62,
	0x63, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x6d,
	0x61, 0x70, 0x6f, 0x2e, 0x76, 0x31, 0x22, 0xb7, 0x01, 0x0a, 0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x72, 0x6f, 0x7a, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x66, 0x72, 0x6f, 0x7a, 0x65, 0x6e, 0x12, 0x20,
	0x0a, 0x0b, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x45, 0x70, 0x6f, 0x63, 0x68,
	0x12, 0x1c, 0x0a, 0x09, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x50,
	0x0a, 0x10, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x62, 0x63, 0x2e, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x6d, 0x61, 0x70, 0x6f,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x10,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x22, 0x0c, 0x0a, 0x0a, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x42, 0x40,
	0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73,
	0x6d, 0x6f, 0x73, 0x2f, 0x69, 0x62, 0x63, 0x2d, 0x67, 0x6f, 0x2f, 0x76, 0x37, 0x2f, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x2d, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x36, 0x31, 0x2d, 0x6d, 0x61, 0x70, 0x6f, 0x3b, 0x6d, 0x61, 0x70, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mapo_proto_rawDescOnce sync.Once
	file_mapo_proto_rawDescData = file_mapo_proto_rawDesc
)

func file_mapo_proto_rawDescGZIP() []byte {
	file_mapo_proto_rawDescOnce.Do(func() {
		file_mapo_proto_rawDescData = protoimpl.X.CompressGZIP(file_mapo_proto_rawDescData)
	})
	return file_mapo_proto_rawDescData
}

var file_mapo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mapo_proto_goTypes = []interface{}{
	(*ClientState)(nil), // 0: ibc.lightclients.mapo.v1.ClientState
	(*Identifier)(nil),  // 1: ibc.lightclients.mapo.v1.Identifier
}
var file_mapo_proto_depIdxs = []int32{
	1, // 0: ibc.lightclients.mapo.v1.ClientState.ClientIdentifier:type_name -> ibc.lightclients.mapo.v1.Identifier
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_mapo_proto_init() }
func file_mapo_proto_init() {
	if File_mapo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mapo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientState); i {
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
		file_mapo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Identifier); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mapo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mapo_proto_goTypes,
		DependencyIndexes: file_mapo_proto_depIdxs,
		MessageInfos:      file_mapo_proto_msgTypes,
	}.Build()
	File_mapo_proto = out.File
	file_mapo_proto_rawDesc = nil
	file_mapo_proto_goTypes = nil
	file_mapo_proto_depIdxs = nil
}
