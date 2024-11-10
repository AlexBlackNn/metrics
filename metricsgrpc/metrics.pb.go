// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: metricsgrpc/metrics.proto

package metricsgrpc_v1

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

// MetricRequest is request for updating or getting a metric.
type MetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`         // Metric ID.
	Type  string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`     // Metric type (gauge or counter).
	Delta int64   `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`  // Delta value for counter metrics.
	Value float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"` // Value for gauge metrics.
}

func (x *MetricRequest) Reset() {
	*x = MetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsgrpc_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricRequest) ProtoMessage() {}

func (x *MetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metricsgrpc_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricRequest.ProtoReflect.Descriptor instead.
func (*MetricRequest) Descriptor() ([]byte, []int) {
	return file_metricsgrpc_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *MetricRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MetricRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MetricRequest) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *MetricRequest) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// MetricsRequest is request for updating several metrics.
type MetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*MetricRequest `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"` // List of metrics to update.
}

func (x *MetricsRequest) Reset() {
	*x = MetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsgrpc_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsRequest) ProtoMessage() {}

func (x *MetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metricsgrpc_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsRequest.ProtoReflect.Descriptor instead.
func (*MetricsRequest) Descriptor() ([]byte, []int) {
	return file_metricsgrpc_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsRequest) GetMetrics() []*MetricRequest {
	if x != nil {
		return x.Metrics
	}
	return nil
}

// MetricResponse is response for getting a metric.
type MetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`         // Metric ID.
	Type  string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`     // Metric type (gauge or counter).
	Delta int64   `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`  // Delta value for counter metrics.
	Value float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"` // Value for gauge metrics.
}

func (x *MetricResponse) Reset() {
	*x = MetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsgrpc_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricResponse) ProtoMessage() {}

func (x *MetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metricsgrpc_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricResponse.ProtoReflect.Descriptor instead.
func (*MetricResponse) Descriptor() ([]byte, []int) {
	return file_metricsgrpc_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *MetricResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MetricResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MetricResponse) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *MetricResponse) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Response is response for updating metrics.
type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"` // Status of the operation.
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`   // Error message (if any).
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsgrpc_metrics_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_metricsgrpc_metrics_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_metricsgrpc_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_metricsgrpc_metrics_proto protoreflect.FileDescriptor

var file_metricsgrpc_metrics_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x22, 0x5f, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x46, 0x0a, 0x0e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x22, 0x60, 0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x38, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xe8, 0x01,
	0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x41, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65, 0x76,
	0x65, 0x72, 0x61, 0x6c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1b, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x47, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_metricsgrpc_metrics_proto_rawDescOnce sync.Once
	file_metricsgrpc_metrics_proto_rawDescData = file_metricsgrpc_metrics_proto_rawDesc
)

func file_metricsgrpc_metrics_proto_rawDescGZIP() []byte {
	file_metricsgrpc_metrics_proto_rawDescOnce.Do(func() {
		file_metricsgrpc_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_metricsgrpc_metrics_proto_rawDescData)
	})
	return file_metricsgrpc_metrics_proto_rawDescData
}

var file_metricsgrpc_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_metricsgrpc_metrics_proto_goTypes = []any{
	(*MetricRequest)(nil),  // 0: metricsgrpc.MetricRequest
	(*MetricsRequest)(nil), // 1: metricsgrpc.MetricsRequest
	(*MetricResponse)(nil), // 2: metricsgrpc.MetricResponse
	(*Response)(nil),       // 3: metricsgrpc.Response
}
var file_metricsgrpc_metrics_proto_depIdxs = []int32{
	0, // 0: metricsgrpc.MetricsRequest.metrics:type_name -> metricsgrpc.MetricRequest
	0, // 1: metricsgrpc.MetricsService.UpdateMetric:input_type -> metricsgrpc.MetricRequest
	1, // 2: metricsgrpc.MetricsService.UpdateSeveralMetrics:input_type -> metricsgrpc.MetricsRequest
	0, // 3: metricsgrpc.MetricsService.GetOneMetric:input_type -> metricsgrpc.MetricRequest
	3, // 4: metricsgrpc.MetricsService.UpdateMetric:output_type -> metricsgrpc.Response
	3, // 5: metricsgrpc.MetricsService.UpdateSeveralMetrics:output_type -> metricsgrpc.Response
	2, // 6: metricsgrpc.MetricsService.GetOneMetric:output_type -> metricsgrpc.MetricResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_metricsgrpc_metrics_proto_init() }
func file_metricsgrpc_metrics_proto_init() {
	if File_metricsgrpc_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metricsgrpc_metrics_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MetricRequest); i {
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
		file_metricsgrpc_metrics_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MetricsRequest); i {
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
		file_metricsgrpc_metrics_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*MetricResponse); i {
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
		file_metricsgrpc_metrics_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_metricsgrpc_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metricsgrpc_metrics_proto_goTypes,
		DependencyIndexes: file_metricsgrpc_metrics_proto_depIdxs,
		MessageInfos:      file_metricsgrpc_metrics_proto_msgTypes,
	}.Build()
	File_metricsgrpc_metrics_proto = out.File
	file_metricsgrpc_metrics_proto_rawDesc = nil
	file_metricsgrpc_metrics_proto_goTypes = nil
	file_metricsgrpc_metrics_proto_depIdxs = nil
}