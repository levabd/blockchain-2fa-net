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
	"fmt"
	"sawtooth_sdk/logging"
	"sawtooth_sdk/processor"
	"sawtooth_sdk/protobuf/processor_pb2"
	"github.com/golang/protobuf/proto"
	"encoding/json"
)

var logger *logging.Logger = logging.Get()

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

func (self *SCHandler) GetFamilyName() string {
	return self.family_name
}

func (self *SCHandler) SetFamilyName(name string) {
	self.family_name = name
}

func (self *SCHandler) SetFamilyVersion(version string) {
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
	payload, err := unpackPayload(request.GetPayload())
	if err != nil {
		return err
	}

	logger.Debugf("service client tp txn %v: action %v",
		request.Signature, payload.GetAction())

	errors := GetPayloadErrors(payload)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	hashedPhoneNumber := Hexdigest(payload.GetPhoneNumber())
	address := self.namespace + hashedPhoneNumber[len(hashedPhoneNumber)-64:]

	switch payload.GetAction() {
	case PayloadType_USER_CREATE:
		return ApplyCreateUser(address, payload.GetPayloadUser(), context)
	case PayloadType_USER_UPDATE:
		return ApplyUpdateUser(address, payload.GetPayloadUser(), context)
	case PayloadType_CODE_GENERATE:
		return ApplyCodeGeneration(
			payload.GetPayloadLog(),
			context,
			address,
			payload.GetPhoneNumber(),
			self.GetFamilyName())
	case PayloadType_CODE_VERIFY:
		return ApplyVerification(address, payload.GetPayloadLog(), context, payload.GetPhoneNumber())
	default:
		return &processor.InternalError{
			Msg: fmt.Sprintf("Verb must be register, update, "+
				"setPushToken ot isVerified: not  %s", payload.GetAction()),
		}
	}

	return nil
}

func unpackPayload(payloadData []byte) (*SCPayload, error) {
	payload := &SCPayload{}
	err := proto.Unmarshal(payloadData, payload)
	if err != nil {
		return nil, &processor.InternalError{
			Msg: fmt.Sprint("Failed to unmarshal 2FA Service Client: %v", err)}
	}
	return payload, nil
}
