package handler

import (
	"sawtooth_sdk/processor"
	"encoding/json"
	"fmt"
)

const (
	RESEND_CODE = "RESEND_CODE"
	SEND_CODE   = "SEND_CODE"
	EXPIRED     = "EXPIRED"
	VALID       = "VALID"
	INVALID     = "INVALID"
	REJECT      = "REJECT"
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
		//UpdatedAt:      user.UpdatedAt,
		//UpdatedBy:      user.UpdatedBy,
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

	//if userUpdateData.UpdatedBy == "" {
	//	return &processor.InvalidTransactionError{Msg: string("UpdatedBy must be set.")}
	//}

	userOld, err := getUserByAddress(address, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	userOld.Name = userUpdateData.Name
	userOld.PhoneNumber = userUpdateData.PhoneNumber
	userOld.Email = userUpdateData.Email
	userOld.Sex = userUpdateData.Sex
	userOld.Uin = userUpdateData.Uin
	// todo finish update tasks
	//userOld.UpdatedAt = userUpdateData.UpdatedAt
	//userOld.UpdatedBy = userUpdateData.UpdatedBy
	userOld.IsVerified = userUpdateData.IsVerified
	userOld.Birthdate = userUpdateData.Birthdate
	userOld.AdditionalData = userUpdateData.AdditionalData

	err = saveUser(address, userOld, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}

func ApplyCodeGeneration(
	log *Log,
	context *processor.Context,
	address, phoneNumber, familyName string) error {

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

	addLogToUser(user, log, phoneNumber, familyName)
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

	err = verify(user, log, phoneNumber)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	err = saveUser(address, user, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}
