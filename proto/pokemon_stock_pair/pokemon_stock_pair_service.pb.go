// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: proto/pokemon_stock_pair/pokemon_stock_pair_service.proto

package pokemon_stock_pair

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto protoreflect.FileDescriptor

var file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_rawDesc = []byte{
	0x0a, 0x39, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f,
	0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d,
	0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x6f, 0x6b,
	0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x1a,
	0x42, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73,
	0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x61, 0x6c,
	0x6c, 0x5f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f,
	0x70, 0x61, 0x69, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x43, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d,
	0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x2f, 0x67, 0x65,
	0x74, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x73, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x3d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61,
	0x69, 0x72, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73,
	0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x3e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69,
	0x72, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x96, 0x02, 0x0a, 0x17, 0x50, 0x6f, 0x6b, 0x65,
	0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x50, 0x6f,
	0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x73, 0x12,
	0x32, 0x2e, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f,
	0x70, 0x61, 0x69, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x50, 0x6f, 0x6b, 0x65, 0x6d,
	0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x50,
	0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x76, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x50,
	0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x12,
	0x2e, 0x2e, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f,
	0x70, 0x61, 0x69, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2f, 0x2e, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f,
	0x70, 0x61, 0x69, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x25, 0x5a, 0x23, 0x70, 0x6f, 0x6b, 0x65, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f,
	0x63, 0x6b, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_goTypes = []any{
	(*GetAllPokemonStockPairsRequest)(nil),  // 0: pokemon_stock_pair.GetAllPokemonStockPairsRequest
	(*GetPokemonStockPairRequest)(nil),      // 1: pokemon_stock_pair.GetPokemonStockPairRequest
	(*GetAllPokemonStockPairsResponse)(nil), // 2: pokemon_stock_pair.GetAllPokemonStockPairsResponse
	(*GetPokemonStockPairResponse)(nil),     // 3: pokemon_stock_pair.GetPokemonStockPairResponse
}
var file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_depIdxs = []int32{
	0, // 0: pokemon_stock_pair.PokemonStockPairService.GetAllPokemonStockPairs:input_type -> pokemon_stock_pair.GetAllPokemonStockPairsRequest
	1, // 1: pokemon_stock_pair.PokemonStockPairService.GetPokemonStockPair:input_type -> pokemon_stock_pair.GetPokemonStockPairRequest
	2, // 2: pokemon_stock_pair.PokemonStockPairService.GetAllPokemonStockPairs:output_type -> pokemon_stock_pair.GetAllPokemonStockPairsResponse
	3, // 3: pokemon_stock_pair.PokemonStockPairService.GetPokemonStockPair:output_type -> pokemon_stock_pair.GetPokemonStockPairResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_init() }
func file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_init() {
	if File_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto != nil {
		return
	}
	file_proto_pokemon_stock_pair_get_all_pokemon_stock_pairs_request_proto_init()
	file_proto_pokemon_stock_pair_get_all_pokemon_stock_pairs_response_proto_init()
	file_proto_pokemon_stock_pair_get_pokemon_stock_pair_request_proto_init()
	file_proto_pokemon_stock_pair_get_pokemon_stock_pair_response_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_goTypes,
		DependencyIndexes: file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_depIdxs,
	}.Build()
	File_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto = out.File
	file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_rawDesc = nil
	file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_goTypes = nil
	file_proto_pokemon_stock_pair_pokemon_stock_pair_service_proto_depIdxs = nil
}
