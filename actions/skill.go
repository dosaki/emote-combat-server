package actions

import (
	"encoding/json"

	"github.com/dosaki/owl_power_server/helpers"
	"github.com/dosaki/owl_power_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

func getSkillBody(c buffalo.Context) models.Skill {
	request := c.Request()
	decoder := json.NewDecoder(request.Body)
	body := models.Skill{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

// SkillCreate default implementation.
func SkillCreate(c buffalo.Context) error {
	body := getSkillBody(c)

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	skill := models.Skill{}
	skill.Name = body.Name
	skill.Description = body.Description
	skill.ParentSkillID = body.ParentSkillID
	skill.Cost = body.Cost

	if tx.Create(&skill) == nil {
		return c.Render(200, r.JSON(skill))
	}
	return c.Render(400, r.JSON(map[string]string{}))
}

// SkillUpdate default implementation.
func SkillUpdate(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	skills := []models.Skill{}
	err := models.DB.Where("id in (?)", uuid).All(&skills)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting skill."}))
	}
	if len(skills) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Skill not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	body := getSkillBody(c)
	skill := skills[0]
	skill.Name = body.Name
	skill.Description = body.Description
	skill.ParentSkillID = body.ParentSkillID
	skill.Cost = body.Cost

	if tx.Save(&skill) == nil {
		return c.Render(200, r.JSON(skill))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// SkillDelete default implementation.
func SkillDelete(c buffalo.Context) error {
	uuid, perr := helpers.Param(c, "id")
	if perr != nil {
		return c.Render(400, r.JSON(map[string]string{"message": "No ID provided."}))
	}

	skills := []models.Skill{}
	err := models.DB.Where("id in (?)", uuid).All(&skills)
	if err != nil {
		return c.Render(500, r.JSON(map[string]string{"message": "Problem getting skill."}))
	}
	if len(skills) == 0 {
		return c.Render(404, r.JSON(map[string]string{"message": "Skill not found."}))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		panic("Unable to get connection")
	}

	skill := skills[0]

	if tx.Destroy(&skill) == nil {
		return c.Render(201, r.JSON(map[string]string{}))
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Unknown error."}))
}

// SkillList default implementation.
func SkillList(c buffalo.Context) error {
	skills := []models.Skill{}
	var err error
	var query *pop.Query

	parentID, pierr := helpers.Param(c, "parent_id")
	if pierr != nil {
		query = models.DB.Where("1=1")
	} else {
		query = models.DB.Where("parent_skill_id = ?", parentID)
	}

	uuid, perr := helpers.Param(c, "id")

	if perr != nil {
		err = query.All(&skills)
		if err == nil {
			return c.Render(201, r.JSON(skills))
		}
	} else {
		err = query.Where("id in (?)", uuid).All(&skills)
		if len(skills) == 0 {
			return c.Render(404, r.JSON(map[string]string{"message": "Skill(s) not found."}))
		}
		if err == nil {
			return c.Render(201, r.JSON(skills[0]))
		}
	}

	return c.Render(500, r.JSON(map[string]string{"message": "Problem getting skill(s)."}))
}
