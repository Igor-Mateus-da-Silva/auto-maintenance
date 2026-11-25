package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"car-maintenance/internal/model"
	"car-maintenance/internal/service"
)

type VehicleMaintenanceHandler struct {
	db *sql.DB
}

func NewVehicleMaintenanceHandler(db *sql.DB) *VehicleMaintenanceHandler {
	return &VehicleMaintenanceHandler{db: db}
}

// ListByVehicle handles GET /api/vehicles/{id}/maintenances
func (h *VehicleMaintenanceHandler) ListByVehicle(w http.ResponseWriter, r *http.Request) {
	vehicleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || vehicleID <= 0 {
		writeError(w, http.StatusBadRequest, "id de veículo inválido")
		return
	}

	exists, err := recordExists(r.Context(), h.db, "vehicles", vehicleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao verificar veículo")
		return
	}
	if !exists {
		writeError(w, http.StatusNotFound, "veículo não encontrado")
		return
	}

	items, err := fetchMaintenanceDetails(r.Context(), h.db, &vehicleID, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao listar manutenções do veículo")
		return
	}

	service.SortMaintenancesByUrgency(items)
	writeJSON(w, http.StatusOK, items)
}

// Create handles POST /api/vehicles/{id}/maintenances
func (h *VehicleMaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	vehicleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || vehicleID <= 0 {
		writeError(w, http.StatusBadRequest, "id de veículo inválido")
		return
	}

	var req model.CreateVehicleMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.MaintenanceTypeID <= 0 {
		writeError(w, http.StatusBadRequest, "tipo de manutenção é obrigatório")
		return
	}

	if _, err := time.Parse("2006-01-02", req.LastPerformedDate); err != nil {
		writeError(w, http.StatusBadRequest, "data inválida, formato esperado YYYY-MM-DD")
		return
	}

	if req.LastPerformedMileage < 0 {
		writeError(w, http.StatusBadRequest, "a quilometragem da manutenção deve ser maior ou igual a zero")
		return
	}

	vehicleExists, err := recordExists(r.Context(), h.db, "vehicles", vehicleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao verificar veículo")
		return
	}
	if !vehicleExists {
		writeError(w, http.StatusNotFound, "veículo não encontrado")
		return
	}

	typeExists, err := recordExists(r.Context(), h.db, "maintenance_types", req.MaintenanceTypeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao verificar tipo de manutenção")
		return
	}
	if !typeExists {
		writeError(w, http.StatusNotFound, "tipo de manutenção não encontrado")
		return
	}

	insertQuery := `
		INSERT INTO vehicle_maintenances (vehicle_id, maintenance_type_id, last_performed_date, last_performed_mileage, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	res, err := h.db.ExecContext(r.Context(), insertQuery, vehicleID, req.MaintenanceTypeID, req.LastPerformedDate, req.LastPerformedMileage, req.Notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao vincular manutenção ao veículo")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao obter id do vínculo criado")
		return
	}

	detail, err := fetchMaintenanceDetailByID(r.Context(), h.db, id, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar manutenção criada")
		return
	}

	writeJSON(w, http.StatusCreated, detail)
}

// Update handles PUT /api/vehicle-maintenances/{id}
func (h *VehicleMaintenanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var req model.UpdateVehicleMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if _, err := time.Parse("2006-01-02", req.LastPerformedDate); err != nil {
		writeError(w, http.StatusBadRequest, "data inválida, formato esperado YYYY-MM-DD")
		return
	}

	if req.LastPerformedMileage < 0 {
		writeError(w, http.StatusBadRequest, "a quilometragem da manutenção deve ser maior ou igual a zero")
		return
	}

	updateQuery := `
		UPDATE vehicle_maintenances
		SET last_performed_date = ?, last_performed_mileage = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, req.LastPerformedDate, req.LastPerformedMileage, req.Notes, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao atualizar manutenção do veículo")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "manutenção do veículo não encontrada")
		return
	}

	detail, err := fetchMaintenanceDetailByID(r.Context(), h.db, id, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar manutenção atualizada")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

// Delete handles DELETE /api/vehicle-maintenances/{id}
func (h *VehicleMaintenanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	res, err := h.db.ExecContext(r.Context(), `DELETE FROM vehicle_maintenances WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao remover manutenção do veículo")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "manutenção do veículo não encontrada")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func recordExists(ctx context.Context, db *sql.DB, table string, id int64) (bool, error) {
	var dummy int
	var query string
	switch table {
	case "vehicles":
		query = "SELECT 1 FROM vehicles WHERE id = ? LIMIT 1"
	case "maintenance_types":
		query = "SELECT 1 FROM maintenance_types WHERE id = ? LIMIT 1"
	default:
		return false, errors.New("invalid table")
	}

	err := db.QueryRowContext(ctx, query, id).Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

const baseMaintenanceDetailQuery = `
	SELECT
		vm.id,
		vm.vehicle_id,
		v.name AS vehicle_name,
		v.current_mileage AS vehicle_current_mileage,
		vm.maintenance_type_id,
		mt.name AS maintenance_type_name,
		mt.description AS maintenance_type_description,
		mt.expiration_type,
		mt.interval_days,
		mt.interval_km,
		vm.last_performed_date,
		vm.last_performed_mileage,
		vm.notes
	FROM vehicle_maintenances vm
	INNER JOIN vehicles v ON v.id = vm.vehicle_id
	INNER JOIN maintenance_types mt ON mt.id = vm.maintenance_type_id
`

func scanMaintenanceDetail(scanner interface {
	Scan(dest ...any) error
}, now time.Time) (model.VehicleMaintenanceDetail, error) {
	var (
		d            model.VehicleMaintenanceDetail
		mtDesc       sql.NullString
		intervalDays sql.NullInt64
		intervalKM   sql.NullInt64
		notes        sql.NullString
	)

	err := scanner.Scan(
		&d.ID,
		&d.VehicleID,
		&d.VehicleName,
		&d.VehicleCurrentMileage,
		&d.MaintenanceTypeID,
		&d.MaintenanceTypeName,
		&mtDesc,
		&d.ExpirationType,
		&intervalDays,
		&intervalKM,
		&d.LastPerformedDate,
		&d.LastPerformedMileage,
		&notes,
	)
	if err != nil {
		return d, err
	}

	if mtDesc.Valid {
		d.MaintenanceTypeDescription = &mtDesc.String
	}
	if intervalDays.Valid {
		days := int(intervalDays.Int64)
		d.IntervalDays = &days
	}
	if intervalKM.Valid {
		km := int(intervalKM.Int64)
		d.IntervalKM = &km
	}
	if notes.Valid {
		d.Notes = &notes.String
	}

	service.CalculateMaintenanceStatus(&d, now)
	return d, nil
}

func fetchMaintenanceDetails(ctx context.Context, db *sql.DB, filterVehicleID *int64, now time.Time) ([]model.VehicleMaintenanceDetail, error) {
	query := baseMaintenanceDetailQuery
	var args []any

	if filterVehicleID != nil {
		query += " WHERE vm.vehicle_id = ?"
		args = append(args, *filterVehicleID)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.VehicleMaintenanceDetail, 0)
	for rows.Next() {
		item, err := scanMaintenanceDetail(rows, now)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func fetchMaintenanceDetailByID(ctx context.Context, db *sql.DB, id int64, now time.Time) (*model.VehicleMaintenanceDetail, error) {
	query := baseMaintenanceDetailQuery + " WHERE vm.id = ?"
	row := db.QueryRowContext(ctx, query, id)
	item, err := scanMaintenanceDetail(row, now)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
