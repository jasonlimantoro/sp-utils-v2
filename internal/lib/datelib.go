package lib

import "time"

func GetWeekday(today time.Time, pastDeltaWeek int, targetWeekday time.Weekday) time.Time {
	currentWeekday := today.Weekday()
	distanceFromTargetWeekday := int(currentWeekday-targetWeekday) + 7*pastDeltaWeek
	return today.AddDate(0, 0, -distanceFromTargetWeekday)
}
