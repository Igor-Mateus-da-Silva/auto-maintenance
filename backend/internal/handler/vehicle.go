package handler

import (
	"context"
	"database/sql"
	"encoding/json"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "o nome do veículo é obrigatório")
		return
	}
	if req.CurrentMileage < 0 {
		writeError(w, http.StatusBadRequest, "a quilometragem atual deve ser maior ou igual a zero")
		return
	}

	insertQuery := `
		INSERT INTO vehicles (name, plate, year, current_mileage, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	res, err := h.db.ExecContext(r.Context(), insertQuery, name, req.Plate, req.Year, req.CurrentMileage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao criar veículo")
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "o nome do veículo é obrigatório")
		return
	}
	if req.CurrentMileage < 0 {
		writeError(w, http.StatusBadRequest, "a quilometragem atual deve ser maior ou igual a zero")
		return
	}

	updateQuery := `
		UPDATE vehicles
		SET name = ?, plate = ?, year = ?, current_mileage = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, name, req.Plate, req.Year, req.CurrentMileage, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao atualizar veículo")
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.CurrentMileage < 0 {
		writeError(w, http.StatusBadRequest, "a quilometragem atual deve ser maior ou igual a zero")
		return
	}

	updateQuery := `
		UPDATE vehicles
		SET current_mileage = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, req.CurrentMileage, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao atualizar quilometragem")
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
