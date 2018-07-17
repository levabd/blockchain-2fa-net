package handler

import (
	"sawtooth_sdk/processor"
	"github.com/golang/protobuf/proto"
	"fmt"
	"strconv"
	"bytes"
	crc16 "github.com/joaojeronimo/go-crc16"
	"crypto/sha512"
	"strings"
	"encoding/hex"
)

func unpackUser(results map[string][]byte, address string, user *User) error {
	payloadData, exists := results[address]
	if exists && len(payloadData) > 0 {
		err := proto.Unmarshal(payloadData, user)
		if err != nil {
			return &processor.InternalError{
				Msg: fmt.Sprint("Failed to unmarshal 2FA Service Client: %v", err)}
		}
	}
	return nil
}

func Hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func addLogToUser(user *User, log *Log, phoneNumber, familyName string) {
	var status = SEND_CODE
	if log.Status == RESEND_CODE {
		status = RESEND_CODE
	}

	log.Status = status

	var buffer bytes.Buffer
	buffer.WriteString(familyName)
	buffer.WriteString(log.Event)
	buffer.WriteString(phoneNumber)
	buffer.WriteString(floatToString(log.ActionTime))

	log.Code = uint32(crc16.Crc16([]byte(buffer.String())))
	user.Logs = append(user.Logs, log)
}

func getUserByAddress(address string, context *processor.Context) (*User, error) {
	results, err := context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	user := User{}
	err = unpackUser(results, address, &user)
	if err != nil {
		return nil, &processor.InternalError{
			Msg: fmt.Sprint("Failed to load user: ", err),
		}
	}
	return &user, nil
}

func verify(user *User, log *Log, phoneNumber string) (error) {
	mapLogsSend := make([]Log, 0)
	for _, item := range user.GetLogs() {
		if item.Status == SEND_CODE || item.Status == RESEND_CODE {
			mapLogsSend = append(mapLogsSend, *item)
		}
	}
	logger.Info(log.Status)
	if log.Status == REJECT {
		user.Logs = append(user.Logs, log)
		return nil
	}
	latestLogWithSendCode := mapLogsSend[len(mapLogsSend)-1]
	if latestLogWithSendCode.ExpiredAt <= log.ActionTime {
		log.Status = EXPIRED
	} else if latestLogWithSendCode.GetCode() == log.Code {
		log.Status = VALID;
	} else {
		log.Status = INVALID;
	}
	user.Logs = append(user.Logs, log)
	return nil
}

func saveUser(address string, user *User, context *processor.Context) error {
	data, err := proto.Marshal(user)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to serialize Account:", err)}
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
