package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dosaki/emote_combat_server/services"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

func parseToken(token *jwt.Token) (interface{}, error) {
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

// RestrictedHandlerMiddleware - handles restricted actions by making sure they ahve a valid token
func RestrictedHandlerMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if len(tokenString) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
		}

		token, err := jwt.Parse(tokenString, parseToken)
		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid user/token pair"))
		}
		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			u, err := services.GetUserByUUID(claims["jti"].(string))
			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid user/token pair"))
			}
			c.Set("user", u)
			return next(c)
		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}
	}
}
