// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: proto/common/pokemon.proto

package common

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Pokemon struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	PokedexNumber int32                  `protobuf:"varint,3,opt,name=pokedexNumber,proto3" json:"pokedexNumber,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	Type1         *PokemonType           `protobuf:"bytes,6,opt,name=type1,proto3" json:"type1,omitempty"`
	Type2         *PokemonType           `protobuf:"bytes,7,opt,name=type2,proto3" json:"type2,omitempty"`
	SpriteUrl     string                 `protobuf:"bytes,8,opt,name=spriteUrl,proto3" json:"spriteUrl,omitempty"`
}

func (x *Pokemon) Reset() {
	*x = Pokemon{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_common_pokemon_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pokemon) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pokemon) ProtoMessage() {}

func (x *Pokemon) ProtoReflect() protoreflect.Message {
	mi := &file_proto_common_pokemon_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pokemon.ProtoReflect.Descriptor instead.
func (*Pokemon) Descriptor() ([]byte, []int) {
	return file_proto_common_pokemon_proto_rawDescGZIP(), []int{0}
}

func (x *Pokemon) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Pokemon) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Pokemon) GetPokedexNumber() int32 {
	if x != nil {
		return x.PokedexNumber
	}
	return 0
}

func (x *Pokemon) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Pokemon) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Pokemon) GetType1() *PokemonType {
	if x != nil {
		return x.Type1
	}
	return nil
}

func (x *Pokemon) GetType2() *PokemonType {
	if x != nil {
		return x.Type2
	}
	return nil
}

func (x *Pokemon) GetSpriteUrl() string {
	if x != nil {
		return x.SpriteUrl
	}
	return ""
}

var File_proto_common_pokemon_proto protoreflect.FileDescriptor

var file_proto_common_pokemon_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70,
	0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x02, 0x0a, 0x07, 0x50, 0x6f, 0x6b, 0x65, 0x6d,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x6f, 0x6b, 0x65, 0x64, 0x65,
	0x78, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x70,
	0x6f, 0x6b, 0x65, 0x64, 0x65, 0x78, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x38, 0x0a, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x29, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x31, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x31, 0x12, 0x29, 0x0a, 0x05, 0x74,
	0x79, 0x70, 0x65, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x05, 0x74, 0x79, 0x70, 0x65, 0x32, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x70, 0x72, 0x69, 0x74, 0x65,
	0x55, 0x72, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x70, 0x72, 0x69, 0x74,
	0x65, 0x55, 0x72, 0x6c, 0x42, 0x19, 0x5a, 0x17, 0x70, 0x6f, 0x6b, 0x65, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_common_pokemon_proto_rawDescOnce sync.Once
	file_proto_common_pokemon_proto_rawDescData = file_proto_common_pokemon_proto_rawDesc
)

func file_proto_common_pokemon_proto_rawDescGZIP() []byte {
	file_proto_common_pokemon_proto_rawDescOnce.Do(func() {
		file_proto_common_pokemon_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_common_pokemon_proto_rawDescData)
	})
	return file_proto_common_pokemon_proto_rawDescData
}

var file_proto_common_pokemon_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_common_pokemon_proto_goTypes = []any{
	(*Pokemon)(nil),               // 0: common.Pokemon
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*PokemonType)(nil),           // 2: common.PokemonType
}
var file_proto_common_pokemon_proto_depIdxs = []int32{
	1, // 0: common.Pokemon.createdAt:type_name -> google.protobuf.Timestamp
	1, // 1: common.Pokemon.updatedAt:type_name -> google.protobuf.Timestamp
	2, // 2: common.Pokemon.type1:type_name -> common.PokemonType
	2, // 3: common.Pokemon.type2:type_name -> common.PokemonType
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_common_pokemon_proto_init() }
func file_proto_common_pokemon_proto_init() {
	if File_proto_common_pokemon_proto != nil {
		return
	}
	file_proto_common_pokemon_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_common_pokemon_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Pokemon); i {
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
			RawDescriptor: file_proto_common_pokemon_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_common_pokemon_proto_goTypes,
		DependencyIndexes: file_proto_common_pokemon_proto_depIdxs,
		MessageInfos:      file_proto_common_pokemon_proto_msgTypes,
	}.Build()
	File_proto_common_pokemon_proto = out.File
	file_proto_common_pokemon_proto_rawDesc = nil
	file_proto_common_pokemon_proto_goTypes = nil
	file_proto_common_pokemon_proto_depIdxs = nil
}
