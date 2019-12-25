package jwt

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken 创建Token
func CreateToken(key string, m map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	for index, val := range m {
		claims[index] = val
	}
	token.Claims = claims
	return token.SignedString([]byte(key))
}

// ParseToken 解析Token
func ParseToken(tokenString string, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err == nil {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			m := make(map[string]interface{})
			for index, val := range claims {
				m[index] = val
			}
			return m, true
		}
	}
	return nil, false
}
