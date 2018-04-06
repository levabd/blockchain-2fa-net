package handler

import (
	"sawtooth_sdk/processor"
	"github.com/golang/protobuf/proto"
	"fmt"
	"sort"
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

func getCode(event, phoneNumber string, actionTime float64) uint32 {
	var buffer bytes.Buffer
	buffer.WriteString(GetFamilyName())
	buffer.WriteString(event)
	buffer.WriteString(phoneNumber)
	buffer.WriteString(floatToString(actionTime))

	return uint32(crc16.Crc16([]byte(buffer.String())))
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func getLastUserLogIndex(user *User) (uint32, error) {
	var insertIndex uint32 = 0
	if len(user.Logs) == 0 {
		user.Logs = make(map[string]*Log)
	} else {
		keys := make([]int, 0, len(user.Logs))
		for k := range user.Logs {
			s, e := strconv.Atoi(k)
			if e != nil {
				return 0, &processor.InternalError{Msg: "Logs key is invalid"}
			}
			keys = append(keys, s)
		}

		sort.Ints(keys)
		if len(keys) == 1 {
			insertIndex = 0
		} else {
			insertIndex = uint32(keys[len(keys)-1])
		}

		insertIndex = insertIndex + 1
	}

	return insertIndex, nil
}

func addLogToUser(user *User, log *Log, phoneNumber string) error {
	insertIndex, err := getLastUserLogIndex(user)
	if err != nil {
		return &processor.InternalError{Msg: err.Error()}
	}

	var status = SEND_CODE
	if log.Status == RESEND_CODE {
		status = RESEND_CODE
	}

	log.Status = status
	log.Code = getCode(log.Event, phoneNumber, log.ActionTime)

	user.Logs[fmt.Sprint(insertIndex)] = log;
	logger.Debugf("code: %d", log.Code)
	return nil
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
	keys := make([]int, 0, len(user.Logs))
	mapLogsSend := make(map[string]Log)
	for k, item := range user.GetLogs() {
		s, e := strconv.Atoi(k)
		if e != nil {
			return &processor.InternalError{Msg: "Logs key is invalid"}
		}
		keys = append(keys, s)
		if item.Status == SEND_CODE || item.Status == RESEND_CODE {
			mapLogsSend[k] = *item
		}
	}

	keys = make([]int, 0, len(mapLogsSend))
	for k := range mapLogsSend {
		s, e := strconv.Atoi(k)
		if e != nil {
			return &processor.InternalError{Msg: "Logs key is invalid"}
		}
		keys = append(keys, s)
	}
	sort.Ints(keys)

	lastLogWithSendStatusIndex := keys[len(keys)-1]
	t := strconv.Itoa(lastLogWithSendStatusIndex)
	latestLogWithSendCode := mapLogsSend[t]

	if latestLogWithSendCode.ExpiredAt <= log.ActionTime {
		log.Status = EXPIRED
	} else if latestLogWithSendCode.GetCode() == log.Code {
		log.Status = VALID;
	} else {
		log.Status = INVALID;
	}

	insertIndex, err := getLastUserLogIndex(user)
	if err != nil {
		return &processor.InternalError{Msg: err.Error()}
	}
	user.Logs[fmt.Sprint(insertIndex)] = log

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
