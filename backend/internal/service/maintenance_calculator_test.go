package service_test

import (
	"testing"
	"time"

	"car-maintenance/internal/model"
	"car-maintenance/internal/service"
)

func intPtr(v int) *int {
	return &v
}

func TestCalculateMaintenanceStatus_Time(t *testing.T) {
	now := time.Date(2026, 6, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		lastDate       string
		intervalDays   int
		wantStatus     string
		wantNextDate   string
		wantDaysRemain int
	}{
		{
			name:           "ok - more than 15 days remaining",
			lastDate:       "2026-06-01",
			intervalDays:   30, // due 2026-07-01 -> 16 days remaining
			wantStatus:     service.StatusOK,
			wantNextDate:   "2026-07-01",
			wantDaysRemain: 16,
		},
		{
			name:           "warning - exactly 15 days remaining",
			lastDate:       "2026-06-01",
			intervalDays:   29, // due 2026-06-30 -> 15 days remaining
			wantStatus:     service.StatusWarning,
			wantNextDate:   "2026-06-30",
			wantDaysRemain: 15,
		},
		{
			name:           "warning - due today (0 days remaining)",
			lastDate:       "2026-05-16",
			intervalDays:   30, // due 2026-06-15 -> 0 days remaining
			wantStatus:     service.StatusWarning,
			wantNextDate:   "2026-06-15",
			wantDaysRemain: 0,
		},
		{
			name:           "overdue - 1 day past due",
			lastDate:       "2026-05-15",
			intervalDays:   30, // due 2026-06-14 -> -1 days remaining
			wantStatus:     service.StatusOverdue,
			wantNextDate:   "2026-06-14",
			wantDaysRemain: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			detail := model.VehicleMaintenanceDetail{
				ExpirationType:    "time",
				IntervalDays:      intPtr(tc.intervalDays),
				LastPerformedDate: tc.lastDate,
			}

			service.CalculateMaintenanceStatus(&detail, now)

			if detail.Status != tc.wantStatus {
				t.Errorf("status = %q, want %q", detail.Status, tc.wantStatus)
			}
			if detail.NextDueDate == nil || *detail.NextDueDate != tc.wantNextDate {
				t.Errorf("next_due_date = %v, want %q", detail.NextDueDate, tc.wantNextDate)
			}
			if detail.DaysRemaining == nil || *detail.DaysRemaining != tc.wantDaysRemain {
				t.Errorf("days_remaining = %v, want %d", detail.DaysRemaining, tc.wantDaysRemain)
			}
			if detail.NextDueMileage != nil || detail.KMRemaining != nil {
				t.Errorf("expected nil mileage fields for 'time' type")
			}
		})
	}
}

func TestCalculateMaintenanceStatus_KM(t *testing.T) {
	now := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		currentKM    int
		lastKM       int
		intervalKM   int
		wantStatus   string
		wantNextKM   int
		wantKMRemain int
	}{
		{
			name:         "ok - 501 km remaining",
			currentKM:    19499,
			lastKM:       10000,
			intervalKM:   10000, // due 20000 -> 501 km left
			wantStatus:   service.StatusOK,
			wantNextKM:   20000,
			wantKMRemain: 501,
		},
		{
			name:         "warning - exactly 500 km remaining",
			currentKM:    19500,
			lastKM:       10000,
			intervalKM:   10000, // due 20000 -> 500 km left
			wantStatus:   service.StatusWarning,
			wantNextKM:   20000,
			wantKMRemain: 500,
		},
		{
			name:         "warning - 0 km remaining",
			currentKM:    20000,
			lastKM:       10000,
			intervalKM:   10000,
			wantStatus:   service.StatusWarning,
			wantNextKM:   20000,
			wantKMRemain: 0,
		},
		{
			name:         "overdue - exceeded by 1 km",
			currentKM:    20001,
			lastKM:       10000,
			intervalKM:   10000,
			wantStatus:   service.StatusOverdue,
			wantNextKM:   20000,
			wantKMRemain: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			detail := model.VehicleMaintenanceDetail{
				ExpirationType:        "km",
				VehicleCurrentMileage: tc.currentKM,
				LastPerformedMileage:  tc.lastKM,
				IntervalKM:            intPtr(tc.intervalKM),
			}

			service.CalculateMaintenanceStatus(&detail, now)

			if detail.Status != tc.wantStatus {
				t.Errorf("status = %q, want %q", detail.Status, tc.wantStatus)
			}
			if detail.NextDueMileage == nil || *detail.NextDueMileage != tc.wantNextKM {
				t.Errorf("next_due_mileage = %v, want %d", detail.NextDueMileage, tc.wantNextKM)
			}
			if detail.KMRemaining == nil || *detail.KMRemaining != tc.wantKMRemain {
				t.Errorf("km_remaining = %v, want %d", detail.KMRemaining, tc.wantKMRemain)
			}
			if detail.NextDueDate != nil || detail.DaysRemaining != nil {
				t.Errorf("expected nil time fields for 'km' type")
			}
		})
	}
}

func TestCalculateMaintenanceStatus_Both(t *testing.T) {
	now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		lastDate   string
		daysInt    int
		currentKM  int
		lastKM     int
		kmInt      int
		wantStatus string
	}{
		{
			name:       "both ok -> ok",
			lastDate:   "2026-06-01",
			daysInt:    30, // 16 days left -> ok
			currentKM:  15000,
			lastKM:     10000,
			kmInt:      10000, // 5000 km left -> ok
			wantStatus: service.StatusOK,
		},
		{
			name:       "time warning + km ok -> warning",
			lastDate:   "2026-06-01",
			daysInt:    20, // due 2026-06-21 -> 6 days left -> warning
			currentKM:  15000,
			lastKM:     10000,
			kmInt:      10000, // 5000 km left -> ok
			wantStatus: service.StatusWarning,
		},
		{
			name:       "time ok + km warning -> warning",
			lastDate:   "2026-06-01",
			daysInt:    60, // ok
			currentKM:  19600,
			lastKM:     10000,
			kmInt:      10000, // 400 km left -> warning
			wantStatus: service.StatusWarning,
		},
		{
			name:       "time overdue + km ok -> overdue",
			lastDate:   "2026-05-01",
			daysInt:    30, // due 2026-05-31 -> overdue
			currentKM:  12000,
			lastKM:     10000,
			kmInt:      10000, // ok
			wantStatus: service.StatusOverdue,
		},
		{
			name:       "time ok + km overdue -> overdue",
			lastDate:   "2026-06-01",
			daysInt:    60, // ok
			currentKM:  20100,
			lastKM:     10000,
			kmInt:      10000, // overdue
			wantStatus: service.StatusOverdue,
		},
		{
			name:       "time warning + km overdue -> overdue",
			lastDate:   "2026-06-01",
			daysInt:    20, // warning
			currentKM:  20500,
			lastKM:     10000,
			kmInt:      10000, // overdue
			wantStatus: service.StatusOverdue,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			detail := model.VehicleMaintenanceDetail{
				ExpirationType:        "both",
				LastPerformedDate:     tc.lastDate,
				IntervalDays:          intPtr(tc.daysInt),
				VehicleCurrentMileage: tc.currentKM,
				LastPerformedMileage:  tc.lastKM,
				IntervalKM:            intPtr(tc.kmInt),
			}

			service.CalculateMaintenanceStatus(&detail, now)

			if detail.Status != tc.wantStatus {
				t.Errorf("status = %q, want %q", detail.Status, tc.wantStatus)
			}
		})
	}
}
