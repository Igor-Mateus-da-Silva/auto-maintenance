package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"car-maintenance/internal/model"
)

type VehicleHandler struct {
	db *sql.DB
}

func NewVehicleHandler(db *sql.DB) *VehicleHandler {
	return &VehicleHandler{db: db}
}

// List handles GET /api/vehicles
func (h *VehicleHandler) List(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, plate, year, current_mileage, created_at, updated_at
		FROM vehicles
		ORDER BY name ASC
	`
	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao listar veículos")
		return
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
			writeError(w, http.StatusInternalServerError, "erro ao ler dados do veículo")
			return
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
		writeError(w, http.StatusInternalServerError, "erro ao iterar veículos")
		return
	}

	writeJSON(w, http.StatusOK, vehicles)
}

// Create handles POST /api/vehicles
func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVehicleRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	name, plate, err := validateVehicleInput(req.Name, req.Plate, req.Year, req.CurrentMileage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	insertQuery := `
		INSERT INTO vehicles (name, plate, year, current_mileage, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	res, err := h.db.ExecContext(r.Context(), insertQuery, name, plate, req.Year, req.CurrentMileage)
	if err != nil {
		writeInternalError(w, "erro ao criar veículo", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao obter id do veículo criado")
		return
	}

	created, err := h.getByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar veículo criado")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// Update handles PUT /api/vehicles/{id}
func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id de veículo inválido")
		return
	}

	var req model.UpdateVehicleRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	name, plate, err := validateVehicleInput(req.Name, req.Plate, req.Year, req.CurrentMileage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateQuery := `
		UPDATE vehicles
		SET name = ?, plate = ?, year = ?, current_mileage = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, name, plate, req.Year, req.CurrentMileage, id)
	if err != nil {
		writeInternalError(w, "erro ao atualizar veículo", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "veículo não encontrado")
		return
	}

	updated, err := h.getByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar veículo atualizado")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// UpdateMileage handles PATCH /api/vehicles/{id}/mileage
func (h *VehicleHandler) UpdateMileage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id de veículo inválido")
		return
	}

	var req model.UpdateMileageRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if req.CurrentMileage < 0 || req.CurrentMileage > 10_000_000 {
		writeError(w, http.StatusBadRequest, "a quilometragem atual deve estar entre 0 e 10.000.000 km")
		return
	}

	updateQuery := `
		UPDATE vehicles
		SET current_mileage = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, req.CurrentMileage, id)
	if err != nil {
		writeInternalError(w, "erro ao atualizar quilometragem", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "veículo não encontrado")
		return
	}

	updated, err := h.getByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar veículo atualizado")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/vehicles/{id}
func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id de veículo inválido")
		return
	}

	res, err := h.db.ExecContext(r.Context(), `DELETE FROM vehicles WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao remover veículo")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "veículo não encontrado")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (h *VehicleHandler) getByID(ctx context.Context, id int64) (*model.Vehicle, error) {
	query := `
		SELECT id, name, plate, year, current_mileage, created_at, updated_at
		FROM vehicles
		WHERE id = ?
	`
	var (
		v     model.Vehicle
		plate sql.NullString
		year  sql.NullInt64
	)
	err := h.db.QueryRowContext(ctx, query, id).Scan(&v.ID, &v.Name, &plate, &year, &v.CurrentMileage, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if plate.Valid {
		v.Plate = &plate.String
	}
	if year.Valid {
		y := int(year.Int64)
		v.Year = &y
	}
	return &v, nil
}

func validateVehicleInput(name string, plate *string, year *int, mileage int) (string, *string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" || len(trimmedName) > 120 {
		return "", nil, errors.New("o nome do veículo é obrigatório e deve ter no máximo 120 caracteres")
	}

	var cleanPlate *string
	if plate != nil {
		trimmedPlate := strings.TrimSpace(*plate)
		if len(trimmedPlate) > 20 {
			return "", nil, errors.New("a placa deve ter no máximo 20 caracteres")
		}
		if trimmedPlate != "" {
			cleanPlate = &trimmedPlate
		}
	}

	if year != nil && (*year < 1900 || *year > 2100) {
		return "", nil, errors.New("o ano do veículo deve estar entre 1900 e 2100")
	}

	if mileage < 0 || mileage > 10_000_000 {
		return "", nil, errors.New("a quilometragem atual deve estar entre 0 e 10.000.000 km")
	}

	return trimmedName, cleanPlate, nil
}
