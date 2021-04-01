// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.6
// source: mirror.proto

package mirror

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_mirror_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pixels []byte `protobuf:"bytes,1,opt,name=Pixels,proto3" json:"Pixels,omitempty"`
	Stride int32  `protobuf:"varint,2,opt,name=Stride,proto3" json:"Stride,omitempty"`
	Width  int32  `protobuf:"varint,3,opt,name=Width,proto3" json:"Width,omitempty"`
	Height int32  `protobuf:"varint,4,opt,name=Height,proto3" json:"Height,omitempty"`
}

func (x *Image) Reset() {
	*x = Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Image) ProtoMessage() {}

func (x *Image) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Image.ProtoReflect.Descriptor instead.
func (*Image) Descriptor() ([]byte, []int) {
	return file_mirror_proto_rawDescGZIP(), []int{1}
}

func (x *Image) GetPixels() []byte {
	if x != nil {
		return x.Pixels
	}
	return nil
}

func (x *Image) GetStride() int32 {
	if x != nil {
		return x.Stride
	}
	return 0
}

func (x *Image) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *Image) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

var File_mirror_proto protoreflect.FileDescriptor

var file_mirror_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x26, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x65,
	0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x69, 0x78, 0x65, 0x6c,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x69, 0x78, 0x65, 0x6c, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x69, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x53, 0x74, 0x72, 0x69, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x57, 0x69, 0x64, 0x74, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x57, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a,
	0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x32, 0x42, 0x0a, 0x06, 0x4d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x38, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x18, 0x2e, 0x6d,
	0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06, 0x6d, 0x69, 0x72,
	0x72, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mirror_proto_rawDescOnce sync.Once
	file_mirror_proto_rawDescData = file_mirror_proto_rawDesc
)

func file_mirror_proto_rawDescGZIP() []byte {
	file_mirror_proto_rawDescOnce.Do(func() {
		file_mirror_proto_rawDescData = protoimpl.X.CompressGZIP(file_mirror_proto_rawDescData)
	})
	return file_mirror_proto_rawDescData
}

var file_mirror_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mirror_proto_goTypes = []interface{}{
	(*SubscribeRequest)(nil), // 0: mirror.SubscribeRequest
	(*Image)(nil),            // 1: mirror.Image
}
var file_mirror_proto_depIdxs = []int32{
	0, // 0: mirror.Mirror.Subscribe:input_type -> mirror.SubscribeRequest
	1, // 1: mirror.Mirror.Subscribe:output_type -> mirror.Image
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mirror_proto_init() }
func file_mirror_proto_init() {
	if File_mirror_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mirror_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeRequest); i {
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
		file_mirror_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Image); i {
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
			RawDescriptor: file_mirror_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mirror_proto_goTypes,
		DependencyIndexes: file_mirror_proto_depIdxs,
		MessageInfos:      file_mirror_proto_msgTypes,
	}.Build()
	File_mirror_proto = out.File
	file_mirror_proto_rawDesc = nil
	file_mirror_proto_goTypes = nil
	file_mirror_proto_depIdxs = nil
}
