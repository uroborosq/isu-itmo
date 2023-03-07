package app

import (
	"errors"
	"github.com/uroborosq/isu/internal/isu/adapters"
	"net/mail"
	"regexp"
)

func ValidateFullInfo(user adapters.UserFullData) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}
	phoneRe := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	if !phoneRe.MatchString(user.PhoneNumber) {
		return errors.New("phone number is in wrong format")
	}

	if user.Role > 2 || user.Role < 0 {
		return errors.New("given role is not supported")
	}
	return nil
}
