package utils

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userId string) (string, string, error) {
	nowTime := time.Now()
	accessTokenExpireTime := nowTime.Add(time.Hour * 720)
	accessTokenClaims := jwt.MapClaims{
		"id":  userId,
		"iat": nowTime.Unix(),
		"exp": accessTokenExpireTime.Unix(),
	}
	refreshTokenExpireTime := nowTime.Add(time.Hour * 720 * 2)
	refreshTokenClaims := jwt.MapClaims{
		"id":  userId,
		"iat": nowTime.Unix(),
		"exp": refreshTokenExpireTime.Unix(),
	}
	accessToken, errAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errAccessToken != nil {
		return "", "", errAccessToken
	}
	refreshToken, errRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errRefreshToken != nil {
		return "", "", errRefreshToken
	}
	return accessToken, refreshToken, nil
}

// ParseToken parsing token
func ParseToken(token string) (*jwt.MapClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.MapClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
