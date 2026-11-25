package service

import (
	"slices"
	"time"

	"car-maintenance/internal/model"
)

const (
	StatusOK      = "ok"
	StatusWarning = "warning"
	StatusOverdue = "overdue"

	WarningDaysThreshold = 15
	WarningKMThreshold   = 500
)

// CalculateMaintenanceStatus computes status, next_due_date, days_remaining,
// next_due_mileage, and km_remaining for a VehicleMaintenanceDetail based on now.
func CalculateMaintenanceStatus(detail *model.VehicleMaintenanceDetail, now time.Time) {
	switch detail.ExpirationType {
	case "time":
		detail.Status = calculateTimeStatus(detail, now)
	case "km":
		detail.Status = calculateKMStatus(detail)
	case "both":
		timeStatus := calculateTimeStatus(detail, now)
		kmStatus := calculateKMStatus(detail)
		detail.Status = combineStatuses(timeStatus, kmStatus)
	default:
		detail.Status = StatusOK
	}
}

func calculateTimeStatus(detail *model.VehicleMaintenanceDetail, now time.Time) string {
	if detail.IntervalDays == nil {
		return StatusOK
	}

	lastDate, err := time.Parse("2006-01-02", detail.LastPerformedDate)
	if err != nil {
		return StatusOK
	}

	nextDue := lastDate.AddDate(0, 0, *detail.IntervalDays)
	nextDueStr := nextDue.Format("2006-01-02")
	detail.NextDueDate = &nextDueStr

	// Truncate hours by normalizing both dates to UTC midnight
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	nextDueDate := time.Date(nextDue.Year(), nextDue.Month(), nextDue.Day(), 0, 0, 0, 0, time.UTC)

	daysRem := int(nextDueDate.Sub(nowDate).Hours() / 24)
	detail.DaysRemaining = &daysRem

	if daysRem < 0 {
		return StatusOverdue
	}
	if daysRem <= WarningDaysThreshold {
		return StatusWarning
	}
	return StatusOK
}

func calculateKMStatus(detail *model.VehicleMaintenanceDetail) string {
	if detail.IntervalKM == nil {
		return StatusOK
	}

	nextKM := detail.LastPerformedMileage + *detail.IntervalKM
	detail.NextDueMileage = &nextKM

	kmRem := nextKM - detail.VehicleCurrentMileage
	detail.KMRemaining = &kmRem

	if kmRem < 0 {
		return StatusOverdue
	}
	if kmRem <= WarningKMThreshold {
		return StatusWarning
	}
	return StatusOK
}

func combineStatuses(a, b string) string {
	if a == StatusOverdue || b == StatusOverdue {
		return StatusOverdue
	}
	if a == StatusWarning || b == StatusWarning {
		return StatusWarning
	}
	return StatusOK
}

func statusUrgencyRank(status string) int {
	switch status {
	case StatusOverdue:
		return 0
	case StatusWarning:
		return 1
	case StatusOK:
		return 2
	default:
		return 3
	}
}

// SortMaintenancesByUrgency sorts maintenances with overdue first, then warning, then ok.
func SortMaintenancesByUrgency(items []model.VehicleMaintenanceDetail) {
	slices.SortStableFunc(items, func(a, b model.VehicleMaintenanceDetail) int {
		rankA := statusUrgencyRank(a.Status)
		rankB := statusUrgencyRank(b.Status)
		if rankA != rankB {
			return rankA - rankB
		}
		return int(a.ID - b.ID)
	})
}
