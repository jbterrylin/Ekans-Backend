package model

import "github.com/oscrud/oscrud"

// Area :
type Area struct {
	ID       int    `json:"ID"`
	AreaName string `json:"AreaName"`
}

// ToCreate :
func (area *Area) ToCreate(ctx oscrud.Context) error {
	return nil
}

// ToResult :
func (area *Area) ToResult(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	return area, nil
}

// ToQuery :
func (area *Area) ToQuery(ctx oscrud.Context, action oscrud.ServiceAction) (interface{}, error) {
	return area, nil
}

// ToPatch :
func (area *Area) ToPatch(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	return nil
}

// ToUpdate :
func (area *Area) ToUpdate(ctx oscrud.Context, incoming oscrud.ServiceModel) error {
	return nil
}

// ToDelete :
func (area *Area) ToDelete(ctx oscrud.Context) error {
	return nil
}
