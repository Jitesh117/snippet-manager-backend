package helper

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/Jitesh117/snippet-manager-backend/constants"
	"github.com/Jitesh117/snippet-manager-backend/models"
)

func ValidateSnippet(snippet models.Snippet) error {
	if snippet.Title == "" {
		return fmt.Errorf(constants.ErrEmptyTitle)
	}
	if snippet.Language == "" {
		return fmt.Errorf(constants.ErrEmptyLanguage)
	}
	if snippet.Content == "" {
		return fmt.Errorf(constants.ErrEmptyContent)
	}
	return nil
}

func validateEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf(constants.ErrPasswordTooShort)
	}
	if len(password) > 20 {
		return fmt.Errorf(constants.ErrPasswordTooLong)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(`!@#$%^&*(),.?":{}|<>`, char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}

	return nil
}

func ValidateUser(user models.User) error {
	if user.UserName == "" {
		return fmt.Errorf(constants.ErrEmptyUsername)
	}
	if len(user.UserName) < 3 {
		return fmt.Errorf(constants.ErrUsernameTooShort)
	}

	if len(user.UserName) > 30 {
		return fmt.Errorf(constants.ErrUsernameTooLong)
	}

	if user.Email == "" {
		return fmt.Errorf(constants.ErrEmptyContent)
	}
	if !validateEmail(user.Email) {
		return fmt.Errorf(constants.ErrInvalidEmailFormat)
	}

	if user.Password == "" {
		return fmt.Errorf(constants.ErrEmptyPassword)
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}

	return nil
}
