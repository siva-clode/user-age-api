package service

import "time"

// CalculateAge gives age in whole years.
func CalculateAge(dob time.Time, at time.Time) int {
	// normalize times to date-only by using UTC and same location
	y1, m1, d1 := dob.Date()
	y2, m2, d2 := at.Date()

	age := y2 - y1
	if m2 < m1 || (m2 == m1 && d2 < d1) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}
