package gokounta

import (
	"net/http"

	"github.com/mholt/binding"
)

//Shift is the struct for a Kounta Shift
type Shift struct {
	StartedAt  string       `json:"started_at"`
	FinishedAt string       `json:"finished_at"`
	Staff      Staff        `json:"staff_member"`
	Breaks     []ShiftBreak `json:"breaks"`
}

type ShiftBreak struct {
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

//FieldMap is required for binding
func (obj *Shift) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&obj.StartedAt:  "started_at",
		&obj.FinishedAt: "finished_at",
		&obj.Staff:      "staff_member",
		&obj.Breaks:     "breaks",
	}
}

//FieldMap is required for binding
func (obj *ShiftBreak) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&obj.StartedAt:  "started_at",
		&obj.FinishedAt: "finished_at",
	}
}
