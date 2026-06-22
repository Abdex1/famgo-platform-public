/*
# STEP 10 — CREATE PASSWORD SERVICE

services/auth-service/internal/domain/services/password_service.go


# FILE: password_service.go
*/
package services

import "golang.org/x/crypto/bcrypt"

type PasswordService struct {}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s *PasswordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(hash), err
}

func (s *PasswordService) Compare(
	hash string,
	password string,
) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
