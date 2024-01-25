// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: lib/goadblock/src/adblock.proto

package parsing

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

type FilterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Classes []string `protobuf:"bytes,1,rep,name=classes,proto3" json:"classes,omitempty"`
	Ids     []string `protobuf:"bytes,2,rep,name=ids,proto3" json:"ids,omitempty"`
	Urls    []string `protobuf:"bytes,3,rep,name=urls,proto3" json:"urls,omitempty"`
	BaseUrl string   `protobuf:"bytes,4,opt,name=base_url,json=baseUrl,proto3" json:"base_url,omitempty"`
}

func (x *FilterRequest) Reset() {
	*x = FilterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_goadblock_src_adblock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterRequest) ProtoMessage() {}

func (x *FilterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lib_goadblock_src_adblock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterRequest.ProtoReflect.Descriptor instead.
func (*FilterRequest) Descriptor() ([]byte, []int) {
	return file_lib_goadblock_src_adblock_proto_rawDescGZIP(), []int{0}
}

func (x *FilterRequest) GetClasses() []string {
	if x != nil {
		return x.Classes
	}
	return nil
}

func (x *FilterRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *FilterRequest) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

func (x *FilterRequest) GetBaseUrl() string {
	if x != nil {
		return x.BaseUrl
	}
	return ""
}

type FilterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Matches        []string `protobuf:"bytes,1,rep,name=matches,proto3" json:"matches,omitempty"`
	BlockedDomains []string `protobuf:"bytes,2,rep,name=blocked_domains,json=blockedDomains,proto3" json:"blocked_domains,omitempty"`
}

func (x *FilterResponse) Reset() {
	*x = FilterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_goadblock_src_adblock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterResponse) ProtoMessage() {}

func (x *FilterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lib_goadblock_src_adblock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterResponse.ProtoReflect.Descriptor instead.
func (*FilterResponse) Descriptor() ([]byte, []int) {
	return file_lib_goadblock_src_adblock_proto_rawDescGZIP(), []int{1}
}

func (x *FilterResponse) GetMatches() []string {
	if x != nil {
		return x.Matches
	}
	return nil
}

func (x *FilterResponse) GetBlockedDomains() []string {
	if x != nil {
		return x.BlockedDomains
	}
	return nil
}

type Rules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []string `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *Rules) Reset() {
	*x = Rules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_goadblock_src_adblock_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rules) ProtoMessage() {}

func (x *Rules) ProtoReflect() protoreflect.Message {
	mi := &file_lib_goadblock_src_adblock_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rules.ProtoReflect.Descriptor instead.
func (*Rules) Descriptor() ([]byte, []int) {
	return file_lib_goadblock_src_adblock_proto_rawDescGZIP(), []int{2}
}

func (x *Rules) GetRules() []string {
	if x != nil {
		return x.Rules
	}
	return nil
}

type RuleGroups struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filters []*Rules `protobuf:"bytes,1,rep,name=filters,proto3" json:"filters,omitempty"`
}

func (x *RuleGroups) Reset() {
	*x = RuleGroups{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_goadblock_src_adblock_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuleGroups) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuleGroups) ProtoMessage() {}

func (x *RuleGroups) ProtoReflect() protoreflect.Message {
	mi := &file_lib_goadblock_src_adblock_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuleGroups.ProtoReflect.Descriptor instead.
func (*RuleGroups) Descriptor() ([]byte, []int) {
	return file_lib_goadblock_src_adblock_proto_rawDescGZIP(), []int{3}
}

func (x *RuleGroups) GetFilters() []*Rules {
	if x != nil {
		return x.Filters
	}
	return nil
}

var File_lib_goadblock_src_adblock_proto protoreflect.FileDescriptor

var file_lib_goadblock_src_adblock_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6c, 0x69, 0x62, 0x2f, 0x67, 0x6f, 0x61, 0x64, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2f,
	0x73, 0x72, 0x63, 0x2f, 0x61, 0x64, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0f, 0x61, 0x64, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x22, 0x6a, 0x0a, 0x0d, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x12, 0x10, 0x0a,
	0x03, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x72, 0x6c, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x73, 0x65, 0x55, 0x72, 0x6c, 0x22, 0x53,
	0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x65, 0x64, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x73, 0x22, 0x1d, 0x0a, 0x05, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x72, 0x75, 0x6c,
	0x65, 0x73, 0x22, 0x3e, 0x0a, 0x0a, 0x52, 0x75, 0x6c, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x12, 0x30, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x61, 0x64, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x73, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x61, 0x72, 0x73,
	0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lib_goadblock_src_adblock_proto_rawDescOnce sync.Once
	file_lib_goadblock_src_adblock_proto_rawDescData = file_lib_goadblock_src_adblock_proto_rawDesc
)

func file_lib_goadblock_src_adblock_proto_rawDescGZIP() []byte {
	file_lib_goadblock_src_adblock_proto_rawDescOnce.Do(func() {
		file_lib_goadblock_src_adblock_proto_rawDescData = protoimpl.X.CompressGZIP(file_lib_goadblock_src_adblock_proto_rawDescData)
	})
	return file_lib_goadblock_src_adblock_proto_rawDescData
}

var file_lib_goadblock_src_adblock_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_lib_goadblock_src_adblock_proto_goTypes = []interface{}{
	(*FilterRequest)(nil),  // 0: adblock.content.FilterRequest
	(*FilterResponse)(nil), // 1: adblock.content.FilterResponse
	(*Rules)(nil),          // 2: adblock.content.Rules
	(*RuleGroups)(nil),     // 3: adblock.content.RuleGroups
}
var file_lib_goadblock_src_adblock_proto_depIdxs = []int32{
	2, // 0: adblock.content.RuleGroups.filters:type_name -> adblock.content.Rules
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_lib_goadblock_src_adblock_proto_init() }
func file_lib_goadblock_src_adblock_proto_init() {
	if File_lib_goadblock_src_adblock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lib_goadblock_src_adblock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilterRequest); i {
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
		file_lib_goadblock_src_adblock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilterResponse); i {
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
		file_lib_goadblock_src_adblock_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rules); i {
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
		file_lib_goadblock_src_adblock_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuleGroups); i {
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
			RawDescriptor: file_lib_goadblock_src_adblock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lib_goadblock_src_adblock_proto_goTypes,
		DependencyIndexes: file_lib_goadblock_src_adblock_proto_depIdxs,
		MessageInfos:      file_lib_goadblock_src_adblock_proto_msgTypes,
	}.Build()
	File_lib_goadblock_src_adblock_proto = out.File
	file_lib_goadblock_src_adblock_proto_rawDesc = nil
	file_lib_goadblock_src_adblock_proto_goTypes = nil
	file_lib_goadblock_src_adblock_proto_depIdxs = nil
}
