package model

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/oscrud/oscrud"
)

// ShiftRequest :
type ShiftRequest struct {
	ID                int       `json:"ID"`
	RequesterID       int       `json:"RequesterID"`
	RequesterDateTime time.Time `json:"RequesterDateTime"`
	RequesterShift    string    `json:"RequesterShift"`
	RequesterArea     int       `json:"RequesterArea"`
	RequesteeID       int       `json:"RequesteeID"`
	RequesteeDateTime time.Time `json:"RequesteeDateTime"`
	RequesteeShift    string    `json:"RequesteeShift"`
	RequesteeArea     int       `json:"RequesteeArea"`
	Remark            string    `json:"Remark"`
	Status            string    `json:"Status"`
	DecideDateTime    time.Time `json:"DecideDateTime"`

	StatusUpdate bool `pg:"-" oscrud:"$StatusUpdate"`
}

// ToCreate :
func (shiftRequest *ShiftRequest) ToCreate(ctx oscrud.Context) error {
	return nil
}

// ToResult :
func (shiftRequest *ShiftRequest) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {

	return shiftRequest, nil
}

// ToQuery :
func (shiftRequest *ShiftRequest) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	if action == oscrud.ServiceActionFind || action == oscrud.ServiceActionGet {
		db := ctx.GetState("db").(*pg.DB)
		var modelArray []ShiftRequest

		db.Model(&modelArray).
			Where("requester_id == ?", shiftRequest.ID).
			WhereOr("requestee_id == ?", shiftRequest.ID).
			Select()
		return modelArray, nil
	}
	return shiftRequest, nil
}

// ToPatch :
func (shiftRequest *ShiftRequest) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	return nil
}

// ToUpdate :
func (shiftRequest *ShiftRequest) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingShiftRequest := incoming.(*ShiftRequest)
	if incomingShiftRequest.StatusUpdate {
		shiftRequest.Status = incomingShiftRequest.Status
		if shiftRequest.Status == "Accept" {
			db := ctx.GetState("db").(*pg.DB)
			//requester => requestee
			requesterRoster := new(Roster)
			db.Model(requesterRoster).
				Where("requester_date_time == ?", incomingShiftRequest.RequesterDateTime).
				Where("requester_area == ?", incomingShiftRequest.RequesterArea).
				Where("requester_shift == ?", incomingShiftRequest.RequesterShift).
				Select()
			db.Model(requesterRoster).Update(exchangedata(requesterRoster, incomingShiftRequest.RequesterID, incomingShiftRequest.RequesteeID))

			//requestee => requester
			requesteeRoster := new(Roster)
			db.Model(requesteeRoster).
				Where("requester_date_time == ?", incomingShiftRequest.RequesteeDateTime).
				Where("requester_area == ?", incomingShiftRequest.RequesteeArea).
				Where("requester_shift == ?", incomingShiftRequest.RequesteeShift).
				Select()
			db.Model(requesterRoster).Update(exchangedata(requesteeRoster, incomingShiftRequest.RequesteeID, incomingShiftRequest.RequesterID))
			shiftRequest.DecideDateTime = incomingShiftRequest.DecideDateTime
		}
	} else {
		shiftRequest.ID = incomingShiftRequest.ID
		shiftRequest.RequesterID = incomingShiftRequest.RequesterID
		shiftRequest.RequesterDateTime = incomingShiftRequest.RequesterDateTime
		shiftRequest.RequesterShift = incomingShiftRequest.RequesterShift
		shiftRequest.RequesterArea = incomingShiftRequest.RequesterArea
		shiftRequest.RequesteeID = incomingShiftRequest.RequesteeID
		shiftRequest.RequesteeDateTime = incomingShiftRequest.RequesteeDateTime
		shiftRequest.RequesteeShift = incomingShiftRequest.RequesteeShift
		shiftRequest.RequesteeArea = incomingShiftRequest.RequesteeArea
		shiftRequest.Remark = incomingShiftRequest.Remark
		shiftRequest.Status = incomingShiftRequest.Status
		shiftRequest.DecideDateTime = incomingShiftRequest.DecideDateTime
	}
	return nil
}

func exchangedata(roster *Roster, tartgetID int, changeID int) interface{} {
	if roster.LeaderID == tartgetID {
		roster.LeaderID = tartgetID
		return roster
	}
	for index, nurseID := range roster.NursesID {
		if nurseID == tartgetID {
			roster.NursesID[index] = changeID
			return roster
		}
	}
	for index, nurseID := range roster.BUNursesID {
		if nurseID == tartgetID {
			roster.NursesID[index] = changeID
			return roster
		}
	}
	for index, nurseID := range roster.PTNursesID {
		if nurseID == tartgetID {
			roster.NursesID[index] = changeID
		}
	}
	return roster
}

// ToDelete :
func (shiftRequest *ShiftRequest) ToDelete(ctx oscrud.Context) error {
	return nil
}
