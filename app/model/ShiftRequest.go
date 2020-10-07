package model

import (
	"errors"
	"log"
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

	RequesteeName string `pg:"-" json:"RequesteeName"`
	RequesterName string `pg:"-" json:"RequesterName"`
	StatusUpdate  bool   `pg:"-" oscrud:"$StatusUpdate"`
	OtherName     string `pg:"-" json:"OtherName"`
}

// ToCreate :
func (shiftRequest *ShiftRequest) ToCreate(ctx oscrud.Context) error {
	db := ctx.GetState("db").(*pg.DB)
	requesterRoster := new(Roster)
	db.Model(requesterRoster).
		Where("date_time = ?", shiftRequest.RequesterDateTime).
		Where("area = ?", shiftRequest.RequesterArea).
		Where("shift = ?", shiftRequest.RequesterShift).
		Select()
	log.Println(requesterRoster)
	log.Println(shiftRequest.RequesteeID)
	if existinroster(requesterRoster, shiftRequest.RequesteeID) {
		return errors.New("Please check roster list")
	}

	requesteeRoster := new(Roster)
	db.Model(requesteeRoster).
		Where("date_time = ?", shiftRequest.RequesteeDateTime).
		Where("area = ?", shiftRequest.RequesteeArea).
		Where("shift = ?", shiftRequest.RequesteeShift).
		Select()
	if existinroster(requesteeRoster, shiftRequest.RequesterID) {
		return errors.New("Please check roster list")
	}
	return nil
}

// ToResult :
func (shiftRequest *ShiftRequest) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	return shiftRequest, nil
}

// ToQuery :
func (shiftRequest *ShiftRequest) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	if action == oscrud.ServiceActionFind || action == oscrud.ServiceActionGet {
		// db := ctx.GetState("db").(*pg.DB)

		// var temp1 []ShiftRequest
		// db.Model(&temp1).
		// 	Where("requester_id = ?", shiftRequest.ID).
		// WhereOr("requestee_id == ?", shiftRequest.ID).
		// 	Select()
		// log.Println(len(temp1))

		// var temp2 []ShiftRequest
		// db.Model(&temp2).
		// 	// Where("requester_id = ?", shiftRequest.ID).
		// 	Where("requestee_id = ?", shiftRequest.ID).
		// 	Select()
		// log.Println(len(temp2))

		// modelArray := append([]ShiftRequest{}, append(temp1, temp2...)...)

		// return modelArray, nil
		// log.Println(len(modelArray))

		// log.Println("hello")
		// if len(modelArray) >= 1 {
		// 	for index, shiftrequest := range modelArray {
		// 		log.Println(shiftrequest.RequesteeID)

		// 		modelUser := new(User)
		// 		err := db.Model(modelUser).Where("id = ?", shiftRequest.RequesteeID).Select()
		// 		log.Println(modelUser.ID)
		// 		log.Println(err.Error())
		// 		modelArray[index].RequesteeName = modelUser.Name

		// 		log.Println(modelArray[index].RequesteeName)

		// 		// db.Model(modelUser).
		// 		// 	Where("id = ?", shiftRequest.RequesterID).
		// 		// 	Select()
		// 		// modelArray[index].RequesterName = modelUser.Name
		// 		// log.Println(modelArray[index].RequesterName)
		// 	}
		// }
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
	db := ctx.GetState("db").(*pg.DB)
	if incomingShiftRequest.StatusUpdate {
		shiftRequest.Status = incomingShiftRequest.Status
		if shiftRequest.Status == "Accept" {
			//requester => requestee
			requesterRoster := new(Roster)
			db.Model(requesterRoster).
				Where("date_time = ?", incomingShiftRequest.RequesterDateTime).
				Where("area = ?", incomingShiftRequest.RequesterArea).
				Where("shift = ?", incomingShiftRequest.RequesterShift).
				Select()
			db.Model(exchangedata(requesterRoster, incomingShiftRequest.RequesterID, incomingShiftRequest.RequesteeID)).Where("id = ?id").Update()

			//requestee => requester
			requesteeRoster := new(Roster)
			db.Model(requesteeRoster).
				Where("date_time = ?", incomingShiftRequest.RequesteeDateTime).
				Where("area = ?", incomingShiftRequest.RequesteeArea).
				Where("shift = ?", incomingShiftRequest.RequesteeShift).
				Select()
			db.Model(exchangedata(requesteeRoster, incomingShiftRequest.RequesteeID, incomingShiftRequest.RequesterID)).Where("id = ?id").Update()
			shiftRequest.DecideDateTime = incomingShiftRequest.DecideDateTime
		} else {
			shiftRequest.Status = incomingShiftRequest.Status
			shiftRequest.DecideDateTime = incomingShiftRequest.DecideDateTime
		}
	} else {
		requesterRoster := new(Roster)
		db.Model(requesterRoster).
			Where("date_time = ?", shiftRequest.RequesterDateTime).
			Where("area = ?", shiftRequest.RequesterArea).
			Where("shift = ?", shiftRequest.RequesterShift).
			Select()
		if existinroster(requesterRoster, shiftRequest.RequesteeID) {
			return errors.New("Please check roster list")
		}

		requesteeRoster := new(Roster)
		db.Model(requesteeRoster).
			Where("date_time = ?", shiftRequest.RequesteeDateTime).
			Where("area = ?", shiftRequest.RequesteeArea).
			Where("shift = ?", shiftRequest.RequesteeShift).
			Select()
		if existinroster(requesteeRoster, shiftRequest.RequesterID) {
			return errors.New("Please check roster list")
		}

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

func existinroster(roster *Roster, tartgetID int) bool {
	if roster.LeaderID == tartgetID {
		return true
	}
	for _, nurseID := range roster.NursesID {
		if nurseID == tartgetID {
			return true
		}
	}
	for _, nurseID := range roster.BUNursesID {
		if nurseID == tartgetID {
			return true
		}
	}
	for _, nurseID := range roster.PTNursesID {
		if nurseID == tartgetID {
			return true
		}
	}
	return false
}

func exchangedata(roster *Roster, tartgetID int, changeID int) interface{} {
	if roster.LeaderID == tartgetID {
		roster.LeaderID = changeID
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
