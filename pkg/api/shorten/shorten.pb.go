// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: api/shorten.proto

package shorten

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

type CreateShortURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // id пользователя
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`                     // исходный url
}

func (x *CreateShortURLRequest) Reset() {
	*x = CreateShortURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShortURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLRequest) ProtoMessage() {}

func (x *CreateShortURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLRequest.ProtoReflect.Descriptor instead.
func (*CreateShortURLRequest) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{0}
}

func (x *CreateShortURLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CreateShortURLRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type CreateShortURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` // короткий url
}

func (x *CreateShortURLResponse) Reset() {
	*x = CreateShortURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShortURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLResponse) ProtoMessage() {}

func (x *CreateShortURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLResponse.ProtoReflect.Descriptor instead.
func (*CreateShortURLResponse) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{1}
}

func (x *CreateShortURLResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type BatchCreateShortURLRequestData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalUrl   string `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`       // исходный url
	CorrelationId string `protobuf:"bytes,2,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"` // уникальный идентификатор
}

func (x *BatchCreateShortURLRequestData) Reset() {
	*x = BatchCreateShortURLRequestData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCreateShortURLRequestData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCreateShortURLRequestData) ProtoMessage() {}

func (x *BatchCreateShortURLRequestData) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCreateShortURLRequestData.ProtoReflect.Descriptor instead.
func (*BatchCreateShortURLRequestData) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{2}
}

func (x *BatchCreateShortURLRequestData) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *BatchCreateShortURLRequestData) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

type BatchCreateShortURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items  []*BatchCreateShortURLRequestData `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	UserId string                            `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // id пользователя
}

func (x *BatchCreateShortURLRequest) Reset() {
	*x = BatchCreateShortURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCreateShortURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCreateShortURLRequest) ProtoMessage() {}

func (x *BatchCreateShortURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCreateShortURLRequest.ProtoReflect.Descriptor instead.
func (*BatchCreateShortURLRequest) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{3}
}

func (x *BatchCreateShortURLRequest) GetItems() []*BatchCreateShortURLRequestData {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *BatchCreateShortURLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type BatchCreateShortURLResponseData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl      string `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`                // короткий url
	CorrelationId string `protobuf:"bytes,2,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"` // уникальный идентификатор
}

func (x *BatchCreateShortURLResponseData) Reset() {
	*x = BatchCreateShortURLResponseData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCreateShortURLResponseData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCreateShortURLResponseData) ProtoMessage() {}

func (x *BatchCreateShortURLResponseData) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCreateShortURLResponseData.ProtoReflect.Descriptor instead.
func (*BatchCreateShortURLResponseData) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{4}
}

func (x *BatchCreateShortURLResponseData) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *BatchCreateShortURLResponseData) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

type BatchCreateShortURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*BatchCreateShortURLResponseData `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *BatchCreateShortURLResponse) Reset() {
	*x = BatchCreateShortURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCreateShortURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCreateShortURLResponse) ProtoMessage() {}

func (x *BatchCreateShortURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCreateShortURLResponse.ProtoReflect.Descriptor instead.
func (*BatchCreateShortURLResponse) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{5}
}

func (x *BatchCreateShortURLResponse) GetItems() []*BatchCreateShortURLResponseData {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetByShortRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // id пользователя
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`                     // короткие url
}

func (x *GetByShortRequest) Reset() {
	*x = GetByShortRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByShortRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByShortRequest) ProtoMessage() {}

func (x *GetByShortRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByShortRequest.ProtoReflect.Descriptor instead.
func (*GetByShortRequest) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{6}
}

func (x *GetByShortRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetByShortRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetByShortResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FullUrl string `protobuf:"bytes,1,opt,name=full_url,json=fullUrl,proto3" json:"full_url,omitempty"` // исходный url
}

func (x *GetByShortResponse) Reset() {
	*x = GetByShortResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByShortResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByShortResponse) ProtoMessage() {}

func (x *GetByShortResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByShortResponse.ProtoReflect.Descriptor instead.
func (*GetByShortResponse) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{7}
}

func (x *GetByShortResponse) GetFullUrl() string {
	if x != nil {
		return x.FullUrl
	}
	return ""
}

type GetUserURLsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // id пользователя
}

func (x *GetUserURLsRequest) Reset() {
	*x = GetUserURLsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserURLsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserURLsRequest) ProtoMessage() {}

func (x *GetUserURLsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserURLsRequest.ProtoReflect.Descriptor instead.
func (*GetUserURLsRequest) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{8}
}

func (x *GetUserURLsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ShortenData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl    string `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`          // короткий url
	OriginalUrl string `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"` // исходный url
}

func (x *ShortenData) Reset() {
	*x = ShortenData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenData) ProtoMessage() {}

func (x *ShortenData) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenData.ProtoReflect.Descriptor instead.
func (*ShortenData) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{9}
}

func (x *ShortenData) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *ShortenData) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type GetUserURLsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ShortenData `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetUserURLsResponse) Reset() {
	*x = GetUserURLsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserURLsResponse) ProtoMessage() {}

func (x *GetUserURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserURLsResponse.ProtoReflect.Descriptor instead.
func (*GetUserURLsResponse) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{10}
}

func (x *GetUserURLsResponse) GetItems() []*ShortenData {
	if x != nil {
		return x.Items
	}
	return nil
}

type DeleteUserURLsBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // id пользователя
	Urls   []string `protobuf:"bytes,2,rep,name=urls,proto3" json:"urls,omitempty"`                   // короткие url
}

func (x *DeleteUserURLsBatchRequest) Reset() {
	*x = DeleteUserURLsBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserURLsBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserURLsBatchRequest) ProtoMessage() {}

func (x *DeleteUserURLsBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserURLsBatchRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserURLsBatchRequest) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteUserURLsBatchRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeleteUserURLsBatchRequest) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

type DeleteUserURLsBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteUserURLsBatchResponse) Reset() {
	*x = DeleteUserURLsBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_shorten_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserURLsBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserURLsBatchResponse) ProtoMessage() {}

func (x *DeleteUserURLsBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_shorten_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserURLsBatchResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserURLsBatchResponse) Descriptor() ([]byte, []int) {
	return file_api_shorten_proto_rawDescGZIP(), []int{12}
}

var File_api_shorten_proto protoreflect.FileDescriptor

var file_api_shorten_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x61, 0x72, 0x74, 0x22, 0x42, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x30, 0x0a,
	0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22,
	0x6a, 0x0a, 0x1e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61,
	0x6c, 0x55, 0x72, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f,
	0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x71, 0x0a, 0x1a, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x65,
	0x0a, 0x1f, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x25,
	0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x5a, 0x0a, 0x1b, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x3e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x22, 0x2f, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x42, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x75, 0x6c, 0x6c, 0x55,
	0x72, 0x6c, 0x22, 0x2d, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x4d, 0x0a, 0x0b, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a,
	0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c,
	0x22, 0x3e, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x22, 0x49, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52,
	0x4c, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x1d, 0x0a, 0x1b, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x93, 0x03, 0x0a, 0x07, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x12, 0x4b, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x1b, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x13, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x2e, 0x63, 0x61, 0x72,
	0x74, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63,
	0x61, 0x72, 0x74, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3f, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x17, 0x2e,
	0x63, 0x61, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e, 0x47, 0x65,
	0x74, 0x42, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x42, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x12,
	0x18, 0x2e, 0x63, 0x61, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52,
	0x4c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x61, 0x72, 0x74,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x20, 0x2e, 0x63, 0x61,
	0x72, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c,
	0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x63, 0x61, 0x72, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x55,
	0x52, 0x4c, 0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x11, 0x5a, 0x0f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_shorten_proto_rawDescOnce sync.Once
	file_api_shorten_proto_rawDescData = file_api_shorten_proto_rawDesc
)

func file_api_shorten_proto_rawDescGZIP() []byte {
	file_api_shorten_proto_rawDescOnce.Do(func() {
		file_api_shorten_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_shorten_proto_rawDescData)
	})
	return file_api_shorten_proto_rawDescData
}

var file_api_shorten_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_api_shorten_proto_goTypes = []interface{}{
	(*CreateShortURLRequest)(nil),           // 0: cart.CreateShortURLRequest
	(*CreateShortURLResponse)(nil),          // 1: cart.CreateShortURLResponse
	(*BatchCreateShortURLRequestData)(nil),  // 2: cart.BatchCreateShortURLRequestData
	(*BatchCreateShortURLRequest)(nil),      // 3: cart.BatchCreateShortURLRequest
	(*BatchCreateShortURLResponseData)(nil), // 4: cart.BatchCreateShortURLResponseData
	(*BatchCreateShortURLResponse)(nil),     // 5: cart.BatchCreateShortURLResponse
	(*GetByShortRequest)(nil),               // 6: cart.GetByShortRequest
	(*GetByShortResponse)(nil),              // 7: cart.GetByShortResponse
	(*GetUserURLsRequest)(nil),              // 8: cart.GetUserURLsRequest
	(*ShortenData)(nil),                     // 9: cart.ShortenData
	(*GetUserURLsResponse)(nil),             // 10: cart.GetUserURLsResponse
	(*DeleteUserURLsBatchRequest)(nil),      // 11: cart.DeleteUserURLsBatchRequest
	(*DeleteUserURLsBatchResponse)(nil),     // 12: cart.DeleteUserURLsBatchResponse
}
var file_api_shorten_proto_depIdxs = []int32{
	2,  // 0: cart.BatchCreateShortURLRequest.items:type_name -> cart.BatchCreateShortURLRequestData
	4,  // 1: cart.BatchCreateShortURLResponse.items:type_name -> cart.BatchCreateShortURLResponseData
	9,  // 2: cart.GetUserURLsResponse.items:type_name -> cart.ShortenData
	0,  // 3: cart.Shorten.CreateShortURL:input_type -> cart.CreateShortURLRequest
	3,  // 4: cart.Shorten.BatchCreateShortURL:input_type -> cart.BatchCreateShortURLRequest
	6,  // 5: cart.Shorten.GetByShort:input_type -> cart.GetByShortRequest
	8,  // 6: cart.Shorten.GetUserURLs:input_type -> cart.GetUserURLsRequest
	11, // 7: cart.Shorten.DeleteUserURLsBatch:input_type -> cart.DeleteUserURLsBatchRequest
	1,  // 8: cart.Shorten.CreateShortURL:output_type -> cart.CreateShortURLResponse
	5,  // 9: cart.Shorten.BatchCreateShortURL:output_type -> cart.BatchCreateShortURLResponse
	7,  // 10: cart.Shorten.GetByShort:output_type -> cart.GetByShortResponse
	10, // 11: cart.Shorten.GetUserURLs:output_type -> cart.GetUserURLsResponse
	12, // 12: cart.Shorten.DeleteUserURLsBatch:output_type -> cart.DeleteUserURLsBatchResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_api_shorten_proto_init() }
func file_api_shorten_proto_init() {
	if File_api_shorten_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_shorten_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateShortURLRequest); i {
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
		file_api_shorten_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateShortURLResponse); i {
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
		file_api_shorten_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCreateShortURLRequestData); i {
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
		file_api_shorten_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCreateShortURLRequest); i {
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
		file_api_shorten_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCreateShortURLResponseData); i {
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
		file_api_shorten_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCreateShortURLResponse); i {
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
		file_api_shorten_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByShortRequest); i {
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
		file_api_shorten_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByShortResponse); i {
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
		file_api_shorten_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserURLsRequest); i {
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
		file_api_shorten_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenData); i {
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
		file_api_shorten_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserURLsResponse); i {
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
		file_api_shorten_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserURLsBatchRequest); i {
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
		file_api_shorten_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserURLsBatchResponse); i {
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
			RawDescriptor: file_api_shorten_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_shorten_proto_goTypes,
		DependencyIndexes: file_api_shorten_proto_depIdxs,
		MessageInfos:      file_api_shorten_proto_msgTypes,
	}.Build()
	File_api_shorten_proto = out.File
	file_api_shorten_proto_rawDesc = nil
	file_api_shorten_proto_goTypes = nil
	file_api_shorten_proto_depIdxs = nil
}