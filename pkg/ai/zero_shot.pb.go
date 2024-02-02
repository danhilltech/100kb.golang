// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: lib/gobert/src/zero_shot.proto

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

type ZeroShotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Texts  []string `protobuf:"bytes,1,rep,name=texts,proto3" json:"texts,omitempty"`
	Labels []string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty"`
}

func (x *ZeroShotRequest) Reset() {
	*x = ZeroShotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZeroShotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZeroShotRequest) ProtoMessage() {}

func (x *ZeroShotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZeroShotRequest.ProtoReflect.Descriptor instead.
func (*ZeroShotRequest) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_zero_shot_proto_rawDescGZIP(), []int{0}
}

func (x *ZeroShotRequest) GetTexts() []string {
	if x != nil {
		return x.Texts
	}
	return nil
}

func (x *ZeroShotRequest) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ZeroShotClassification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Label []byte  `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *ZeroShotClassification) Reset() {
	*x = ZeroShotClassification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZeroShotClassification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZeroShotClassification) ProtoMessage() {}

func (x *ZeroShotClassification) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZeroShotClassification.ProtoReflect.Descriptor instead.
func (*ZeroShotClassification) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_zero_shot_proto_rawDescGZIP(), []int{1}
}

func (x *ZeroShotClassification) GetLabel() []byte {
	if x != nil {
		return x.Label
	}
	return nil
}

func (x *ZeroShotClassification) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type ZeroShotClassifications struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Classifications []*ZeroShotClassification `protobuf:"bytes,1,rep,name=classifications,proto3" json:"classifications,omitempty"`
}

func (x *ZeroShotClassifications) Reset() {
	*x = ZeroShotClassifications{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZeroShotClassifications) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZeroShotClassifications) ProtoMessage() {}

func (x *ZeroShotClassifications) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZeroShotClassifications.ProtoReflect.Descriptor instead.
func (*ZeroShotClassifications) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_zero_shot_proto_rawDescGZIP(), []int{2}
}

func (x *ZeroShotClassifications) GetClassifications() []*ZeroShotClassification {
	if x != nil {
		return x.Classifications
	}
	return nil
}

type ZeroShotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sentences []*ZeroShotClassifications `protobuf:"bytes,1,rep,name=sentences,proto3" json:"sentences,omitempty"`
}

func (x *ZeroShotResponse) Reset() {
	*x = ZeroShotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZeroShotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZeroShotResponse) ProtoMessage() {}

func (x *ZeroShotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lib_gobert_src_zero_shot_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZeroShotResponse.ProtoReflect.Descriptor instead.
func (*ZeroShotResponse) Descriptor() ([]byte, []int) {
	return file_lib_gobert_src_zero_shot_proto_rawDescGZIP(), []int{3}
}

func (x *ZeroShotResponse) GetSentences() []*ZeroShotClassifications {
	if x != nil {
		return x.Sentences
	}
	return nil
}

var File_lib_gobert_src_zero_shot_proto protoreflect.FileDescriptor

var file_lib_gobert_src_zero_shot_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x6c, 0x69, 0x62, 0x2f, 0x67, 0x6f, 0x62, 0x65, 0x72, 0x74, 0x2f, 0x73, 0x72, 0x63,
	0x2f, 0x7a, 0x65, 0x72, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x61, 0x69, 0x2e, 0x7a, 0x65, 0x72, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x74, 0x22, 0x3f,
	0x0a, 0x0f, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x65, 0x78, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x65, 0x78, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x22,
	0x44, 0x0a, 0x16, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f, 0x74, 0x43, 0x6c, 0x61, 0x73, 0x73,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x69, 0x0a, 0x17, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f,
	0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x4e, 0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61, 0x69, 0x2e, 0x7a,
	0x65, 0x72, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f,
	0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x22, 0x57, 0x0a, 0x10, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x69, 0x2e, 0x7a, 0x65, 0x72,
	0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x5a, 0x65, 0x72, 0x6f, 0x53, 0x68, 0x6f, 0x74, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x09,
	0x73, 0x65, 0x6e, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lib_gobert_src_zero_shot_proto_rawDescOnce sync.Once
	file_lib_gobert_src_zero_shot_proto_rawDescData = file_lib_gobert_src_zero_shot_proto_rawDesc
)

func file_lib_gobert_src_zero_shot_proto_rawDescGZIP() []byte {
	file_lib_gobert_src_zero_shot_proto_rawDescOnce.Do(func() {
		file_lib_gobert_src_zero_shot_proto_rawDescData = protoimpl.X.CompressGZIP(file_lib_gobert_src_zero_shot_proto_rawDescData)
	})
	return file_lib_gobert_src_zero_shot_proto_rawDescData
}

var file_lib_gobert_src_zero_shot_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_lib_gobert_src_zero_shot_proto_goTypes = []interface{}{
	(*ZeroShotRequest)(nil),         // 0: ai.zero_shot.ZeroShotRequest
	(*ZeroShotClassification)(nil),  // 1: ai.zero_shot.ZeroShotClassification
	(*ZeroShotClassifications)(nil), // 2: ai.zero_shot.ZeroShotClassifications
	(*ZeroShotResponse)(nil),        // 3: ai.zero_shot.ZeroShotResponse
}
var file_lib_gobert_src_zero_shot_proto_depIdxs = []int32{
	1, // 0: ai.zero_shot.ZeroShotClassifications.classifications:type_name -> ai.zero_shot.ZeroShotClassification
	2, // 1: ai.zero_shot.ZeroShotResponse.sentences:type_name -> ai.zero_shot.ZeroShotClassifications
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_lib_gobert_src_zero_shot_proto_init() }
func file_lib_gobert_src_zero_shot_proto_init() {
	if File_lib_gobert_src_zero_shot_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lib_gobert_src_zero_shot_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZeroShotRequest); i {
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
		file_lib_gobert_src_zero_shot_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZeroShotClassification); i {
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
		file_lib_gobert_src_zero_shot_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZeroShotClassifications); i {
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
		file_lib_gobert_src_zero_shot_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZeroShotResponse); i {
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
			RawDescriptor: file_lib_gobert_src_zero_shot_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lib_gobert_src_zero_shot_proto_goTypes,
		DependencyIndexes: file_lib_gobert_src_zero_shot_proto_depIdxs,
		MessageInfos:      file_lib_gobert_src_zero_shot_proto_msgTypes,
	}.Build()
	File_lib_gobert_src_zero_shot_proto = out.File
	file_lib_gobert_src_zero_shot_proto_rawDesc = nil
	file_lib_gobert_src_zero_shot_proto_goTypes = nil
	file_lib_gobert_src_zero_shot_proto_depIdxs = nil
}
