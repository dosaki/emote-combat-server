package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/dosaki/emote_combat_server/services"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken default implementation.
func GenerateToken(c buffalo.Context) error {
	u := getUserAuthBody(c)

	var users []models.User
	var err error
	err = models.DB.Where("email in (?)", strings.ToLower(strings.TrimSpace(u.Email))).All(&users)
	if err != nil || len(users) == 0 || bcrypt.CompareHashAndPassword([]byte(users[0].PasswordHash), []byte(u.Password)) != nil {
		return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"message": "Unable to authenticate."}))
	}

	expiry := time.Now().Add(time.Minute * 60).Unix()
	claims := jwt.StandardClaims{
		ExpiresAt: expiry,
		Id:        users[0].ID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", ""))
	if err != nil {
		fmt.Println("could not open jwt key", err)
		return c.Render(http.StatusServiceUnavailable, r.JSON(map[string]string{"message": "Token generation is unavailable. Please contact the administrator."}))
	}
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Println("could not sign token", err)
		return c.Render(http.StatusServiceUnavailable, r.JSON(map[string]string{"message": "Token generation is unavailable. Please contact the administrator."}))
	}

	return c.Render(201, r.JSON(map[string]string{
		"token":     tokenString,
		"expiresAt": strconv.FormatInt(expiry, 10),
		"playerId":  users[0].ID.String(),
	}))
}

// DestroyToken - Destroys the current token effectively logging out the user
func DestroyToken(c buffalo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if len(tokenString) == 0 {
		return c.Render(http.StatusBadRequest, r.JSON(map[string]string{"message": "Attempting to log out a user who's already logged out."}))
	}

	_, err := jwt.Parse(tokenString, services.ParseToken)
	if err != nil {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("Invalid user/token pair"))
	}
	services.InvalidateToken(tokenString)
	return c.Render(http.StatusOK, r.JSON(map[string]string{}))
}
