package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateFullname = regexp.MustCompile(`^[a-zA-z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("value: %v must contain %v-%v characters ", value, minLength, maxLength)
	}
	return nil
}

func ValidateUserName(name string) error {
	if err := ValidateString(name, 3, 10); err != nil {
		return err
	}
	if !isValidateUsername(name) {
		return fmt.Errorf("value: %v, must contain only lowercase letters, digits, or underscore ", name)
	}
	return nil
}

func ValidatePassword(pasword string) error {
	return ValidateString(pasword, 6, 10)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("value: %v is not valid email address", email)
	}
	return nil
}

func ValidateFullName(name string) error {
	if err := ValidateString(name, 3, 20); err != nil {
		return err
	}
	if !isValidateFullname(name) {
		return fmt.Errorf("value: %v is not allowed, must contain only letters and space", name)
	}
	return nil
}
