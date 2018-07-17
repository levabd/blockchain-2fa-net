package handler

import (
	"sawtooth_sdk/processor"
	"github.com/golang/protobuf/proto"
	"fmt"
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
