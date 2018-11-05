package actions

import (
	"encoding/json"

	"github.com/dosaki/owl_power_server/helpers"
	"github.com/dosaki/owl_power_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

func getPlayerBody(c buffalo.Context) models.Player {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := models.Player{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

// PlayerCreate default implementation.
func PlayerCreate(c buffalo.Context) error {
	body := getPlayerBody(c)

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	player := models.Player{}
	player.Name = body.Name

	if tx.Create(&player) == nil {
		return c.Render(200, r.JSON(player))
	}
	return c.Render(400, r.JSON(map[string]string{}))
}

// PlayerUpdate default implementation.
func PlayerUpdate(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	players := []models.Player{}
	err := models.DB.Where("id in (?)", uuid).All(&players)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting player."}))
	}
	if len(players) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Player not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	body := getPlayerBody(c)
	player := players[0]
	player.Name = body.Name

	if tx.Save(&player) == nil {
		return c.Render(200, r.JSON(player))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// PlayerDelete default implementation.
func PlayerDelete(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	players := []models.Player{}
	err := models.DB.Where("id in (?)", uuid).All(&players)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting player."}))
	}
	if len(players) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Player not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	player := players[0]

	if tx.Destroy(&player) == nil {
		return c.Render(201, r.JSON(map[string]string{}))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// PlayerList default implementation.
func PlayerList(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	players := []models.Player{}
	var err error

	if perr != nil {
		err = models.DB.All(&players)
		if err == nil {
			return c.Render(201, r.JSON(players))
		}
	} else {
		err = models.DB.Where("id in (?)", uuid).All(&players)
		if len(players) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Player not found."}))
		}
		if err == nil {
			return c.Render(201, r.JSON(players[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Problem getting player(s)."}))
}
