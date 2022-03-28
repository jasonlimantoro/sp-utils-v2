package getweeklyupdates

import "time"

func getMondayDate(today time.Time, deltaWeek int) time.Time {
	currentWeekday := today.Weekday()
	distanceFromMonday := int(currentWeekday-time.Monday) + 7*deltaWeek
	return today.AddDate(0, 0, -distanceFromMonday)
}
