// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: lib/gobert/src/keywords.proto

package ai

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

type KeywordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Texts []string `protobuf:"bytes,1,rep,name=texts,proto3" json:"texts,omitempty"`
}

func (x *KeywordRequest) Reset() {
	*x = KeywordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_keywords_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeywordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeywordRequest) ProtoMessage() {}

func (x *KeywordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_keywords_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeywordRequest.ProtoReflect.Descriptor instead.
func (*KeywordRequest) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_keywords_proto_rawDescGZIP(), []int{0}
}

func (x *KeywordRequest) GetTexts() []string {
	if x != nil {
		return x.Texts
	}
	return nil
}

type Keyword struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text  []byte  `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *Keyword) Reset() {
	*x = Keyword{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_keywords_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keyword) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keyword) ProtoMessage() {}

func (x *Keyword) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_keywords_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keyword.ProtoReflect.Descriptor instead.
func (*Keyword) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_keywords_proto_rawDescGZIP(), []int{1}
}

func (x *Keyword) GetText() []byte {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *Keyword) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type Keywords struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keywords []*Keyword `protobuf:"bytes,1,rep,name=keywords,proto3" json:"keywords,omitempty"`
}

func (x *Keywords) Reset() {
	*x = Keywords{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_keywords_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keywords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keywords) ProtoMessage() {}

func (x *Keywords) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_keywords_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keywords.ProtoReflect.Descriptor instead.
func (*Keywords) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_keywords_proto_rawDescGZIP(), []int{2}
}

func (x *Keywords) GetKeywords() []*Keyword {
	if x != nil {
		return x.Keywords
	}
	return nil
}

type KeywordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Texts []*Keywords `protobuf:"bytes,1,rep,name=texts,proto3" json:"texts,omitempty"`
}

func (x *KeywordResponse) Reset() {
	*x = KeywordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_keywords_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeywordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeywordResponse) ProtoMessage() {}

func (x *KeywordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_keywords_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeywordResponse.ProtoReflect.Descriptor instead.
func (*KeywordResponse) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_keywords_proto_rawDescGZIP(), []int{3}
}

func (x *KeywordResponse) GetTexts() []*Keywords {
	if x != nil {
		return x.Texts
	}
	return nil
}

var File_lib_gobert_src_keywords_proto protoreflect.FileDescriptor

var file_lib_gobert_src_keywords_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6c, 0x69, 0x62, 0x2f, 0x67, 0x6f, 0x62, 0x65, 0x72, 0x74, 0x2f, 0x73, 0x72, 0x63,
	0x2f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x61, 0x69, 0x2e, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x26, 0x0a, 0x0e,
	0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x65, 0x78, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x65, 0x78, 0x74, 0x73, 0x22, 0x33, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x3c, 0x0a, 0x08, 0x4b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x30, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x69, 0x2e, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x08, 0x6b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x3e, 0x0a, 0x0f, 0x4b, 0x65, 0x79, 0x77, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x74, 0x65,
	0x78, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x69, 0x2e, 0x6b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x05, 0x74, 0x65, 0x78, 0x74, 0x73, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lib_gobert_src_keywords_proto_rawDescOnce sync.Once
	file_lib_gobert_src_keywords_proto_rawDescData = file_lib_gobert_src_keywords_proto_rawDesc
)

func file_lib_gobert_src_keywords_proto_rawDescGZIP() []byte {
	file_lib_gobert_src_keywords_proto_rawDescOnce.Do(func() {
		file_lib_gobert_src_keywords_proto_rawDescData = protoimpl.X.CompressGZIP(file_lib_gobert_src_keywords_proto_rawDescData)
	})
	return file_lib_gobert_src_keywords_proto_rawDescData
}

var file_lib_gobert_src_keywords_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_lib_gobert_src_keywords_proto_goTypes = []interface{}{
	(*KeywordRequest)(nil),  // 0: ai.keywords.KeywordRequest
	(*Keyword)(nil),         // 1: ai.keywords.Keyword
	(*Keywords)(nil),        // 2: ai.keywords.Keywords
	(*KeywordResponse)(nil), // 3: ai.keywords.KeywordResponse
}
var file_lib_gobert_src_keywords_proto_depIdxs = []int32{
	1, // 0: ai.keywords.Keywords.keywords:type_name -> ai.keywords.Keyword
	2, // 1: ai.keywords.KeywordResponse.texts:type_name -> ai.keywords.Keywords
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_lib_gobert_src_keywords_proto_init() }
func file_lib_gobert_src_keywords_proto_init() {
	if File_lib_gobert_src_keywords_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lib_gobert_src_keywords_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeywordRequest); i {
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
		file_lib_gobert_src_keywords_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Keyword); i {
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
		file_lib_gobert_src_keywords_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Keywords); i {
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
		file_lib_gobert_src_keywords_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeywordResponse); i {
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
			RawDescriptor: file_lib_gobert_src_keywords_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lib_gobert_src_keywords_proto_goTypes,
		DependencyIndexes: file_lib_gobert_src_keywords_proto_depIdxs,
		MessageInfos:      file_lib_gobert_src_keywords_proto_msgTypes,
	}.Build()
	File_lib_gobert_src_keywords_proto = out.File
	file_lib_gobert_src_keywords_proto_rawDesc = nil
	file_lib_gobert_src_keywords_proto_goTypes = nil
	file_lib_gobert_src_keywords_proto_depIdxs = nil
}
