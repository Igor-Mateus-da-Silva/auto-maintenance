package model

type VehicleMaintenanceDetail struct {
	ID                         int64   `json:"id"`
	VehicleID                  int64   `json:"vehicle_id"`
	VehicleName                string  `json:"vehicle_name"`
	VehicleCurrentMileage      int     `json:"vehicle_current_mileage"`
	MaintenanceTypeID          int64   `json:"maintenance_type_id"`
	MaintenanceTypeName        string  `json:"maintenance_type_name"`
	MaintenanceTypeDescription *string `json:"maintenance_type_description"`
	ExpirationType             string  `json:"expiration_type"`
	IntervalDays               *int    `json:"interval_days"`
	IntervalKM                 *int    `json:"interval_km"`
	LastPerformedDate          string  `json:"last_performed_date"`
	LastPerformedMileage       int     `json:"last_performed_mileage"`
	Notes                      *string `json:"notes"`
	Status                     string  `json:"status"`
	NextDueDate                *string `json:"next_due_date"`
	DaysRemaining              *int    `json:"days_remaining"`
	NextDueMileage             *int    `json:"next_due_mileage"`
	KMRemaining                *int    `json:"km_remaining"`
}

type CreateVehicleMaintenanceRequest struct {
	MaintenanceTypeID    int64   `json:"maintenance_type_id"`
	LastPerformedDate    string  `json:"last_performed_date"`
	LastPerformedMileage int     `json:"last_performed_mileage"`
	Notes                *string `json:"notes"`
}

type UpdateVehicleMaintenanceRequest struct {
	LastPerformedDate    string  `json:"last_performed_date"`
	LastPerformedMileage int     `json:"last_performed_mileage"`
	Notes                *string `json:"notes"`
}
