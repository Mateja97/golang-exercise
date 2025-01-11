package models

type Event struct {
	State     State
	CompanyID string
	Company   *Company
}

type State int8

const (
	StateInvalid State = iota
	StateCreated
	StateUpdated
	StateDeleted
)
