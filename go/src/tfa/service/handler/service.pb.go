// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

/*
Package handler is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	User
	SCPayload
*/
package handler

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PayloadType int32

const (
	PayloadType_USER_CREATE PayloadType = 0
	PayloadType_USER_UPDATE PayloadType = 1
)

var PayloadType_name = map[int32]string{
	0: "USER_CREATE",
	1: "USER_UPDATE",
}
var PayloadType_value = map[string]int32{
	"USER_CREATE": 0,
	"USER_UPDATE": 1,
}

func (x PayloadType) String() string {
	return proto.EnumName(PayloadType_name, int32(x))
}
func (PayloadType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type User struct {
	PhoneNumber string  `protobuf:"bytes,1,opt,name=PhoneNumber" json:"PhoneNumber,omitempty"`
	Uin         uint32  `protobuf:"varint,2,opt,name=Uin" json:"Uin,omitempty"`
	Name        string  `protobuf:"bytes,3,opt,name=Name" json:"Name,omitempty"`
	IsVerified  bool    `protobuf:"varint,4,opt,name=IsVerified" json:"IsVerified,omitempty"`
	Email       string  `protobuf:"bytes,5,opt,name=Email" json:"Email,omitempty"`
	Sex         string  `protobuf:"bytes,6,opt,name=Sex" json:"Sex,omitempty"`
	Birthdate   float64 `protobuf:"fixed64,7,opt,name=Birthdate" json:"Birthdate,omitempty"`
	PushToken   string  `protobuf:"bytes,8,opt,name=PushToken" json:"PushToken,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *User) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *User) GetUin() uint32 {
	if m != nil {
		return m.Uin
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetIsVerified() bool {
	if m != nil {
		return m.IsVerified
	}
	return false
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetSex() string {
	if m != nil {
		return m.Sex
	}
	return ""
}

func (m *User) GetBirthdate() float64 {
	if m != nil {
		return m.Birthdate
	}
	return 0
}

func (m *User) GetPushToken() string {
	if m != nil {
		return m.PushToken
	}
	return ""
}

type SCPayload struct {
	Action      PayloadType `protobuf:"varint,1,opt,name=Action,enum=PayloadType" json:"Action,omitempty"`
	PhoneNumber string      `protobuf:"bytes,2,opt,name=PhoneNumber" json:"PhoneNumber,omitempty"`
	PayloadUser *User       `protobuf:"bytes,3,opt,name=PayloadUser" json:"PayloadUser,omitempty"`
}

func (m *SCPayload) Reset()                    { *m = SCPayload{} }
func (m *SCPayload) String() string            { return proto.CompactTextString(m) }
func (*SCPayload) ProtoMessage()               {}
func (*SCPayload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SCPayload) GetAction() PayloadType {
	if m != nil {
		return m.Action
	}
	return PayloadType_USER_CREATE
}

func (m *SCPayload) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *SCPayload) GetPayloadUser() *User {
	if m != nil {
		return m.PayloadUser
	}
	return nil
}

func init() {
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*SCPayload)(nil), "SCPayload")
	proto.RegisterEnum("PayloadType", PayloadType_name, PayloadType_value)
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0xbf, 0xe9, 0xff, 0xdc, 0x7c, 0xd5, 0x32, 0xb8, 0x98, 0x85, 0x48, 0x28, 0x82, 0xc1,
	0x45, 0x85, 0xfa, 0x04, 0x6d, 0xcd, 0xc2, 0x4d, 0x09, 0x93, 0xc4, 0x85, 0x1b, 0x99, 0x36, 0x57,
	0x32, 0x98, 0x64, 0xca, 0x24, 0x15, 0x0b, 0x3e, 0xa5, 0x4f, 0x24, 0x33, 0x86, 0x36, 0xe8, 0xee,
	0xdc, 0xdf, 0xb9, 0x67, 0x86, 0x33, 0x03, 0xe3, 0x0a, 0xf5, 0xbb, 0xdc, 0xe2, 0x6c, 0xa7, 0x55,
	0xad, 0xa6, 0x5f, 0x04, 0x7a, 0x49, 0x85, 0x9a, 0x7a, 0xe0, 0x86, 0x99, 0x2a, 0x71, 0xbd, 0x2f,
	0x36, 0xa8, 0x19, 0xf1, 0x88, 0xef, 0xf0, 0x36, 0xa2, 0x13, 0xe8, 0x26, 0xb2, 0x64, 0x1d, 0x8f,
	0xf8, 0x63, 0x6e, 0x24, 0xa5, 0xd0, 0x5b, 0x8b, 0x02, 0x59, 0xd7, 0x2e, 0x5b, 0x4d, 0xaf, 0x00,
	0x1e, 0xab, 0x27, 0xd4, 0xf2, 0x55, 0x62, 0xca, 0x7a, 0x1e, 0xf1, 0x47, 0xbc, 0x45, 0xe8, 0x05,
	0xf4, 0x83, 0x42, 0xc8, 0x9c, 0xf5, 0x6d, 0xe8, 0x67, 0x30, 0x67, 0x47, 0xf8, 0xc1, 0x06, 0x96,
	0x19, 0x49, 0x2f, 0xc1, 0x59, 0x4a, 0x5d, 0x67, 0xa9, 0xa8, 0x91, 0x0d, 0x3d, 0xe2, 0x13, 0x7e,
	0x02, 0xc6, 0x0d, 0xf7, 0x55, 0x16, 0xab, 0x37, 0x2c, 0xd9, 0xc8, 0xa6, 0x4e, 0x60, 0xfa, 0x09,
	0x4e, 0xb4, 0x0a, 0xc5, 0x21, 0x57, 0x22, 0xa5, 0xd7, 0x30, 0x58, 0x6c, 0x6b, 0xa9, 0x4a, 0xdb,
	0xe9, 0x6c, 0xfe, 0x7f, 0xd6, 0x38, 0xf1, 0x61, 0x87, 0xbc, 0xf1, 0x7e, 0xd7, 0xef, 0xfc, 0xad,
	0x7f, 0x03, 0x6e, 0x13, 0x34, 0xef, 0x65, 0x3b, 0xbb, 0xf3, 0xfe, 0xcc, 0x0c, 0xbc, 0xed, 0xdc,
	0xde, 0x1d, 0x17, 0xcd, 0x0d, 0xf4, 0x1c, 0xdc, 0x24, 0x0a, 0xf8, 0xcb, 0x8a, 0x07, 0x8b, 0x38,
	0x98, 0xfc, 0x3b, 0x82, 0x24, 0x7c, 0x30, 0x80, 0x2c, 0x9d, 0xe7, 0x61, 0x26, 0xca, 0x34, 0x47,
	0xbd, 0x19, 0xd8, 0x5f, 0xb9, 0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x4e, 0x88, 0xca, 0xa6,
	0x01, 0x00, 0x00,
}