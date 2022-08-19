package util

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey []byte

type claimValue struct {
	Username string
	jwt.StandardClaims
}

func GenerateJWT(username string, id int) (token string, err error) {
	expirationTime := time.Now().Add(time.Hour * 12).Unix() // 12Hrs.
	claimValues := &claimValue{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    strconv.Itoa(id),
			ExpiresAt: expirationTime,
		},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claimValues)
	return claim.SignedString(secretKey)
}

func ValidateToken(signToken string) (claim jwt.MapClaims, err error) {
	token, err := jwt.Parse(
		signToken,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				err = errors.New("Unexpected signing token")
				return nil, err
			}
			return []byte(secretKey), nil
		},
	)

	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claim, nil
	}

	return nil, errors.New(err.Error())
}
