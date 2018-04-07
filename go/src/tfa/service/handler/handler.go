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
	"regexp"
	"github.com/golang/protobuf/proto"
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

var RE = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

const (
	NAME        = "tfa"
	VERSION     = "0.1"
)

func (self *SCHandler) FamilyName() string {
	return NAME
}

func (self *SCHandler) FamilyVersions() []string {
	return []string{VERSION}
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

	phoneNumber := payload.PhoneNumber

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

	hashedPhoneNumber := Hexdigest(phoneNumber)
	address := self.namespace + hashedPhoneNumber[len(hashedPhoneNumber)-64:]

	switch payload.Action {
	case PayloadType_USER_CREATE:
		return ApplyCreate(address, payload.GetPayloadUser(), context)
	case PayloadType_USER_UPDATE:
		return ApplyUpdate(address, payload.GetPayloadUser(), context)
	default:
		return &processor.InternalError{
			Msg: fmt.Sprintf("Verb must be register, update, setPushToken ot isVerified: not  %s", payload.Action),
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
