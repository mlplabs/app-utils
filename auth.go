package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// ValidPassword сравнивает хэш и пароль
func ValidPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = errors.New("wrong username or password")
	}
	return err == nil, err
}

// CreateHash возвращает хэш по строке пароля
func CreateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CreateTokenPair возвращает пару токенов
//   secretKey - ключ подписи
//   accExpSec, refExpSec - время жизни токенов соответственно
//   payload - дополнительные поля
func CreateTokenPair(payload map[string]string, secretKey string, accExpSec int64, refExpSec int64) (*TokenPair, error) {
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["exp"] = time.Now().Add(time.Duration(accExpSec) * time.Second).Unix() //Token expires after 15 minutes

	refreshClaims := jwt.MapClaims{}
	refreshClaims["exp"] = time.Now().Add(time.Duration(refExpSec) * time.Second).Unix() //Token expires after 12 hour

	if payload != nil {
		for k, v := range payload {
			accessClaims[k] = v
			refreshClaims[k] = v
		}
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("create signed access token string %v", err)
	}
	refreshTokenSting, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("create signed refresh token string %v", err)
	}
	return &TokenPair{accessTokenString, "bearer", accExpSec, refreshTokenSting}, nil
}

func ReadToken(secretKey string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}
	return nil, err
}

func TokenValid(r *http.Request) error {
	// TODO: tokenSecret!
	tokenSecret := ""
	tokenString := ExtractToken(r)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
