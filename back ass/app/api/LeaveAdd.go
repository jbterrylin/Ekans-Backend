package api

import (
	"ekans/app/model"
	"net/http"

	"github.com/oscrud/oscrud"
)

// LeaveAdd :
func (ep *API) LeaveAdd(ctx oscrud.Context) oscrud.Context {
	model := new(model.Leave)
	ctx.BindAll(model)

	// if err := ep.db.Model(model).
	// 	Where("id = ?", filter.ID).
	// 	Select(model); err != nil {
	// 	return ctx.JSON(http.StatusOK, Rbody(model, false, Poem(false)))
	// }

	return ctx.JSON(http.StatusOK, Rbody(model, true, Poem(true)))
}
