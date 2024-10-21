package main

import (
	"fmt"
	"math"
	"time"
)

const totalAnnualLeave = 14

// Function to calculate leave eligibility
func canTakeLeave(joinDate string, plannedLeaveDate string, publicHolidays, requestedDays int) (bool, string) {
	join, err1 := time.Parse("2006-01-02", joinDate)
	leave, err2 := time.Parse("2006-01-02", plannedLeaveDate)
	if err1 != nil || err2 != nil {
		return false, "Invalid date format. Use YYYY-MM-DD."
	}

	personalLeave := totalAnnualLeave - publicHolidays
	probationEnd := join.AddDate(0, 0, 180)

	// Check if the planned leave is within the first 180 days
	if leave.Before(probationEnd) {
		return false, "Cannot take personal leave within the first 180 days."
	}

	endOfYear := time.Date(probationEnd.Year(), 12, 31, 0, 0, 0, 0, probationEnd.Location())
	daysRemaining := endOfYear.Sub(probationEnd).Hours() / 24

	proRatedLeave := int(math.Floor(daysRemaining / 365 * float64(personalLeave)))

	// Check if the requested leave exceeds the pro-rated entitlement
	if requestedDays > proRatedLeave {
		return false, fmt.Sprintf("Only %d day(s) of personal leave available for the rest of the year.", proRatedLeave)
	}

	// Check if the requested leave exceeds 3 consecutive days
	if requestedDays > 3 {
		return false, "Cannot take more than 3 consecutive days of personal leave."
	}

	// If all conditions are met, return true
	return true, "Leave request approved."
}

func main() {
	// Example 1: False, within 180 days
	publicHolidays := 7
	joinDate := "2021-05-01"
	plannedLeaveDate := "2021-07-05"
	requestedDays := 1
	canTake, reason := canTakeLeave(joinDate, plannedLeaveDate, publicHolidays, requestedDays)
	fmt.Println(canTake, reason) // Output: false, "Cannot take personal leave within the first 180 days."

	// Example 2: False, only 1 day allowed
	joinDate = "2021-05-01"
	plannedLeaveDate = "2021-11-05"
	requestedDays = 3
	canTake, reason = canTakeLeave(joinDate, plannedLeaveDate, publicHolidays, requestedDays)
	fmt.Println(canTake, reason) // Output: false, "Only 1 day(s) of personal leave available for the rest of the year."

	// Example 3: True, allowed
	joinDate = "2021-01-05"
	plannedLeaveDate = "2021-12-18"
	requestedDays = 1
	canTake, reason = canTakeLeave(joinDate, plannedLeaveDate, publicHolidays, requestedDays)
	fmt.Println(canTake, reason) // Output: true, "Leave request approved."
}
