// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: mirror/v1/server.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Instance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Host        string `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	HostName    string `protobuf:"bytes,3,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	Port        int32  `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	ServiceName string `protobuf:"bytes,5,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Meta        string `protobuf:"bytes,6,opt,name=meta,proto3" json:"meta,omitempty"`
}

func (x *Instance) Reset() {
	*x = Instance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_v1_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instance) ProtoMessage() {}

func (x *Instance) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_v1_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instance.ProtoReflect.Descriptor instead.
func (*Instance) Descriptor() ([]byte, []int) {
	return file_mirror_v1_server_proto_rawDescGZIP(), []int{0}
}

func (x *Instance) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Instance) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Instance) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *Instance) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Instance) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *Instance) GetMeta() string {
	if x != nil {
		return x.Meta
	}
	return ""
}

type RegisterReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authentication string    `protobuf:"bytes,1,opt,name=authentication,proto3" json:"authentication,omitempty"`
	Instance       *Instance `protobuf:"bytes,2,opt,name=instance,proto3" json:"instance,omitempty"`
	Filter         []string  `protobuf:"bytes,3,rep,name=filter,proto3" json:"filter,omitempty"`
}

func (x *RegisterReq) Reset() {
	*x = RegisterReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_v1_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterReq) ProtoMessage() {}

func (x *RegisterReq) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_v1_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterReq.ProtoReflect.Descriptor instead.
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return file_mirror_v1_server_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterReq) GetAuthentication() string {
	if x != nil {
		return x.Authentication
	}
	return ""
}

func (x *RegisterReq) GetInstance() *Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

func (x *RegisterReq) GetFilter() []string {
	if x != nil {
		return x.Filter
	}
	return nil
}

type RegisterRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool  `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Ttl     int32 `protobuf:"varint,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *RegisterRes) Reset() {
	*x = RegisterRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_v1_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRes) ProtoMessage() {}

func (x *RegisterRes) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_v1_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRes.ProtoReflect.Descriptor instead.
func (*RegisterRes) Descriptor() ([]byte, []int) {
	return file_mirror_v1_server_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterRes) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RegisterRes) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

type UnRegisterReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authentication string    `protobuf:"bytes,1,opt,name=authentication,proto3" json:"authentication,omitempty"`
	Instance       *Instance `protobuf:"bytes,2,opt,name=instance,proto3" json:"instance,omitempty"`
}

func (x *UnRegisterReq) Reset() {
	*x = UnRegisterReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_v1_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnRegisterReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnRegisterReq) ProtoMessage() {}

func (x *UnRegisterReq) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_v1_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnRegisterReq.ProtoReflect.Descriptor instead.
func (*UnRegisterReq) Descriptor() ([]byte, []int) {
	return file_mirror_v1_server_proto_rawDescGZIP(), []int{3}
}

func (x *UnRegisterReq) GetAuthentication() string {
	if x != nil {
		return x.Authentication
	}
	return ""
}

func (x *UnRegisterReq) GetInstance() *Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

type UnRegisterRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UnRegisterRes) Reset() {
	*x = UnRegisterRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mirror_v1_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnRegisterRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnRegisterRes) ProtoMessage() {}

func (x *UnRegisterRes) ProtoReflect() protoreflect.Message {
	mi := &file_mirror_v1_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnRegisterRes.ProtoReflect.Descriptor instead.
func (*UnRegisterRes) Descriptor() ([]byte, []int) {
	return file_mirror_v1_server_proto_rawDescGZIP(), []int{4}
}

func (x *UnRegisterRes) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_mirror_v1_server_proto protoreflect.FileDescriptor

var file_mirror_v1_server_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x96, 0x01, 0x0a, 0x08, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x22, 0x7b, 0x0a, 0x0b, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x26, 0x0a, 0x0e, 0x61, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2c, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x39, 0x0a, 0x0b, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x74, 0x74,
	0x6c, 0x22, 0x65, 0x0a, 0x0d, 0x55, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x12, 0x26, 0x0a, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x08, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d,
	0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x08,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x29, 0x0a, 0x0d, 0x55, 0x6e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x32, 0x7e, 0x0a, 0x06, 0x4d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x36, 0x0a,
	0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x6d, 0x69, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x13,
	0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0a, 0x55, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x55, 0x6e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x6d, 0x69, 0x72,
	0x72, 0x6f, 0x72, 0x2e, 0x55, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x22, 0x00, 0x42, 0x13, 0x5a, 0x11, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x69, 0x72, 0x72,
	0x6f, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mirror_v1_server_proto_rawDescOnce sync.Once
	file_mirror_v1_server_proto_rawDescData = file_mirror_v1_server_proto_rawDesc
)

func file_mirror_v1_server_proto_rawDescGZIP() []byte {
	file_mirror_v1_server_proto_rawDescOnce.Do(func() {
		file_mirror_v1_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_mirror_v1_server_proto_rawDescData)
	})
	return file_mirror_v1_server_proto_rawDescData
}

var file_mirror_v1_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_mirror_v1_server_proto_goTypes = []any{
	(*Instance)(nil),      // 0: mirror.Instance
	(*RegisterReq)(nil),   // 1: mirror.RegisterReq
	(*RegisterRes)(nil),   // 2: mirror.RegisterRes
	(*UnRegisterReq)(nil), // 3: mirror.UnRegisterReq
	(*UnRegisterRes)(nil), // 4: mirror.UnRegisterRes
}
var file_mirror_v1_server_proto_depIdxs = []int32{
	0, // 0: mirror.RegisterReq.instance:type_name -> mirror.Instance
	0, // 1: mirror.UnRegisterReq.instance:type_name -> mirror.Instance
	1, // 2: mirror.Mirror.Register:input_type -> mirror.RegisterReq
	3, // 3: mirror.Mirror.UnRegister:input_type -> mirror.UnRegisterReq
	2, // 4: mirror.Mirror.Register:output_type -> mirror.RegisterRes
	4, // 5: mirror.Mirror.UnRegister:output_type -> mirror.UnRegisterRes
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_mirror_v1_server_proto_init() }
func file_mirror_v1_server_proto_init() {
	if File_mirror_v1_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mirror_v1_server_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Instance); i {
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
		file_mirror_v1_server_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterReq); i {
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
		file_mirror_v1_server_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterRes); i {
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
		file_mirror_v1_server_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UnRegisterReq); i {
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
		file_mirror_v1_server_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*UnRegisterRes); i {
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
			RawDescriptor: file_mirror_v1_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mirror_v1_server_proto_goTypes,
		DependencyIndexes: file_mirror_v1_server_proto_depIdxs,
		MessageInfos:      file_mirror_v1_server_proto_msgTypes,
	}.Build()
	File_mirror_v1_server_proto = out.File
	file_mirror_v1_server_proto_rawDesc = nil
	file_mirror_v1_server_proto_goTypes = nil
	file_mirror_v1_server_proto_depIdxs = nil
}
