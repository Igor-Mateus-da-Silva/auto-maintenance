package model

type MaintenanceType struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	ExpirationType string  `json:"expiration_type"` // "time", "km", "both"
	IntervalDays   *int    `json:"interval_days"`
	IntervalKM     *int    `json:"interval_km"`
	CreatedAt      string  `json:"created_at"`
}

type CreateMaintenanceTypeRequest struct {
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	ExpirationType string  `json:"expiration_type"`
	IntervalDays   *int    `json:"interval_days"`
	IntervalKM     *int    `json:"interval_km"`
}

type UpdateMaintenanceTypeRequest struct {
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	ExpirationType string  `json:"expiration_type"`
	IntervalDays   *int    `json:"interval_days"`
	IntervalKM     *int    `json:"interval_km"`
}
