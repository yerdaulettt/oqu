package service

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func verifyPassword(hashPswd string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPswd), []byte(password))
}
