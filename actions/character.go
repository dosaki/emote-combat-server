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
    userID, plerr := helpers.Param(c, "player_id")
    if plerr != nil {
        return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
    }

    body := getCharacterBody(c)

    tx, ok := c.Value("tx").(*pop.Connection)
    if !ok {
        panic("Unable to get connection")
    }

    playerUUID, puuiderr := uuid.FromString(userID)
    if puuiderr != nil {
        return c.Render(400, r.JSON(map[string]string{"message": "Bad player UUID."}))
    }

    character := models.Character{}
    character.Name = body.Name
    character.PlayerID = playerUUID
    character.Race = body.Race
    character.Gender = body.Gender
    character.IngameName = body.IngameName
    character.Server = body.Server

    var users []models.User
    err := models.DB.Where("id = ?", body.PlayerID).All(&users)
    if err != nil {
        fmt.Println(err)
        return c.Render(404, r.JSON(map[string]string{"message": "Unable to find associated player."}))
    }

    if tx.Create(&character) == nil {

        var skills []models.Skill
        skillsErr := models.DB.All(&skills)

        if skillsErr == nil {
            for _, skill := range skills {
                skillEntry := models.CharacterSheetEntry{
                    CharacterID: character.ID,
                    SkillID:     skill.ID,
                    Value:       skill.StartingValue,
                }
                if tx.Create(&skillEntry) != nil {
                    return c.Render(400, r.JSON(map[string]string{}))
                }
            }
        } else {
            return c.Render(400, r.JSON(map[string]string{}))
        }
        return c.Render(201, r.JSON(character))
    }
    return c.Render(400, r.JSON(map[string]string{}))
}

// CharacterUpdate default implementation.
func CharacterUpdate(c buffalo.Context) error {
    puuid, perr := helpers.Param(c, "id")
    if perr != nil {
        return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
    }

    userID, plerr := helpers.Param(c, "player_id")
    if plerr != nil {
        return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
    }

    var characters []models.Character
    err := models.DB.Where("player_id = ?", userID).Where("id = ?", puuid).All(&characters)
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
    character.Race = body.Race
    character.Gender = body.Gender
    character.IngameName = body.IngameName
    character.Server = body.Server

    var users []models.User
    aperr := models.DB.Where("id = ?", body.PlayerID).All(&users)
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
    puuid, perr := helpers.Param(c, "id")
    if perr != nil {
        return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
    }

    var characters []models.Character
    err := models.DB.Where("id = ?", puuid).All(&characters)
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
        var skillEntries []models.CharacterSheetEntry
        skillsErr := models.DB.Where("character_id = ?", character.ID).All(&skillEntries)
        if skillsErr != nil {
            return c.Render(500, r.JSON(map[string]string{"message": "Something went wrong while finding all the skill entries to delete."}))
        }
        for _, skillEntry := range skillEntries {
            if tx.Destroy(&skillEntry) != nil {
                return c.Render(500, r.JSON(map[string]string{"message": "Something went wrong while deleting all the skill entries."}))
            }
        }
        return c.Render(201, r.JSON(map[string]string{}))
    }

    return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// CharacterList default implementation.
func CharacterList(c buffalo.Context) error {
    puuid, perr := helpers.Param(c, "id")
    var characters []models.Character
    var err error
    var query *pop.Query

    userID, pierr := helpers.Param(c, "player_id")
    if pierr != nil {
        query = models.DB.Where("1=1")
    } else {
        query = models.DB.Where("player_id = ?", userID)
    }

    if perr != nil {
        err = query.All(&characters)
        if err == nil {
            return c.Render(200, r.JSON(characters))
        }
    } else {
        err = query.Where("id = ?", puuid).All(&characters)
        if len(characters) == 0 {
            return c.Render(404, r.JSON(map[string]string{"message": "Character not found."}))
        }
        if err == nil {
            return c.Render(200, r.JSON(characters[0]))
        }
    }
    fmt.Println(err)
    return c.Render(500, r.JSON(map[string]string{"message": "Problem getting character(s)."}))
}

// CharacterList default implementation.
func CharacterListByNameAndServer(c buffalo.Context) error {
    name, nerr := helpers.Param(c, "name")
    server, serr := helpers.Param(c, "server")

    if nerr != nil && serr != nil {
        var characters []models.Character
        var err error
        var query *pop.Query

        err = query.Where("name = ? and server = ?", name, server).All(&characters)
        if len(characters) == 0 {
            return c.Render(404, r.JSON(map[string]string{"message": "Character not found."}))
        }
        if err == nil {
            return c.Render(200, r.JSON(characters[0]))
        }
    }

    fmt.Println(nerr, serr)
    return c.Render(500, r.JSON(map[string]string{"message": "Problem getting character."}))
}
