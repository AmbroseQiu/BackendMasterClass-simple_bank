package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateFullname = regexp.MustCompile(`^[a-zA-z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %v-%v characters ", minLength, maxLength)
	}
	return nil
}

func ValidateUserName(name string) error {
	if err := ValidateString(name, 3, 10); err != nil {
		return err
	}
	if !isValidateUsername(name) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore ")
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
		return fmt.Errorf("is not valid email address")
	}
	return nil
}

func ValidateFullName(name string) error {
	if err := ValidateString(name, 3, 10); err != nil {
		return err
	}
	if !isValidateFullname(name) {
		return fmt.Errorf("must contain only letters and space")
	}
	return nil
}
