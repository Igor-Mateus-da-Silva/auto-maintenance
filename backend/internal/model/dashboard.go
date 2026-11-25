package model

type DashboardSummary struct {
	Total   int `json:"total"`
	OK      int `json:"ok"`
	Warning int `json:"warning"`
	Overdue int `json:"overdue"`
}

type DashboardResponse struct {
	Summary      DashboardSummary           `json:"summary"`
	Vehicles     []Vehicle                  `json:"vehicles"`
	Maintenances []VehicleMaintenanceDetail `json:"maintenances"`
}
