package api

import (
	"ekans/app/model"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/oscrud/oscrud"
)

// CalanderFilter :
type CalanderFilter struct {
	ID          int       `json:"ID"`
	DateTime    time.Time `json:"DateTime"`
	ForExchange bool      `json:"ForExchange"`
}

// LeaveFilter :
type LeaveFilter struct {
	NurseID  int       `json:"NurseID"`
	DateTime time.Time `json:"DateTime"`
	Shift    string    `json:"Shift"`
	Area     int       `json:"Area"`
}

// Result :
type Result struct {
	DayList []Day    `json:"DayList"`
	Meta    MetaData `json:"Meta"`
}

// MetaData :
type MetaData struct {
	Total       int `json:"Total"`
	MC          int `json:"MC"`
	AnnualLeave int `json:"AnnualLeave"`
	Attend      int `json:"Attend"`
}

// Day :
type Day struct {
	Number     int      `json:"Number"`
	MatterList []Matter `json:"Matter"`
}

// Matter :
type Matter struct {
	Shift     string `json:"Shift"`
	Area      int    `json:"Area"`
	Situation string `json:"Situation"`
}

// Calander :
func (ep *API) Calander(ctx oscrud.Context) oscrud.Context {
	filter := new(CalanderFilter)
	ctx.BindAll(filter)

	month := [2]string{}
	for i := 0; i < 2; i++ {
		temp := int(filter.DateTime.Month()) + i
		if temp < 10 {
			month[i] = "0" + strconv.Itoa(temp)
		} else {
			month[i] = strconv.Itoa(temp)
		}
	}

	// get that month roster (Attend/ MC/ Leave)
	rosterArray := make([]model.Roster, 0)
	if filter.ForExchange == true {
		ep.db.Model(&rosterArray).
			// Where("date_time == ?", strconv.Itoa(filter.DateTime.Year())+"-"+month[0]+"-"+strconv.Itoa(filter.DateTime.Day())).
			// WhereStruct(filter).
			Where("date_time = ?", filter.DateTime).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr("leader_id = ?", filter.ID).
					WhereOr("nurses_id::jsonb @> '?'", filter.ID).
					WhereOr("bu_nurses_id::jsonb @> '?'", filter.ID).
					WhereOr("pt_nurses_id::jsonb @> '?'", filter.ID)
				return q, nil
			}).
			Select()
	} else {
		ep.db.Model(&rosterArray).
			Where("date_time >= ?", strconv.Itoa(filter.DateTime.Year())+"-"+month[0]+"-01").
			Where("date_time < ?", strconv.Itoa(filter.DateTime.Year())+"-"+month[1]+"-01").
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr("leader_id = ?", filter.ID).
					WhereOr("nurses_id::jsonb @> '?'", filter.ID).
					WhereOr("bu_nurses_id::jsonb @> '?'", filter.ID).
					WhereOr("pt_nurses_id::jsonb @> '?'", filter.ID)
				return q, nil
			}).
			Select()
	}

	meta := new(MetaData)
	worklist := make([]Day, 0)
	if len(rosterArray) > 0 {
		for _, v := range rosterArray {
			// get that day status (Attend/ MC/ Leave)
			matter := new(Matter)
			matter.Area = v.Area
			matter.Shift = v.Shift
			for _, attandance := range v.Attendances {
				if attandance == filter.ID {
					matter.Situation = "Attend"
					meta.Attend = meta.Attend + 1
				}
			}
			if matter.Situation == "" {
				leave := new(model.Leave)
				leavefilter := new(LeaveFilter)
				leavefilter.NurseID = filter.ID
				leavefilter.DateTime = v.DateTime
				leavefilter.Shift = v.Shift
				ep.db.Model(leave).WhereStruct(leavefilter).Select()
				if leave.ID != 0 {
					matter.Situation = leave.Type
					if leave.Type == "MC" {
						meta.MC = meta.MC + 1
					} else {
						meta.AnnualLeave = meta.AnnualLeave + 1
					}
				}
			}
			if matter.Situation == "" {
				matter.Situation = "Not Come"
			}
			if filter.ForExchange == false || (filter.ForExchange == true && matter.Situation == "Not Come") {
				// put to worklist
				counter := false
				for index, work := range worklist {
					if work.Number == v.DateTime.Day() {
						worklist[index].MatterList = append(worklist[index].MatterList, *matter)
						counter = true
					}
				}
				if counter == false {
					day := new(Day)
					day.Number = v.DateTime.Day()
					day.MatterList = append(day.MatterList, *matter)
					worklist = append(worklist, *day)
				}
			}
		}
	}
	meta.Total = meta.AnnualLeave + meta.Attend + meta.MC
	result := new(Result)
	result.DayList = worklist
	result.Meta = *meta
	return ctx.JSON(http.StatusOK, Rbody(result, true, Poem(false)))
}
