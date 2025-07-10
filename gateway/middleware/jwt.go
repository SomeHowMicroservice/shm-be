package middleware

import (
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("lỗi parse token: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, customErr.ErrInvalidToken
}

func ExtractToken(claims jwt.MapClaims) (string, []string, error) {
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", nil, customErr.ErrUserIdNotFound
	}
	rawRoles, ok := claims["roles"].([]interface{}) 
	if !ok {
		return "", nil, customErr.ErrRolesNotFound
	}
	userRoles := make([]string, len(rawRoles)) 
	for i, r := range rawRoles {
		userRoles[i] = fmt.Sprintf("%v", r)
	}
	return userID, userRoles, nil
}
