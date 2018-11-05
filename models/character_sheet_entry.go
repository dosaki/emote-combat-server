package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type CharacterSheetEntry struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CharacterID uuid.UUID `json:"character_id" db:"character_id"`
	SkillID     uuid.UUID `json:"skill_id" db:"skill_id"`
	Value       int       `json:"value" db:"value"`
	Note        string    `json:"note" db:"note"`
}

// String is not required by pop and may be deleted
func (c CharacterSheetEntry) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CharacterSheetEntries is not required by pop and may be deleted
type CharacterSheetEntries []CharacterSheetEntry

// String is not required by pop and may be deleted
func (c CharacterSheetEntries) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *CharacterSheetEntry) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: c.Value, Name: "Value"},
		&validators.StringIsPresent{Field: c.Note, Name: "Note"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *CharacterSheetEntry) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *CharacterSheetEntry) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
