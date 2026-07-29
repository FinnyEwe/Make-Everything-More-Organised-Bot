package model

type Reminder struct {
	ID uint
	Description string
	Date string
	Discord bool
	Gcal bool
}