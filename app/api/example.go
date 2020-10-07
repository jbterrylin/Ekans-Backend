package api

import (
	"ekans/app/model"
	"net/http"
	"time"

	"github.com/oscrud/oscrud"
)

// GetExample :
func (ep *API) GetExample(ctx oscrud.Context) oscrud.Context {
	model1 := new(model.User)
	model1.Name = "Lee Wang Lin"
	model1.Password = "123"
	model1.PhoneNumber = "0123456789"
	model1.Job = "Admin"
	model1.RegisterDateTime = time.Now()
	a, err := ep.db.Model(model1).Insert()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err)
	}
	return ctx.JSON(http.StatusOK, a)
}

// GetUsers :
func (ep *API) GetUsers(ctx oscrud.Context) oscrud.Context {
	model := new(model.User)
	if err := ep.db.Model(model).Select(); err != nil {
		return ctx.Error(69, err)
	}
	return ctx.JSON(http.StatusOK, model)
}
