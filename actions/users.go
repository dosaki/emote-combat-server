package actions

import (
	"encoding/json"

	"github.com/dosaki/emote_combat_server/helpers"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func getUserBody(c buffalo.Context) models.User {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := models.User{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(400, r.JSON(map[string]string{}))
	}

	c.Session().Set("current_user_id", u.ID)

	return c.Render(201, r.JSON(u))
}

// UserUpdate default implementation.
func UserUpdate(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	users := []models.User{}
	err := models.DB.Where("id in (?)", uuid).All(&users)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting player."}))
	}
	if len(users) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Player not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	body := getUserBody(c)
	user := users[0]
	user.Name = body.Name
	user.Email = body.Name

	if tx.Save(&user) == nil {
		return c.Render(200, r.JSON(user))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// UserList default implementation.
func UserList(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	users := []models.User{}
	var err error

	if perr != nil {
		err = models.DB.All(&users)
		if err == nil {
			return c.Render(200, r.JSON(users))
		}
	} else {
		err = models.DB.Where("id in (?)", uuid).All(&users)
		if len(users) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Player not found."}))
		}
		if err == nil {
			return c.Render(200, r.JSON(users[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Problem getting player(s)."}))
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			return c.Render(403, r.JSON(map[string]string{}))
		}
		return next(c)
	}
}
