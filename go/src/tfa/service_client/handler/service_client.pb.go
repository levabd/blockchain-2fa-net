// Code generated by protoc-gen-go.
// source: service_client.proto
// DO NOT EDIT!

/*
Package handler is a generated protocol buffer package.

It is generated from these files:
	service_client.proto

It has these top-level messages:
	Log
	User
	SCPayload
*/
package handler

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type PayloadType int32

const (
	PayloadType_USER_CREATE   PayloadType = 0
	PayloadType_USER_UPDATE   PayloadType = 1
	PayloadType_CODE_GENERATE PayloadType = 2
	PayloadType_CODE_VERIFY   PayloadType = 3
)

var PayloadType_name = map[int32]string{
	0: "USER_CREATE",
	1: "USER_UPDATE",
	2: "CODE_GENERATE",
	3: "CODE_VERIFY",
}
var PayloadType_value = map[string]int32{
	"USER_CREATE":   0,
	"USER_UPDATE":   1,
	"CODE_GENERATE": 2,
	"CODE_VERIFY":   3,
}

func (x PayloadType) String() string {
	return proto.EnumName(PayloadType_name, int32(x))
}

type Log struct {
	Event      string  `protobuf:"bytes,1,opt" json:"Event,omitempty"`
	Status     string  `protobuf:"bytes,2,opt" json:"Status,omitempty"`
	Code       uint32  `protobuf:"varint,3,opt" json:"Code,omitempty"`
	ExpiredAt  float64 `protobuf:"fixed64,4,opt" json:"ExpiredAt,omitempty"`
	Embeded    bool    `protobuf:"varint,5,opt" json:"Embeded,omitempty"`
	ActionTime float64 `protobuf:"fixed64,6,opt" json:"ActionTime,omitempty"`
	Method     string  `protobuf:"bytes,7,opt" json:"Method,omitempty"`
	Cert       string  `protobuf:"bytes,8,opt" json:"Cert,omitempty"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}

type User struct {
	PhoneNumber    string  `protobuf:"bytes,1,opt" json:"PhoneNumber,omitempty"`
	Uin            string  `protobuf:"bytes,2,opt" json:"Uin,omitempty"`
	Name           string  `protobuf:"bytes,3,opt" json:"Name,omitempty"`
	IsVerified     bool    `protobuf:"varint,4,opt" json:"IsVerified,omitempty"`
	Email          string  `protobuf:"bytes,5,opt" json:"Email,omitempty"`
	Sex            string  `protobuf:"bytes,6,opt" json:"Sex,omitempty"`
	Birthdate      float64 `protobuf:"fixed64,7,opt" json:"Birthdate,omitempty"`
	AdditionalData string  `protobuf:"bytes,8,opt" json:"AdditionalData,omitempty"`
	UpdatedAt      float64 `protobuf:"fixed64,9,opt" json:"UpdatedAt,omitempty"`
	UpdatedBy      string  `protobuf:"bytes,10,opt" json:"UpdatedBy,omitempty"`
	Logs           []*Log  `protobuf:"bytes,11,rep" json:"Logs,omitempty"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}

func (m *User) GetLogs() []*Log {
	if m != nil {
		return m.Logs
	}
	return nil
}

// todo: переименовать sex в gender
type SCPayload struct {
	Action      PayloadType `protobuf:"varint,1,opt,enum=PayloadType" json:"Action,omitempty"`
	PhoneNumber string      `protobuf:"bytes,2,opt" json:"PhoneNumber,omitempty"`
	PayloadUser *User       `protobuf:"bytes,3,opt" json:"PayloadUser,omitempty"`
	PayloadLog  *Log        `protobuf:"bytes,4,opt" json:"PayloadLog,omitempty"`
}

func (m *SCPayload) Reset()         { *m = SCPayload{} }
func (m *SCPayload) String() string { return proto.CompactTextString(m) }
func (*SCPayload) ProtoMessage()    {}

func (m *SCPayload) GetPayloadUser() *User {
	if m != nil {
		return m.PayloadUser
	}
	return nil
}

func (m *SCPayload) GetPayloadLog() *Log {
	if m != nil {
		return m.PayloadLog
	}
	return nil
}

func init() {
	proto.RegisterEnum("PayloadType", PayloadType_name, PayloadType_value)
}
