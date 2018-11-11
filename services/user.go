package services

import (
	"errors"

	"github.com/dosaki/emote_combat_server/models"
	UUID "github.com/gobuffalo/uuid"
)

// GetUserByUUID - returns a user based on a UUID
func GetUserByUUID(uuidString string) (models.User, error) {
	uuid, err := UUID.FromString(uuidString)
	users := []models.User{}

	if err == nil {
		err = models.DB.Where("id in (?)", uuid).All(&users)
		if len(users) != 0 && err == nil {
			return users[0], nil
		}
	}

	return models.User{}, errors.New("Unable to find user")
}
