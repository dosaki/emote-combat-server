package actions

import (
	"encoding/json"
	"errors"

	"github.com/dosaki/emote_combat_server/helpers"
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	UUID "github.com/gobuffalo/uuid"
)

func getSheetEntryBody(c buffalo.Context) models.CharacterSheetEntry {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := models.CharacterSheetEntry{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

func getSheetEntriesBody(c buffalo.Context) []models.CharacterSheetEntry {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := []models.CharacterSheetEntry{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

func createOne(c buffalo.Context, body models.CharacterSheetEntry, characterID string) (models.CharacterSheetEntry, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	characterUUID, puuiderr := UUID.FromString(characterID)
	if puuiderr != nil {
		return models.CharacterSheetEntry{}, errors.New("bad character uuid")
	}

	sheetEntry := models.CharacterSheetEntry{}
	sheetEntry.CharacterID = characterUUID
	sheetEntry.SkillID = body.SkillID
	sheetEntry.Value = body.Value
	sheetEntry.Note = body.Note

	if tx.Create(&sheetEntry) == nil {
		return sheetEntry, nil
	}
	return models.CharacterSheetEntry{}, errors.New("unknown")
}

func updateOne(c buffalo.Context, body models.CharacterSheetEntry, characterID string, uuid string) (models.CharacterSheetEntry, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	characterUUID, puuiderr := UUID.FromString(characterID)
	if puuiderr != nil {
		return models.CharacterSheetEntry{}, errors.New("Bad character UUID")
	}

	sheetEntries := []models.CharacterSheetEntry{}
	err := models.DB.Where("character_id = ?", characterID).Where("id = ?", uuid).All(&sheetEntries)
	if err != nil {
		return models.CharacterSheetEntry{}, errors.New("Problem getting sheet entry")
	}
	if len(sheetEntries) == 0 {
		return models.CharacterSheetEntry{}, errors.New("Sheet entry not found")
	}

	sheetEntry := sheetEntries[0]
	sheetEntry.CharacterID = characterUUID
	sheetEntry.SkillID = body.SkillID
	sheetEntry.Value = body.Value
	sheetEntry.Note = body.Note

	if tx.Save(&sheetEntry) == nil {
		return sheetEntry, nil
	}
	return sheetEntry, errors.New("unknown")
}

// SheetEntryCreate default implementation.
func SheetEntryCreate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	characters := []models.Character{}
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
	}

	body := getSheetEntryBody(c)
	sheetEntry, err := createOne(c, body, characterID)
	if err == nil {
		return c.Render(201, r.JSON(sheetEntry))
	}
	return c.Render(400, r.JSON(map[string]string{"message": err.Error()}))
}

// SheetEntriesCreate default implementation.
func SheetEntriesCreate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	characters := []models.Character{}
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
	}

	bodies := getSheetEntriesBody(c)
	sheetEntries := []models.CharacterSheetEntry{}
	for _, body := range bodies {
		sheetEntry, err := createOne(c, body, characterID)
		if err != nil {
			return c.Render(400, r.JSON(map[string]string{"message": err.Error()}))
		}
		sheetEntries = append(sheetEntries, sheetEntry)
	}
	return c.Render(200, r.JSON(sheetEntries))
}

// SheetEntryUpdate default implementation.
func SheetEntryUpdate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	characters := []models.Character{}
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
	}

	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	sheetEntry, seError := updateOne(c, getSheetEntryBody(c), characterID, uuid)

	if seError == nil {
		return c.Render(200, r.JSON(sheetEntry))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// SheetEntriesUpdate default implementation.
func SheetEntriesUpdate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	characters := []models.Character{}
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
	}

	bodies := getSheetEntriesBody(c)
	sheetEntries := []models.CharacterSheetEntry{}
	for _, body := range bodies {
		sheetEntry, err := updateOne(c, body, characterID, body.ID.String())
		if err != nil {
			return c.Render(400, r.JSON(map[string]string{"message": err.Error()}))
		}
		sheetEntries = append(sheetEntries, sheetEntry)
	}
	return c.Render(200, r.JSON(sheetEntries))
}

// SheetEntryDelete default implementation.
func SheetEntryDelete(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No player ID provided."}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	characters := []models.Character{}
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
	}

	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	sheetEntries := []models.CharacterSheetEntry{}
	err := models.DB.Where("character_id = ?", characterID).Where("id = ?", uuid).All(&sheetEntries)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting sheet entry."}))
	}
	if len(sheetEntries) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Sheet entry not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	sheetEntry := sheetEntries[0]

	if tx.Destroy(&sheetEntry) == nil {
		return c.Render(201, r.JSON(map[string]string{}))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// SheetEntryList default implementation.
func SheetEntryList(c buffalo.Context) error {
	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No character ID provided."}))
	}

	playerID, pierr := helpers.Param(c, "player_id")
	if pierr == nil {
		characters := []models.Character{}
		cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
		if cerr != nil || len(characters) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Unable to find that player's character."}))
		}
	}

	var query *pop.Query
	var err error
	sheetEntries := []models.CharacterSheetEntry{}
	query = models.DB.Where("character_id = ?", characterID)

	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		err = query.All(&sheetEntries)
		if err == nil {
			return c.Render(200, r.JSON(sheetEntries))
		}
	} else {
		err = query.Where("id = ?", uuid).All(&sheetEntries)
		if len(sheetEntries) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Sheet entry not found."}))
		}
		if err == nil {
			return c.Render(200, r.JSON(sheetEntries[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Problem getting sheet entry(s)."}))
}
