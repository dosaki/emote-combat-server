package actions

import (
	"encoding/json"
	"fmt"

	"github.com/dosaki/emote_combat_server/helpers"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

func getCharacterBody(c buffalo.Context) models.Character {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := models.Character{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

// CharacterCreate default implementation.
func CharacterCreate(c buffalo.Context) error {
	playerID, plerr := helpers.Param(c, "player_id")
	if plerr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	body := getCharacterBody(c)

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	playerUUID, puuiderr := uuid.FromString(playerID)
	if puuiderr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "Bad player UUID."}))
	}

	character := models.Character{}
	character.Name = body.Name
	character.PlayerID = playerUUID

	players := []models.Player{}
	err := models.DB.Where("id in (?)", body.PlayerID).All(&players)
	if err != nil {
		fmt.Println(err)
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find associated player."}))
	}

	if tx.Create(&character) == nil {
		return c.Render(200, r.JSON(character))
	}
	return c.Render(400, r.JSON(map[string]string{}))
}

// CharacterUpdate default implementation.
func CharacterUpdate(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	playerID, plerr := helpers.Param(c, "player_id")
	if plerr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characters := []models.Character{}
	err := models.DB.Where("player_id = ?", playerID).Where("id = ?", uuid).All(&characters)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting character."}))
	}
	if len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Character not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	body := getCharacterBody(c)
	character := characters[0]
	character.Name = body.Name
	character.PlayerID = body.PlayerID

	players := []models.Player{}
	aperr := models.DB.Where("id in (?)", body.PlayerID).All(&players)
	if aperr != nil {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find associated player."}))
	}

	if tx.Save(&character) == nil {
		return c.Render(200, r.JSON(character))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// CharacterDelete default implementation.
func CharacterDelete(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	characters := []models.Character{}
	err := models.DB.Where("id in (?)", uuid).All(&characters)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting character."}))
	}
	if len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Character not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	character := characters[0]

	if tx.Destroy(&character) == nil {
		return c.Render(201, r.JSON(map[string]string{}))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// CharacterList default implementation.
func CharacterList(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	characters := []models.Character{}
	var err error
	var query *pop.Query

	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		query = models.DB.Where("1=1")
	} else {
		query = models.DB.Where("player_id = ?", playerID)
	}

	if perr != nil {
		err = query.All(&characters)
		if err == nil {
			return c.Render(201, r.JSON(characters))
		}
	} else {
		err = query.Where("id in (?)", uuid).All(&characters)
		if len(characters) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Character not found."}))
		}
		if err == nil {
			return c.Render(201, r.JSON(characters[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Problem getting character(s)."}))
}
