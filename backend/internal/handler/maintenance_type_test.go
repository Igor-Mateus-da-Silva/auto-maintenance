package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"car-maintenance/internal/handler"
	"car-maintenance/internal/model"
)

func TestMaintenanceType_CRUD_And_Validations(t *testing.T) {
	db := setupTestDB(t)
	h := handler.NewMaintenanceTypeHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/maintenance-types", h.List)
	mux.HandleFunc("POST /api/maintenance-types", h.Create)
	mux.HandleFunc("PUT /api/maintenance-types/{id}", h.Update)
	mux.HandleFunc("DELETE /api/maintenance-types/{id}", h.Delete)

	days180 := 180
	daysZero := 0
	km10000 := 10000
	kmZero := 0

	// 1. Validation test cases for POST
	cases := []struct {
		name       string
		payload    model.CreateMaintenanceTypeRequest
		wantStatus int
	}{
		{
			name:       "empty name",
			payload:    model.CreateMaintenanceTypeRequest{Name: "", ExpirationType: "time", IntervalDays: &days180},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid expiration type",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "invalid", IntervalDays: &days180},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "time without interval_days",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "time", IntervalDays: nil},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "time with interval_days <= 0",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "time", IntervalDays: &daysZero},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "km without interval_km",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "km", IntervalKM: nil},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "km with interval_km <= 0",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "km", IntervalKM: &kmZero},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "both missing interval_km",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "both", IntervalDays: &days180, IntervalKM: nil},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "both missing interval_days",
			payload:    model.CreateMaintenanceTypeRequest{Name: "Troca", ExpirationType: "both", IntervalDays: nil, IntervalKM: &km10000},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/maintenance-types", bytes.NewReader(b))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d: %s", tc.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}

	// 2. Successful creation: "time" (interval_km must be saved as NULL)
	descTime := "Revisão por tempo"
	timePayload, _ := json.Marshal(model.CreateMaintenanceTypeRequest{
		Name:           "Revisão Anual",
		Description:    &descTime,
		ExpirationType: "time",
		IntervalDays:   &days180,
		IntervalKM:     &km10000, // Should be ignored/cleared to nil
	})
	req := httptest.NewRequest(http.MethodPost, "/api/maintenance-types", bytes.NewReader(timePayload))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var timeCreated model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&timeCreated)
	if timeCreated.IntervalKM != nil {
		t.Errorf("expected interval_km to be nil for 'time' type, got %v", *timeCreated.IntervalKM)
	}
	if timeCreated.IntervalDays == nil || *timeCreated.IntervalDays != 180 {
		t.Errorf("expected interval_days to be 180, got %v", timeCreated.IntervalDays)
	}

	// 3. Successful creation: "km" (interval_days must be saved as NULL)
	kmPayload, _ := json.Marshal(model.CreateMaintenanceTypeRequest{
		Name:           "Troca de Pastilhas",
		ExpirationType: "km",
		IntervalDays:   &days180, // Should be ignored/cleared to nil
		IntervalKM:     &km10000,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/maintenance-types", bytes.NewReader(kmPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var kmCreated model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&kmCreated)
	if kmCreated.IntervalDays != nil {
		t.Errorf("expected interval_days to be nil for 'km' type, got %v", *kmCreated.IntervalDays)
	}
	if kmCreated.IntervalKM == nil || *kmCreated.IntervalKM != 10000 {
		t.Errorf("expected interval_km to be 10000, got %v", kmCreated.IntervalKM)
	}

	// 4. Successful creation: "both" (both must be saved)
	bothPayload, _ := json.Marshal(model.CreateMaintenanceTypeRequest{
		Name:           "Troca de Óleo e Filtro",
		ExpirationType: "both",
		IntervalDays:   &days180,
		IntervalKM:     &km10000,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/maintenance-types", bytes.NewReader(bothPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var bothCreated model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&bothCreated)
	if bothCreated.IntervalDays == nil || *bothCreated.IntervalDays != 180 || bothCreated.IntervalKM == nil || *bothCreated.IntervalKM != 10000 {
		t.Errorf("expected both intervals to be set, got %+v", bothCreated)
	}

	// 5. List should return all 3 types ordered by name
	req = httptest.NewRequest(http.MethodGet, "/api/maintenance-types", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var list []model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}

	// 6. Update (PUT)
	newDays := 365
	updatePayload, _ := json.Marshal(model.UpdateMaintenanceTypeRequest{
		Name:           "Revisão Anual Completa",
		ExpirationType: "time",
		IntervalDays:   &newDays,
	})
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance-types/%d", timeCreated.ID), bytes.NewReader(updatePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&updated)
	if updated.Name != "Revisão Anual Completa" || *updated.IntervalDays != 365 {
		t.Errorf("unexpected updated data: %+v", updated)
	}

	// 7. Update non-existent returns 404
	req = httptest.NewRequest(http.MethodPut, "/api/maintenance-types/9999", bytes.NewReader(updatePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 on updating non-existent, got %d", rec.Code)
	}

	// 8. Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/maintenance-types/%d", timeCreated.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", rec.Code)
	}

	// 9. Delete non-existent returns 404
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/maintenance-types/%d", timeCreated.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 on deleting already deleted item, got %d", rec.Code)
	}
}
