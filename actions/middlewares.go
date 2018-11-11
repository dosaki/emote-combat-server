package actions

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dosaki/emote_combat_server/helpers"
	"github.com/dosaki/emote_combat_server/services"
	"github.com/gobuffalo/buffalo"
)

func getToken(c buffalo.Context) (*jwt.Token, error) {
	tokenString := c.Request().Header.Get("Authorization")
	if len(tokenString) == 0 {
		return nil, c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
	}

	if !services.TokenIsValid(tokenString) {
		return nil, c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid token pair"))
	}

	token, err := jwt.Parse(tokenString, services.ParseToken)
	if err != nil {
		return nil, c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid user/token pair"))
	}
	return token, nil
}

func checkClaims(c buffalo.Context, token *jwt.Token, checkUser bool) bool {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u, err := services.GetUserByUUID(claims["jti"].(string))
		if err != nil {
			return false
		}

		if checkUser {
			uuid, perr := helpers.Param(c, "player_id")
			if perr == nil && u.ID.String() == uuid {
				c.Set("user", u)
				return true
			}
			return false
		}
		c.Set("user", u)
		return true
	}
	return false
}

// PlayerRestrictedHandlerMiddleware - handles restricted actions by making sure they have a valid token and are acting on the correct player
func PlayerRestrictedHandlerMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		token, buffaloErr := getToken(c)
		if buffaloErr != nil {
			return buffaloErr
		}

		if checkClaims(c, token, true) {
			return next(c)
		}

		return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid token or unauthorized action"))
	}
}

// RestrictedHandlerMiddleware - handles restricted actions by making sure they have a valid token and are acting on the correct player
func RestrictedHandlerMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		token, buffaloErr := getToken(c)
		if buffaloErr != nil {
			return buffaloErr
		}

		if checkClaims(c, token, false) {
			return next(c)
		}

		return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid token or unauthorized action"))
	}
}
