package model

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/oscrud/oscrud"
)

// Leave :
type Leave struct {
	ID         int       `json:"ID"`
	DateTime   time.Time `json:"DateTime"`
	Shift      string    `json:"Shift"`
	NurseID    int       `json:"NurseID"`
	Type       string    `json:"Type"`
	Remark     string    `json:"Remark"`
	ApprovedBy int       `json:"ApprovedBy"`

	Check       bool   `pg:"-" json:"Check"`
	NurseName   string `pg:"-" json:"NurseName"`
	ManagerName string `pg:"-" json:"ManagerName"`
}

// ToCreate :
func (leave *Leave) ToCreate(ctx oscrud.Context) error {
	db := ctx.GetState("db").(*pg.DB)
	db.Model(leave).
		Where("date_time = ?", leave.DateTime).
		Where("shift = ?", leave.Shift).
		Where("nurse_id = ?", leave.NurseID).
		Select()
	if leave.ID == 0 {
		if leave.Check {
			return fmt.Errorf("Not Exist but it is check, so didn't post it")
		}
		return nil
	}
	return fmt.Errorf("%vJBTERRYLIN%v", leave.DateTime, leave.Shift)
}

// ToResult :
func (leave *Leave) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	if action == oscrud.ServiceActionFind || action == oscrud.ServiceActionGet {
		db := ctx.GetState("db").(*pg.DB)
		model := new(User)
		db.Model(model).
			Where("id = ?", leave.NurseID).
			Select()
		leave.NurseName = model.Name

		db.Model(model).
			Where("id = ?", leave.ApprovedBy).
			Select()
		leave.ManagerName = model.Name
	}
	return leave, nil
}

// ToQuery :
func (leave *Leave) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	return leave, nil
}

// ToPatch :
func (leave *Leave) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	return nil
}

// ToUpdate :
func (leave *Leave) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	return nil
}

// ToDelete :
func (leave *Leave) ToDelete(ctx oscrud.Context) error {
	return nil
}
