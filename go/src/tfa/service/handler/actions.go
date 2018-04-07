package handler

import (
	"sawtooth_sdk/processor"
	"encoding/json"
)


func ApplyCreate(address string, user *User, context *processor.Context) error {
	errors := GetUserValidationErrors(user)
	if len(errors) != 0 {
		errorsMarshaled, _ := json.Marshal(errors)
		return &processor.InvalidTransactionError{Msg: string(errorsMarshaled)}
	}

	newUser := &User{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		IsVerified:  user.IsVerified,
		Birthdate:   user.Birthdate,
		Sex:         user.Sex,
		Uin:         user.Uin,
	}

	err := saveUser(address, newUser, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}

func ApplyUpdate(address string, userUpdateData *User, context *processor.Context) error {
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

	err = saveUser(address, userOld, context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return nil
}