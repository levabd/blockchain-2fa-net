package handler

import (
	"sawtooth_sdk/processor"
	"encoding/json"
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

const (
	RESEND_CODE = "RESEND_CODE"
	SEND_CODE   = "SEND_CODE"
	EXPIRED     = "EXPIRED"
	VALID       = "VALID"
	INVALID     = "INVALID"
)

func ApplyCreateUser(address string, user *User, context *processor.Context) error {
	errors := GetUserValidationErrors(user)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	newUser := &User{
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Email:          user.Email,
		IsVerified:     user.IsVerified,
		Birthdate:      user.Birthdate,
		Sex:            user.Sex,
		Logs:           user.Logs,
		Uin:            user.Uin,
		AdditionalData: user.AdditionalData,
	}

	err := saveUser(address, newUser, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}

func ApplyUpdateUser(address string, userUpdateData *User, context *processor.Context) error {
	errors := GetUserValidationErrors(userUpdateData)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	userOld, err := getUserByAddress(address, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	userOld.Name = userUpdateData.Name
	userOld.PhoneNumber = userUpdateData.PhoneNumber
	userOld.Email = userUpdateData.Email
	userOld.Sex = userUpdateData.Sex
	userOld.Birthdate = userUpdateData.Birthdate
	userOld.AdditionalData = userUpdateData.AdditionalData

	err = saveUser(address, userOld, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}

func ApplyCodeGeneration(address string, log *Log, context *processor.Context, phoneNumber string) error {
	if log == (&Log{}) {
		return &processor.InternalError{
			Msg: fmt.Sprint("Payload does not contain Log model"),
		}
	}

	errors := GetLogValidationErrors(log)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	user, err := getUserByAddress(address, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	err = addLogToUser(user, log, phoneNumber)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	err = saveUser(address, user, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}

func ApplyVerification(address string, log *Log, context *processor.Context, phoneNumber string) error {
	if log == (&Log{}) {
		return &processor.InternalError{
			Msg: fmt.Sprint("Payload does not contain Log model"),
		}
	}

	errors := GetLogValidationErrors(log)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	user, err := getUserByAddress(address, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	err = Verify(user, log, phoneNumber)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	err = saveUser(address, user, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
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

func addLogToUser(user *User, log *Log, phoneNumber string) error {
	logIndex := 0
	if len(user.Logs) == 0 {
		user.Logs = make(map[string]*Log)
	} else {
		keys := make([]int, 0, len(user.Logs))
		for k := range user.Logs {
			s, e := strconv.Atoi(k)
			if e != nil {
				return &processor.InternalError{Msg: "Logs key is invalid"}
			}
			keys = append(keys, s)
		}
		sort.Ints(keys)
		logIndex = keys[len(keys)-1]
	}

	var status = SEND_CODE
	if log.Status == RESEND_CODE {
		status = RESEND_CODE
	}

	log.Status = status
	log.Code = getCode(log.Event, phoneNumber, log.ActionTime)

	if logIndex == 0 {
		user.Logs["0"] = log;
	} else {
		t := strconv.Itoa(logIndex + 1)
		user.Logs[t] = log;
	}

	logger.Info("code: %v", log.Code)
	return nil
}

func Verify(user *User, log *Log, phoneNumber string) (error) {
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
	sort.Ints(keys)
	lastLogIndex := keys[len(keys)-1]

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

	var buffer bytes.Buffer
	buffer.WriteString(GetFamilyName())
	buffer.WriteString(log.Event)
	buffer.WriteString(phoneNumber)
	buffer.WriteString(strconv.Itoa(int(log.ActionTime)))
	log.Code = uint32(crc16.Crc16([]byte(buffer.String())))

	if latestLogWithSendCode.ExpiredAt <= log.ActionTime {
		log.Status = EXPIRED
	} else if latestLogWithSendCode.GetCode() == getCode(log.GetEvent(), phoneNumber, log.GetActionTime()) {
		log.Status = VALID;
	} else {
		log.Status = INVALID;
	}

	t = strconv.Itoa(lastLogIndex)
	user.Logs[t] = log

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
	buffer.WriteString(FloatToString(actionTime))

	return uint32(crc16.Crc16([]byte(buffer.String())))
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
