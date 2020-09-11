package model

import (
	"time"

	"github.com/oscrud/oscrud"
	"golang.org/x/crypto/bcrypt"
)

// User :
type User struct {
	ID               int       `json:"ID" oscrud:"$id"`
	Name             string    `json:"Name"`
	Password         string    `json:"Password"`
	PhoneNumber      string    `json:"PhoneNumber"`
	Job              string    `json:"Job"`
	RegisterDateTime time.Time `json:"RegisterDateTime"`
	ResignDateTime   time.Time `json:"ResignDateTime"`
	ResignReason     string    `json:"ResignReason"`
	Picture          string    `json:"Picture"`
	IsResign         bool      `json:"IsResign"`

	Exist       bool   `pg:"-" oscrud:"$Exist"`
	Now         bool   `pg:"-" oscrud:"$Now"`
	JobCategory string `pg:"-" oscrud:"$JobCategory"`
}

// ToCreate :
func (user *User) ToCreate(ctx oscrud.Context) error {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	encodePW := string(hashPassword)
	user.Password = encodePW
	return nil
}

// ToResult :
func (user *User) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	return user, nil
}

// ToQuery :
func (user *User) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	// if user.ID == 0 {
	// 	return nil, nil
	// }
	// if user.Now {
	// 	model := new(User)
	// 	model.ResignDateTime = nil
	// 	return model, nil
	// }

	temp := []string{}
	if user.Now {
		temp = append(temp, "is_resign IS NULL")
	}
	if len(user.JobCategory) >= 1 {
		temp = append(temp, "job IN "+user.JobCategory)
	}
	if len(temp) >= 1 {
		return temp, nil
	}

	return user, nil
	//return "ResignDateTime IS NOT NULL", nil
	//return ["ResignDateTime IS NOT NULL", "ID =0"], nil
}

// ToPatch :
func (user *User) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingUser := incoming.(*User)
	user.ID = incomingUser.ID
	return nil
}

// ToUpdate :
func (user *User) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingUser := incoming.(*User)
	//user.ID = incomingUser.ID
	user.Name = incomingUser.Name
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(incomingUser.Password), bcrypt.MinCost)
	encodePW := string(hashPassword)
	user.Password = encodePW

	user.PhoneNumber = incomingUser.PhoneNumber
	user.Job = incomingUser.Job
	user.RegisterDateTime = incomingUser.RegisterDateTime
	user.ResignDateTime = incomingUser.ResignDateTime
	user.ResignReason = incomingUser.ResignReason
	user.Picture = incomingUser.Picture
	user.IsResign = incomingUser.IsResign
	return nil
}

// ToDelete :
func (user *User) ToDelete(ctx oscrud.Context) error {
	return nil
}
