package model

import (
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/oscrud/oscrud"
)

// Group :
type Group struct {
	ID          int    `json:"ID"`
	GroupName   string `json:"GroupName"`
	LeaderID    int    `json:"LeaderID"`
	TeammatesID []int  `json:"TeammatesID"`
	IsDeleted   bool   `json:"IsDeleted"`

	IDAddName   []string `pg:"-" json:"IDAddName"`
	Details     []User   `pg:"-" json:"Details"`
	Now         bool     `pg:"-" json:"Now" oscrud:"$Now"`
	SearchValue string   `pg:"-" oscrud:"SearchValue"`
	TargetUser  int      `pg:"-" oscrud:"TargetUser"`
}

// ToCreate :
func (group *Group) ToCreate(ctx oscrud.Context) error {
	return nil
}

// ToResult :
func (group *Group) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	if action == oscrud.ServiceActionFind || action == oscrud.ServiceActionGet {
		db := ctx.GetState("db").(*pg.DB)
		// relation to user table
		model := new(User)
		temp := make([]User, 0)

		db.Model(model).Where("id = ?", group.LeaderID).Select()
		if model.ID == 0 {
			model.ID = group.LeaderID
			model.Exist = false
			group.IDAddName = append(group.IDAddName, strconv.Itoa(model.ID)+" - Not Exist")
		} else {
			group.IDAddName = append(group.IDAddName, strconv.Itoa(model.ID)+" - "+model.Name)
			model.Exist = true
		}
		temp = append(temp, *model)

		for _, v := range group.TeammatesID {
			db.Model(model).Where("id = ?", v).Select()
			if model.ID == 0 {
				model.ID = group.LeaderID
				model.Exist = false
				group.IDAddName = append(group.IDAddName, strconv.Itoa(model.ID)+" - Not Exist")
			} else {
				model.Exist = true
				group.IDAddName = append(group.IDAddName, strconv.Itoa(model.ID)+" - "+model.Name)
			}
			temp = append(temp, *model)
		}

		group.Details = temp
	}
	return group, nil
}

// ToQuery :
func (group *Group) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	temp := []string{}
	if group.SearchValue != "" {
		temp = append(temp, "group_name LIKE '%"+group.SearchValue+"%'")
	}
	if group.TargetUser != 0 {
		temp = append(temp, "(leader_id = "+strconv.Itoa(group.TargetUser)+") OR (\"teammates_id\"::jsonb @> '"+strconv.Itoa(group.TargetUser)+"')")
	}
	if group.Now {
		temp = append(temp, "is_deleted IS NULL")
	}
	if len(temp) > 0 {
		return temp, nil
	}
	return group, nil
}

// ToPatch :
func (group *Group) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingGroup := incoming.(*Group)
	group.ID = incomingGroup.ID
	return nil
}

// ToUpdate :
func (group *Group) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	incomingGroup := incoming.(*Group)
	if incomingGroup.IsDeleted {
		group.IsDeleted = incomingGroup.IsDeleted
		return nil
	}

	group.ID = incomingGroup.ID
	group.GroupName = incomingGroup.GroupName
	group.LeaderID = incomingGroup.LeaderID
	group.TeammatesID = incomingGroup.TeammatesID
	group.IsDeleted = incomingGroup.IsDeleted
	return nil
}

// ToDelete :
func (group *Group) ToDelete(ctx oscrud.Context) error {
	return nil
}
