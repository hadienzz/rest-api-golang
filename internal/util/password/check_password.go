package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}

	return nil
}
