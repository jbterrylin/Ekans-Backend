package api

import (
	"ekans/app/model"
	"net/http"

	"github.com/oscrud/oscrud"
)

// RosterRosterFilter :
type RosterRosterFilter struct {
	ID int `json:"ID"`
}

// Roster :
func (ep *API) Roster(ctx oscrud.Context) oscrud.Context {
	filter := new(RosterRosterFilter)
	ctx.BindAll(filter)
	// if want all Roster
	if filter.ID == 0 {
		var modelArray []model.Roster
		if err := ep.db.Model(&modelArray).
			Select(); err != nil {
			return ctx.JSON(http.StatusOK, Rbody(modelArray, false, Poem(false)))
		}
		return ctx.JSON(http.StatusOK, Rbody(modelArray, true, Poem(true)))
	}
	// if want 1 people Roster
	model := new(model.Roster)
	if err := ep.db.Model(model).
		WhereStruct(filter).
		Select(); err != nil {
		return ctx.JSON(http.StatusOK, Rbody(model, false, Poem(false)))
	}
	return ctx.JSON(http.StatusOK, Rbody(model, true, Poem(true)))
}
