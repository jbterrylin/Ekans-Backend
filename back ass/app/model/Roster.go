package model

import (
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/oscrud/oscrud"
)

// Roster :
type Roster struct {
	ID             int       `json:"ID"`
	DateTime       time.Time `json:"DateTime" oscrud:"DateTime"`
	Shift          string    `json:"Shift" oscrud:"Shift"`
	Area           int       `json:"Area" oscrud:"Area"`
	GroupID        int       `json:"GroupID"`
	LeaderID       int       `json:"LeaderID"`
	NursesID       []int     `json:"NursesID"`
	BUNursesID     []int     `json:"BUNursesID"`
	PTNursesID     []int     `json:"PTNursesID"`
	Attendances    []int     `json:"Attendances"`
	CreateDateTime time.Time `json:"CreateDateTime"`

	GroupName         string   `pg:"-" json:"GroupName"`
	LeaderIDAddName   string   `pg:"-" json:"LeaderIDAddName"`
	NursesIDAddName   []string `pg:"-" json:"NursesIDAddName"`
	BUNursesIDAddName []string `pg:"-" json:"BUNursesIDAddName"`
	PTNursesIDAddName []string `pg:"-" json:"PTNursesIDAddName"`
	Details           []User   `pg:"-" json:"Details"`
}

// ToCreate :
func (roster *Roster) ToCreate(ctx oscrud.Context) error {
	return nil
}

// ToResult :
func (roster *Roster) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	if action == oscrud.ServiceActionFind || action == oscrud.ServiceActionGet {
		db := ctx.GetState("db").(*pg.DB)
		db.Model(roster).
			Where("date_time = ?", roster.DateTime).
			Where("shift = ?", roster.Shift).
			Where("area = ?", roster.Area).
			Select()

		if roster.ID != 0 {
			if roster.GroupID != 0 {
				model1 := new(Group)
				db.Model(model1).Where("id = ?", roster.GroupID).Select()
				roster.GroupName = model1.GroupName
			}

			model2 := new(User)
			db.Model(model2).Where("id = ?", roster.LeaderID).Select()
			roster.LeaderIDAddName = strconv.Itoa(model2.ID) + " - " + model2.Name
			roster.NursesIDAddName = append(roster.NursesIDAddName, strconv.Itoa(model2.ID)+" - "+model2.Name)
			roster.Details = append(roster.Details, *model2)

			for _, v := range roster.NursesID {
				model1 := new(User)
				db.Model(model1).Where("id = ?", v).Select()
				roster.NursesIDAddName = append(roster.NursesIDAddName, strconv.Itoa(model1.ID)+" - "+model1.Name)
				roster.Details = append(roster.Details, *model1)
			}

			for _, v := range roster.BUNursesID {
				model1 := new(User)
				db.Model(model1).Where("id = ?", v).Select()
				roster.BUNursesIDAddName = append(roster.BUNursesIDAddName, strconv.Itoa(model1.ID)+" - "+model1.Name)
				roster.Details = append(roster.Details, *model1)
			}

			for _, v := range roster.PTNursesID {
				model1 := new(User)
				db.Model(model1).Where("id = ?", v).Select()
				roster.PTNursesIDAddName = append(roster.PTNursesIDAddName, strconv.Itoa(model1.ID)+" - "+model1.Name)
				roster.Details = append(roster.Details, *model1)
			}
		}
	}
	return roster, nil
}

// ToQuery :
func (roster *Roster) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	temp := []string{}
	if !roster.DateTime.IsZero() && roster.Shift != "" && roster.Area != 0 {
		temp = append(temp, "date_time = '"+roster.DateTime.Format("2006-01-02 15:04:05")+"+08'")
		temp = append(temp, "shift = '"+roster.Shift+"'")
		temp = append(temp, "area = "+strconv.Itoa(roster.Area))
		return temp, nil
	}
	return roster, nil
}

// ToPatch :
func (roster *Roster) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingRoster := incoming.(*Roster)
	// roster.ID = incomingRoster.ID
	roster.Attendances = incomingRoster.Attendances
	return nil
}

// ToUpdate :
func (roster *Roster) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingRoster := incoming.(*Roster)

	roster.ID = incomingRoster.ID
	roster.DateTime = incomingRoster.DateTime
	roster.Shift = incomingRoster.Shift
	roster.Area = incomingRoster.Area
	roster.GroupID = incomingRoster.GroupID
	roster.LeaderID = incomingRoster.LeaderID
	roster.NursesID = incomingRoster.NursesID
	roster.BUNursesID = incomingRoster.BUNursesID
	roster.PTNursesID = incomingRoster.PTNursesID
	roster.Attendances = incomingRoster.Attendances
	roster.CreateDateTime = incomingRoster.CreateDateTime
	return nil
}

// ToDelete :
func (roster *Roster) ToDelete(ctx oscrud.Context) error {
	return nil
}
