package services

import (
	"fmt"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

var invalidTokens = []string{}

// TokenIsValid - Checks against a list to make sure a token is valid
func TokenIsValid(token string) bool {
	for i := range invalidTokens {
		if invalidTokens[i] == token {
			return false
		}
	}
	return true
}

// InvalidateToken - adds token to the list of invalid tokens
func InvalidateToken(token string) {
	invalidTokens = append(invalidTokens, token)
	fmt.Println(invalidTokens)
}

func ParseToken(token *jwt.Token) (interface{}, error) {
	// check signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	// read the key
	mySignedKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", ""))
	if err != nil {
		return nil, fmt.Errorf("could not open jwt key, %v", err)
	}
	return mySignedKey, nil
}
