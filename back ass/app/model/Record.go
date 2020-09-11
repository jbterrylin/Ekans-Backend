package model

import "time"

// Record :
type Record struct {
	ID         int       `json:"ID"`
	NurseID    int       `json:"NursesID"`
	DateTime   time.Time `json:"DateTime"`
	TotalShift int       `json:"TotalShift"`
	LeaveShift int       `json:"LeaveShift"`
	McShift    int       `json:"McShift"`
}
