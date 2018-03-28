package handler

import (
	"regexp"
	"encoding/json"
)

var RE_EMAIL = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_" + "{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

func GetUserValidationErrors(user User) []string {
	var errors = []string{}

	if user.Name == "" {
		errors = append(errors, "Name is required")
	}

	if user.PhoneNumber == "" {
		errors = append(errors, "PhoneNumber is required")
	}

	if ! RE.MatchString(user.PhoneNumber) {
		errors = append(errors, "PhoneNumber format is invalid")
	}

	if user.Sex == "" {
		errors = append(errors, "Sex is required")
	}

	if user.Sex != "female" && user.Sex != "male" {
		errors = append(errors, "Sex word is wrong")
	}

	if user.Email == "" {
		errors = append(errors, "Email is required")
	}

	if ! RE_EMAIL.MatchString(user.Email) {
		errors = append(errors, "Email format is invalid")
	}

	if user.AdditionalData != "" && !isJSONString(user.AdditionalData) {
		errors = append(errors, "AdditionalData must be a valid json string ")
	}

	return errors;
}

func GetLogValidationErrors(log Log) []string {
	var errors = []string{}

	if log.Event == "" {
		errors = append(errors, "Event is required")
	}

	if log.Status == "" {
		errors = append(errors, "Status is required")
	}

	if log.ExpiredAt == 0 {
		errors = append(errors, "ExpiredAt is required")
	}

	if log.ActionTime == 0 {
		errors = append(errors, "ActionTime is required")
	}

	if log.Embeded && log.Cert == "" {
		errors = append(errors, "Cert is required then Embeded true")
	}

	return errors;
}
