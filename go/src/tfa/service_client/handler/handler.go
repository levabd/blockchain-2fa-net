/**
 * Copyright 2017 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package handler

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	cbor "github.com/brianolson/cbor_go"
	"sawtooth_sdk/logging"
	"sawtooth_sdk/processor"
	"sawtooth_sdk/protobuf/processor_pb2"
	"strings"
	"regexp"
	"encoding/json"
	"sort"
	crc16 "github.com/joaojeronimo/go-crc16"
	"bytes"
	"strconv"
)

var logger *logging.Logger = logging.Get()

type Log struct {
	Event      string `json:"Event" bson:"Event" binding:"required"`
	Status     string `json:"Status" bson:"Status" binding:"required"`
	Code       uint16 `json:"Code" bson:"Code" binding:"required"`
	ExpiredAt  uint64 `json:"ExpiredAt" bson:"ExpiredAt" binding:"required"`
	Embeded    bool   `json:"Embeded" bson:",omitempty"`
	ActionTime uint64 `json:"ActionTime" bson:"ActionTime" binding:"required"`
	Method     string `json:"Method" bson:"Method" binding:"required"`
	Cert       string `json:"Cert" bson:"Cert" bson:",omitempty"`
}

type User struct {
	PhoneNumber    string      `json:"PhoneNumber" bson:"PhoneNumber" binding:"required"`
	Uin            uint64      `json:"Uin" bson:"Uin" binding:"required"`
	Name           string      `json:"Name" bson:"Name" binding:"required"`
	IsVerified     bool        `json:",omitempty"`
	Email          string      `json:"Email" bson:"Email" binding:"required"`
	Sex            string      `json:"Sex" bson:"Sex" binding:"required"`
	Birthdate      uint64      `json:"Birthdate" bson:"Birthdate" binding:"required"`
	AdditionalData string      `json:",omitempty"`
	Logs           map[int]Log `json:",omitempty"`
}

type SCPayload struct {
	Action      string
	PhoneNumber string
	User        User
	Log         Log `json:",omitempty"`
}

type SCHandler struct {
	namespace      string
	family_name    string
	family_version string
}

func NewHandler(namespace string) *SCHandler {
	return &SCHandler{
		namespace: namespace,
	}
}

var RE = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

const (
	RESEND_CODE = "RESEND_CODE"
	SEND_CODE   = "SEND_CODE"
	EXPIRED   = "EXPIRED"
	VALID   = "VALID"
	INVALID   = "INVALID"
)

var family_name, family_version string

func (self *SCHandler) SetFamilyName(name string) {
	self.family_name = name
}

func (self *SCHandler) SetFamilyVersion(version string)  {
	self.family_version = version
}

func (self *SCHandler) FamilyName() string {
	return self.family_name
}

func (self *SCHandler) FamilyVersions() []string {
	return []string{self.family_version}
}

func (self *SCHandler) Namespaces() []string {
	return []string{self.namespace}
}

func (self *SCHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	payloadData := request.GetPayload()
	if payloadData == nil {
		return &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}
	var payload SCPayload
	err := DecodeCBOR(payloadData, &payload)
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprint("Failed to decode payload: ", err),
		}
	}

	if err != nil {
		logger.Error("Bad payload: ", payloadData)
		return &processor.InternalError{Msg: fmt.Sprint("Failed to decode payload: ", err)}
	}

	action := payload.Action
	phoneNumber := payload.PhoneNumber

	if action == "" {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Action is required")}
	}

	if phoneNumber == "" {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("PhoneNumber is required"),
		}
	}

	if ! RE.MatchString(phoneNumber) {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("PhoneNumber %v has invalid format", phoneNumber),
		}
	}

	if !(action == "create" || action == "update" || action == "delete" || action == "addLog" || action == "verify" || action == "setPushToken") {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid action: %v", action)}
	}

	hashedPhoneNumber := Hexdigest(phoneNumber)
	address := self.namespace + hashedPhoneNumber[len(hashedPhoneNumber)-64:]

	results, err := context.GetState([]string{address})
	if err != nil {
		return err
	}

	var user = User{}
	switch action {
	case "create":
		errors := GetUserValidationErrors(payload.User)
		if len(errors) != 0 {
			errorsMarshaled, _ := json.Marshal(errors)
			return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
		}
		user = payload.User
		break
	case "update":
		errors := GetUserValidationErrors(payload.User)
		if len(errors) != 0 {
			errorsMarshaled, _ := json.Marshal(errors)
			return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
		}

		data, exists := results[address]
		if exists && len(data) > 0 {
			fmt.Print("Data: ", data)
			err = DecodeCBOR(data, &user)
			if err != nil {
				return &processor.InternalError{
					Msg: fmt.Sprint("Failed to decode data: ", err),
				}
			}
		} else {
			user = User{}
		}

		user = self.UpdateUser(user, payload.User)
		break
	case "addLog":
		if payload.Log == (Log{}) {
			return &processor.InternalError{
				Msg: fmt.Sprint("Payload does not contain Log model"),
			}
		}
		errors := GetLogValidationErrors(payload.Log)
		if len(errors) != 0 {
			errorsMarshaled, _ := json.Marshal(errors)
			return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
		}

		user = self.AddLogToUser(user, payload.Log, phoneNumber)
		break
	case "varify":
		if payload.Log == (Log{}) {
			return &processor.InternalError{
				Msg: fmt.Sprint("Payload does not contain Log model"),
			}
		}
		errors := GetLogValidationErrors(payload.Log)
		if len(errors) != 0 {
			errorsMarshaled, _ := json.Marshal(errors)
			return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
		}
		user = self.Verify(user, payload.Log, phoneNumber)
		break
	default:
		return &processor.InternalError{
			Msg: fmt.Sprintf("Verb must be register, update, setPushToken ot isVerified: not  %s", action),
		}
	}

	data, err := EncodeCBOR(user)
	if err != nil {
		return &processor.InternalError{
			Msg: fmt.Sprint("Failed to encode new map:", err),
		}
	}

	addresses, err := context.SetState(map[string][]byte{
		address: data,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}

	return nil
}

func (self *SCHandler) UpdateUser(userOld, userUpdateData User) User {
	userOld.Name = userUpdateData.Name
	userOld.PhoneNumber = userUpdateData.PhoneNumber
	userOld.Email = userUpdateData.Email
	userOld.Sex = userUpdateData.Sex
	userOld.Birthdate = userUpdateData.Birthdate
	userOld.AdditionalData = userUpdateData.AdditionalData
	return userOld
}

func (self *SCHandler) Verify(user User, log Log, phoneNumber string) User {
	keys := make([]int, 0, len(user.Logs))
	mapLogsSend := make(map[int]Log)
	for k, item := range user.Logs {
		keys = append(keys, k)
		if item.Status ==SEND_CODE || item.Status == RESEND_CODE {
			mapLogsSend[k]=item
		}
	}
	sort.Ints(keys)
	lastLogIndex := keys[len(keys)-1]

	keys = make([]int, 0, len(mapLogsSend))
	for k := range mapLogsSend {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	lastLogWithSendStatusIndex := keys[len(keys)-1]
	latestLogWithSendCode := mapLogsSend[lastLogWithSendStatusIndex]


	var buffer bytes.Buffer
	buffer.WriteString(self.FamilyName())
	buffer.WriteString(log.Event)
	buffer.WriteString(phoneNumber)
	buffer.WriteString(strconv.Itoa(int(log.ActionTime)))
	log.Code = crc16.Crc16([]byte(buffer.String()))

	if latestLogWithSendCode.ExpiredAt <= log.ActionTime{
		log.Status = EXPIRED
	} else if latestLogWithSendCode.Code == self.getCode(log.Event,phoneNumber, log.ActionTime){
		log.Status = VALID;
	} else{
		log.Status = INVALID;
	}

	user.Logs[lastLogIndex]=log

	return user
}

func (self *SCHandler) AddLogToUser(user User, log Log, phoneNumber string) User {
	logIndex := 0
	if len(user.Logs) == 0 {
		user.Logs = make(map[int]Log)
	} else {
		keys := make([]int, 0, len(user.Logs))
		for k := range user.Logs {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		logIndex = keys[len(keys)-1]
	}

	var status = SEND_CODE
	if log.Status == RESEND_CODE {
		status = RESEND_CODE
	}

	log.Status = status

	log.Code = self.getCode(log.Event,phoneNumber, log.ActionTime)

	if logIndex == 0 {
		user.Logs[0] = log;
	} else {
		user.Logs[ logIndex+1] = log;
	}
	logger.Debugf("code: %v", log.Code)
	return user
}

func (self *SCHandler) getCode(event, phoneNumber string, actionTime uint64) uint16{
	var buffer bytes.Buffer
	buffer.WriteString(self.FamilyName())
	buffer.WriteString(event)
	buffer.WriteString(phoneNumber)
	buffer.WriteString(strconv.Itoa(int(actionTime)))

	return  crc16.Crc16([]byte(buffer.String()))
}

func EncodeCBOR(value interface{}) ([]byte, error) {
	data, err := cbor.Dumps(value)
	return data, err
}

func DecodeCBOR(data []byte, pointer interface{}) error {
	defer func() error {
		if recover() != nil {
			return &processor.InvalidTransactionError{Msg: "Failed to decode payload"}
		}
		return nil
	}()
	err := cbor.Loads(data, pointer)
	if err != nil {
		return err
	}
	return nil
}

func Hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}
