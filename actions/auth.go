package actions

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		return c.Render(422, r.JSON(map[string]string{"message": "invalid email/password"}))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	c.Session().Set("current_user_id", u.ID)

	claims := jwt.MapClaims{}
	claims["userid"] = u.ID
	expiry := time.Now().Add(time.Minute * 5).Unix()
	claims["exp"] = expiry

	key := envy.Get("JWT_SECRET", "ERRONEOUS KEY")
	fmt.Println(key)
	if key == "ERRONEOUS KEY" {
		log.Fatal(errors.Wrap(err, "bad key"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		log.Fatal(errors.Wrap(err, "error obtaining signed key"))
	}
	return c.Render(201, r.JSON(map[string]string{
		"token":    tokenString,
		"playerId": u.ID.String(),
		"expiry":   strconv.FormatInt(expiry, 10),
	}))
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	return c.Render(200, r.JSON(map[string]string{}))
}
