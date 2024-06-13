package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

type CryptographyService struct {
	secretKey []byte
}

func NewCryptographyService(secretKey string) (CryptographyService, error) {
	if len(secretKey) != 32 {
		return CryptographyService{}, errors.New("invalid secretKey length: must be exactly 32 bytes for AES-256")
	}
	return CryptographyService{
		secretKey: []byte(secretKey),
	}, nil
}

func (c CryptographyService) EncryptPlainText(plainText string) (string, error) {
	block, err := aes.NewCipher(c.secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (c CryptographyService) DecryptCipherText(cipherText string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func (c CryptographyService) HashPassword(unhashedPassword string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(unhashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (c CryptographyService) ComparePasswords(hashedPassword, passwordAttempt string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordAttempt))
}
