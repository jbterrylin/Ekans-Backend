package api

import (
	"ekans/app/model"
	"fmt"
	"net/http"

	"github.com/oscrud/oscrud"
	"golang.org/x/crypto/bcrypt"
)

// LoginUserFilter :
type LoginUserFilter struct {
	ID       int    `json:"ID"`
	Password string `json:"Password"`
}

// Login :
func (ep *API) Login(ctx oscrud.Context) oscrud.Context {
	model := new(model.User)
	filter := new(LoginUserFilter)
	ctx.BindAll(filter)

	if err := ep.db.Model(model).
		Where("id = ?", filter.ID).
		Select(model); err != nil {
		return ctx.JSON(http.StatusOK, Rbody(model, false, Poem(false)))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(filter.Password)); err != nil {
		fmt.Println("pw wrong")
		return ctx.JSON(http.StatusOK, Rbody(err, false, Poem(false)))
	}

	model.Password = "Secret"

	return ctx.JSON(http.StatusOK, Rbody(model, true, Poem(true)))
}
