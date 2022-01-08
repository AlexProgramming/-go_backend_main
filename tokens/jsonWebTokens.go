package tokens

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt"
)

// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.

// from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb @ denisbrodbeck
func GenerateRandomString(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return []byte(""), err
		}
		ret[i] = letters[num.Int64()]
	}

	return []byte(string(ret)), nil
}

var SecretKey = []byte("")

func CreateJWT(issuer string) (string, error) {
	err := new(error)
	SecretKey, *err = GenerateRandomString(32)
	if *err != nil {
		return "", *err
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(), // align with authorizedController login param
		Issuer:    issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss := new(string)
	*ss, *err = token.SignedString(SecretKey)

	return *ss, *err
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if ok && token.Valid {
		return claims.Issuer, nil
	}

	return "", err
}
