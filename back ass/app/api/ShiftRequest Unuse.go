package api

import (
	"ekans/app/model"
	"net/http"

	"github.com/oscrud/oscrud"
)

// ShiftRequestFilter :
type ShiftRequestFilter struct {
	ID int `json:"ID"`
}

// ShiftRequest :
func (ep *API) ShiftRequest(ctx oscrud.Context) oscrud.Context {
	filter := new(ShiftRequestFilter)
	ctx.BindAll(filter)
	// if want all group
	if filter.ID == 0 {
		var modelArray []model.ShiftRequest
		if err := ep.db.Model(&modelArray).
			Select(); err != nil {
			return ctx.JSON(http.StatusOK, Rbody(modelArray, false, Poem(false)))
		}
		return ctx.JSON(http.StatusOK, Rbody(modelArray, true, Poem(true)))
	}
	// if want 1 people group
	model := new(model.ShiftRequest)
	if err := ep.db.Model(model).
		WhereStruct(filter).
		Select(); err != nil {
		return ctx.JSON(http.StatusOK, Rbody(model, false, Poem(false)))
	}
	return ctx.JSON(http.StatusOK, Rbody(model, true, Poem(true)))
}
