package service

import (
	"1dv027/aad/internal/dto"
	"1dv027/aad/internal/model"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

type JwtService struct {
	signingKey string
}

func NewJwtService(signingKey string) JwtService {
	return JwtService{
		signingKey: signingKey,
	}
}

func (j JwtService) GenerateJwt(username string, id int, userType model.UserRole) (string, error) {
	var mySigningKey = []byte(j.signingKey)

	claims := CustomClaims{
		Username: username,
		UserType: string(userType),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "DogAdoptionApp",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			Subject:   fmt.Sprintf("%d", id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j JwtService) ValidateToken(tokenString string) (dto.UserCredentials, error) {
	claims := &CustomClaims{}

	validatedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})

	if err != nil {
		return dto.UserCredentials{}, err
	}

	if !validatedToken.Valid {
		return dto.UserCredentials{}, fmt.Errorf("invalid token")
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return dto.UserCredentials{}, fmt.Errorf("could not format subject to integer")
	}

	userRole, err := model.StringToUserRole(claims.UserType)
	if err != nil {
		return dto.UserCredentials{}, err
	}

	userCredentials := dto.UserCredentials{
		Id:       id,
		Username: claims.Username,
		UserRole: userRole,
	}

	return userCredentials, nil
}
