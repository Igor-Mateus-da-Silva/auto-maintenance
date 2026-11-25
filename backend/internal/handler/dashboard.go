package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"car-maintenance/internal/model"
	"car-maintenance/internal/service"
)

type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// Get handles GET /api/dashboard
func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	var filterVehicleID *int64
	if vehicleIDStr := r.URL.Query().Get("vehicle_id"); vehicleIDStr != "" {
		id, err := strconv.ParseInt(vehicleIDStr, 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "vehicle_id inválido")
			return
		}
		filterVehicleID = &id
	}

	// Fetch all vehicles for the selector and quick mileage update
	vehicles, err := h.listVehicles(r)
	if err != nil {
		writeInternalError(w, "erro ao listar veículos para o dashboard", err)
		return
	}

	// Fetch maintenances (filtered if vehicle_id provided)
	maintenances, err := fetchMaintenanceDetails(r.Context(), h.db, filterVehicleID, time.Now())
	if err != nil {
		writeInternalError(w, "erro ao carregar manutenções para o dashboard", err)
		return
	}

	service.SortMaintenancesByUrgency(maintenances)

	summary := model.DashboardSummary{
		Total: len(maintenances),
	}
	for _, m := range maintenances {
		switch m.Status {
		case service.StatusOK:
			summary.OK++
		case service.StatusWarning:
			summary.Warning++
		case service.StatusOverdue:
			summary.Overdue++
		}
	}

	resp := model.DashboardResponse{
		Summary:      summary,
		Vehicles:     vehicles,
		Maintenances: maintenances,
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) listVehicles(r *http.Request) ([]model.Vehicle, error) {
	query := `
		SELECT id, name, plate, year, current_mileage, created_at, updated_at
		FROM vehicles
		ORDER BY name ASC
	`
	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := make([]model.Vehicle, 0)
	for rows.Next() {
		var (
			v     model.Vehicle
			plate sql.NullString
			year  sql.NullInt64
		)
		if err := rows.Scan(&v.ID, &v.Name, &plate, &year, &v.CurrentMileage, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if plate.Valid {
			v.Plate = &plate.String
		}
		if year.Valid {
			y := int(year.Int64)
			v.Year = &y
		}
		vehicles = append(vehicles, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}
