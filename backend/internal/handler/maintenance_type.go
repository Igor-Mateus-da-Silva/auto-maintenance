package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"car-maintenance/internal/model"
)

type MaintenanceTypeHandler struct {
	db *sql.DB
}

func NewMaintenanceTypeHandler(db *sql.DB) *MaintenanceTypeHandler {
	return &MaintenanceTypeHandler{db: db}
}

// List handles GET /api/maintenance-types
func (h *MaintenanceTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, description, expiration_type, interval_days, interval_km, created_at
		FROM maintenance_types
		ORDER BY name ASC
	`
	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao listar tipos de manutenção")
		return
	}
	defer rows.Close()

	items := make([]model.MaintenanceType, 0)
	for rows.Next() {
		var (
			mt           model.MaintenanceType
			desc         sql.NullString
			intervalDays sql.NullInt64
			intervalKM   sql.NullInt64
		)
		if err := rows.Scan(&mt.ID, &mt.Name, &desc, &mt.ExpirationType, &intervalDays, &intervalKM, &mt.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "erro ao ler dados do tipo de manutenção")
			return
		}
		if desc.Valid {
			mt.Description = &desc.String
		}
		if intervalDays.Valid {
			days := int(intervalDays.Int64)
			mt.IntervalDays = &days
		}
		if intervalKM.Valid {
			km := int(intervalKM.Int64)
			mt.IntervalKM = &km
		}
		items = append(items, mt)
	}

	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao iterar tipos de manutenção")
		return
	}

	writeJSON(w, http.StatusOK, items)
}

// Create handles POST /api/maintenance-types
func (h *MaintenanceTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateMaintenanceTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	cleanDays, cleanKM, err := h.validate(req.Name, req.ExpirationType, req.IntervalDays, req.IntervalKM)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	insertQuery := `
		INSERT INTO maintenance_types (name, description, expiration_type, interval_days, interval_km, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	res, err := h.db.ExecContext(r.Context(), insertQuery, strings.TrimSpace(req.Name), req.Description, req.ExpirationType, cleanDays, cleanKM)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao criar tipo de manutenção")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao obter id do tipo de manutenção")
		return
	}

	created, err := h.getByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tipo de manutenção criado")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// Update handles PUT /api/maintenance-types/{id}
func (h *MaintenanceTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id de tipo de manutenção inválido")
		return
	}

	var req model.UpdateMaintenanceTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	cleanDays, cleanKM, err := h.validate(req.Name, req.ExpirationType, req.IntervalDays, req.IntervalKM)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateQuery := `
		UPDATE maintenance_types
		SET name = ?, description = ?, expiration_type = ?, interval_days = ?, interval_km = ?
		WHERE id = ?
	`
	res, err := h.db.ExecContext(r.Context(), updateQuery, strings.TrimSpace(req.Name), req.Description, req.ExpirationType, cleanDays, cleanKM, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao atualizar tipo de manutenção")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "tipo de manutenção não encontrado")
		return
	}

	updated, err := h.getByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tipo de manutenção atualizado")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/maintenance-types/{id}
func (h *MaintenanceTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id de tipo de manutenção inválido")
		return
	}

	res, err := h.db.ExecContext(r.Context(), `DELETE FROM maintenance_types WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao remover tipo de manutenção")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "tipo de manutenção não encontrado")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (h *MaintenanceTypeHandler) validate(name, expirationType string, days, km *int) (*int, *int, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil, errors.New("o nome do tipo de manutenção é obrigatório")
	}

	switch expirationType {
	case "time":
		if days == nil || *days <= 0 {
			return nil, nil, errors.New("para o tipo 'time', o intervalo em dias é obrigatório e deve ser maior que zero")
		}
		return days, nil, nil
	case "km":
		if km == nil || *km <= 0 {
			return nil, nil, errors.New("para o tipo 'km', o intervalo em km é obrigatório e deve ser maior que zero")
		}
		return nil, km, nil
	case "both":
		if days == nil || *days <= 0 {
			return nil, nil, errors.New("para o tipo 'both', o intervalo em dias é obrigatório e deve ser maior que zero")
		}
		if km == nil || *km <= 0 {
			return nil, nil, errors.New("para o tipo 'both', o intervalo em km é obrigatório e deve ser maior que zero")
		}
		return days, km, nil
	default:
		return nil, nil, fmt.Errorf("tipo de expiração inválido: deve ser 'time', 'km' ou 'both'")
	}
}

func (h *MaintenanceTypeHandler) getByID(ctx context.Context, id int64) (*model.MaintenanceType, error) {
	query := `
		SELECT id, name, description, expiration_type, interval_days, interval_km, created_at
		FROM maintenance_types
		WHERE id = ?
	`
	var (
		mt           model.MaintenanceType
		desc         sql.NullString
		intervalDays sql.NullInt64
		intervalKM   sql.NullInt64
	)
	err := h.db.QueryRowContext(ctx, query, id).Scan(&mt.ID, &mt.Name, &desc, &mt.ExpirationType, &intervalDays, &intervalKM, &mt.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if desc.Valid {
		mt.Description = &desc.String
	}
	if intervalDays.Valid {
		days := int(intervalDays.Int64)
		mt.IntervalDays = &days
	}
	if intervalKM.Valid {
		km := int(intervalKM.Int64)
		mt.IntervalKM = &km
	}
	return &mt, nil
}
