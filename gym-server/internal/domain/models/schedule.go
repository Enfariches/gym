package models

type Schedule struct {
	ID            int
	DayOfWeek     int
	Hour          int
	Minute        int
	VideoID       string
	DepartamentID string
	Force         bool
}
