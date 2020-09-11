package api

import (
	"ekans/app/model"
	"net/http"
	"time"

	"github.com/oscrud/oscrud"
	"golang.org/x/crypto/bcrypt"
)

// Register :
func (ep *API) Register(ctx oscrud.Context) oscrud.Context {
	model := new(model.User)
	ctx.BindAll(model)

	model.RegisterDateTime = time.Now()
	model.IsResign = false

	var temp string

	if !isInt(model.PhoneNumber) {
		temp = temp + "PhoneNumber,"
	}
	if len(temp) > 0 {
		return ctx.JSON(http.StatusOK, Rbody(temp, false, Poem(false)))
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.MinCost)
	encodePW := string(hashPassword)
	model.Password = encodePW

	if _, err := ep.db.Model(model).Insert(); err != nil {
		return ctx.JSON(http.StatusOK, Rbody(model, false, Poem(false)))
	}

	return ctx.JSON(http.StatusOK, Rbody(model, true, Poem(true)))
}
