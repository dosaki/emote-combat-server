package actions

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken default implementation.
func GenerateToken(c buffalo.Context) error {
	u := getUserAuthBody(c)

	users := []models.User{}
	var err error
	err = models.DB.Where("email in (?)", strings.ToLower(strings.TrimSpace(u.Email))).All(&users)
	if err != nil || len(users) == 0 || bcrypt.CompareHashAndPassword([]byte(users[0].PasswordHash), []byte(u.Password)) != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "Unable to authenticate."}))
	}

	expiry := time.Now().Add(time.Minute * 60).Unix()
	claims := jwt.StandardClaims{
		ExpiresAt: expiry,
		Id:        users[0].ID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", ""))
	if err != nil {
		return fmt.Errorf("could not open jwt key, %v", err)
	}
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return fmt.Errorf("could not sign token, %v", err)
	}

	return c.Render(201, r.JSON(map[string]string{
		"token":     tokenString,
		"expiresAt": strconv.FormatInt(expiry, 10),
		"playerId":  users[0].ID.String(),
	}))
}
