package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)


func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3 ,100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("%s Not a valid username", value)
	}
	
	return nil
}
func ValidateFullUsername(value string) error {
	if err := ValidateString(value, 3 ,100); err != nil {
		return err
	}

	if !isValidFullname(value) {
		return fmt.Errorf("%s Not a valid full name", value)
	}
	
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3 ,100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("%s is not a valid email address", value)
	}
	return nil
}