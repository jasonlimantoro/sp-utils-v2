package createdraft

import "time"

func getFridayDate(today time.Time, deltaWeek int) time.Time {
	currentWeekday := today.Weekday()
	distanceFromFriday := int(currentWeekday-time.Friday) + 7*deltaWeek
	return today.AddDate(0, 0, -distanceFromFriday)
}
