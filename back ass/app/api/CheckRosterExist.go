package api

import (
	"ekans/app/model"
	"net/http"
	"time"

	"github.com/oscrud/oscrud"
)

// CheckRosterExistModel :
type CheckRosterExistModel struct {
	DateTime time.Time `json:"DateTime"`
}

// HaventCreateRoster :
type HaventCreateRoster struct {
	DateTime   time.Time     `json:"DateTime"`
	AreasShift []AreaDeShift `json:"AreasShift"`
}

// AreaDeShift :
type AreaDeShift struct {
	Area  string   `json:"Area"`
	Shift []string `json:"Shift"`
}

// CheckRosterExist :
func (ep *API) CheckRosterExist(ctx oscrud.Context) oscrud.Context {
	filter := new(CheckRosterExistModel)
	ctx.BindAll(filter)

	arealist := make([]model.Area, 0)
	ep.db.Model(&arealist).Select()
	shiftlist := [3]string{"Q1", "Q2", "Q3"}

	rosterArray := make([]model.Roster, 0)

	ep.db.Model(&rosterArray).
		Where("date_time >= ?", filter.DateTime).
		Where("date_time < ?", filter.DateTime.AddDate(0, 0, 15)).
		Select()

	// areadeshift := make([]AreaDeShift, 0)
	total := make([]HaventCreateRoster, 0)
	for i := 0; i < 14; i++ {
		haventcreateroster := new(HaventCreateRoster)
		haventcreateroster.DateTime = filter.DateTime.AddDate(0, 0, i)
		for _, area := range arealist {
			areadata := new(AreaDeShift)
			areadata.Area = area.AreaName
			for _, shift := range shiftlist {
				counter := false
				for _, databaseroster := range rosterArray {
					if databaseroster.DateTime == filter.DateTime.AddDate(0, 0, i) &&
						databaseroster.Area == area.ID &&
						databaseroster.Shift == shift {
						counter = true
					}
				}
				if counter == false {
					areadata.Shift = append(areadata.Shift, shift)
				}
			}
			haventcreateroster.AreasShift = append(haventcreateroster.AreasShift, *areadata)
		}
		total = append(total, *haventcreateroster)
	}

	return ctx.JSON(http.StatusOK, Rbody(total, true, Poem(true)))
}
