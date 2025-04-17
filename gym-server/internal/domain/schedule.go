package domain

type Schedule struct{
	ID string
	DayOfWeek int
	Hour int
	Minute int
	VideoID string
	Force bool
}