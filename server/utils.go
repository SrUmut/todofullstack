package server

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func checkRegisterReq(username, password, password_repeated string) error {
	if username == "" || password == "" || password_repeated == "" {
		return fmt.Errorf("all fields must be filled")
	}

	if password != password_repeated {
		return fmt.Errorf("passwords don't match")
	}

	if len(password) > 72 {
		return fmt.Errorf("password too long")
	}

	if utf8.RuneCountInString(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if strings.Contains(password, username) {
		return fmt.Errorf("password should not include the username")
	}

	return nil

}
