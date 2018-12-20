package actions

import (
	"encoding/json"
	"errors"
	"github.com/dosaki/emote_combat_server/messages"

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
	var body []models.CharacterSheetEntry
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

func createOne(c buffalo.Context, body models.CharacterSheetEntry, characterID string) (models.CharacterSheetEntry, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic(messages.NoConnectionError)
	}

	characterUUID, puuiderr := UUID.FromString(characterID)
	if puuiderr != nil {
		return models.CharacterSheetEntry{}, errors.New(messages.BadUUIDError)
	}

	sheetEntry := models.CharacterSheetEntry{}
	sheetEntry.CharacterID = characterUUID
	sheetEntry.SkillID = body.SkillID
	sheetEntry.Value = body.Value
	sheetEntry.Note = body.Note

	if tx.Create(&sheetEntry) == nil {
		return sheetEntry, nil
	}
	return models.CharacterSheetEntry{}, errors.New(messages.UnknownError)
}

func updateOne(c buffalo.Context, body models.CharacterSheetEntry, characterID string, uuid string) (models.CharacterSheetEntry, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic(messages.NoConnectionError)
	}

	characterUUID, puuiderr := UUID.FromString(characterID)
	if puuiderr != nil {
		return models.CharacterSheetEntry{}, errors.New(messages.BadUUIDError)
	}

	var sheetEntries []models.CharacterSheetEntry
	err := models.DB.Where("character_id = ?", characterID).Where("id = ?", uuid).All(&sheetEntries)
	if err != nil {
		return models.CharacterSheetEntry{}, errors.New(messages.ProblemGettingSheetEntryError)
	}
	if len(sheetEntries) == 0 {
		return models.CharacterSheetEntry{}, errors.New(messages.SheetNotFoundError)
	}

	sheetEntry := sheetEntries[0]
	sheetEntry.CharacterID = characterUUID
	sheetEntry.SkillID = body.SkillID
	sheetEntry.Value = body.Value
	sheetEntry.Note = body.Note

	if tx.Save(&sheetEntry) == nil {
		return sheetEntry, nil
	}
	return sheetEntry, errors.New(messages.UnknownError)
}

// SheetEntryCreate default implementation.
func SheetEntryCreate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoPlayerIDError}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	var characters []models.Character
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
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
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoPlayerIDError}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	var characters []models.Character
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
	}

	bodies := getSheetEntriesBody(c)
	var sheetEntries []models.CharacterSheetEntry
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
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoPlayerIDError}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	var characters []models.Character
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
	}

	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoSheetIDError}))
	}

	sheetEntry, seError := updateOne(c, getSheetEntryBody(c), characterID, uuid)

	if seError == nil {
		return c.Render(200, r.JSON(sheetEntry))
	}

	return c.Render(500, r.JSON(map[string]string{"message": messages.UnknownError}))
}

// SheetEntriesUpdate default implementation.
func SheetEntriesUpdate(c buffalo.Context) error {
	playerID, pierr := helpers.Param(c, "player_id")
	if pierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoPlayerIDError}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	var characters []models.Character
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
	}

	bodies := getSheetEntriesBody(c)
	var sheetEntries []models.CharacterSheetEntry
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
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoPlayerIDError}))
	}

	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	var characters []models.Character
	cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
	if cerr != nil || len(characters) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
	}

	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoSheetIDError}))
	}

	var sheetEntries []models.CharacterSheetEntry
	err := models.DB.Where("character_id = ?", characterID).Where("id = ?", uuid).All(&sheetEntries)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": messages.ProblemGettingSheetEntryError}))
	}
	if len(sheetEntries) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": messages.SheetNotFoundError}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic(messages.NoConnectionError)
	}

	sheetEntry := sheetEntries[0]

	if tx.Destroy(&sheetEntry) == nil {
		return c.Render(201, r.JSON(map[string]string{}))
	}

	return c.Render(500, r.JSON(map[string]string{"message": messages.UnknownError}))
}

// SheetEntryList default implementation.
func SheetEntryList(c buffalo.Context) error {
	characterID, cierr := helpers.Param(c, "character_id")
	if cierr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": messages.NoCharacterIDError}))
	}

	playerID, pierr := helpers.Param(c, "player_id")
	if pierr == nil {
		var characters []models.Character
		cerr := models.DB.Where("player_id = ?", playerID).Where("id = ?", characterID).All(&characters)
		if cerr != nil || len(characters) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": messages.PlayerCharacterNotFoundError}))
		}
	}

	var query *pop.Query
	var err error
	var sheetEntries []models.CharacterSheetEntry
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
			return c.Render(404, r.JSON(map[string]string{"message": messages.SheetNotFoundError}))
		}
		if err == nil {
			return c.Render(200, r.JSON(sheetEntries[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": messages.ProblemGettingSheetEntryError}))
}
