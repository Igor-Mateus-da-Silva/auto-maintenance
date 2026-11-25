package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"car-maintenance/internal/database"
	"car-maintenance/internal/handler"
	"car-maintenance/internal/model"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func TestVehicle_CRUD_And_Mileage(t *testing.T) {
	db := setupTestDB(t)
	h := handler.NewVehicleHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/vehicles", h.List)
	mux.HandleFunc("POST /api/vehicles", h.Create)
	mux.HandleFunc("PUT /api/vehicles/{id}", h.Update)
	mux.HandleFunc("PATCH /api/vehicles/{id}/mileage", h.UpdateMileage)
	mux.HandleFunc("DELETE /api/vehicles/{id}", h.Delete)

	// 1. Initial list should be empty
	req := httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var list []model.Vehicle
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("failed to decode list: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 vehicles, got %d", len(list))
	}

	// 2. Validation error on Create: empty name
	badPayload, _ := json.Marshal(model.CreateVehicleRequest{Name: "", CurrentMileage: 1000})
	req = httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewReader(badPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty name, got %d", rec.Code)
	}

	// 3. Validation error on Create: negative mileage
	badPayload, _ = json.Marshal(model.CreateVehicleRequest{Name: "Corolla", CurrentMileage: -5})
	req = httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewReader(badPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative mileage, got %d", rec.Code)
	}

	// 4. Successful Create
	plate := "ABC-1234"
	year := 2022
	createPayload, _ := json.Marshal(model.CreateVehicleRequest{
		Name:           "Corolla XEi",
		Plate:          &plate,
		Year:           &year,
		CurrentMileage: 45000,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewReader(createPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var created model.Vehicle
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode created vehicle: %v", err)
	}
	if created.ID == 0 || created.Name != "Corolla XEi" || *created.Plate != "ABC-1234" || created.CurrentMileage != 45000 {
		t.Errorf("unexpected created vehicle data: %+v", created)
	}

	// 5. List should now contain 1 vehicle
	req = httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil || len(list) != 1 {
		t.Fatalf("expected 1 vehicle in list, got %d", len(list))
	}

	// 6. Fast mileage update (PATCH)
	mileagePayload, _ := json.Marshal(model.UpdateMileageRequest{CurrentMileage: 50000})
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/vehicles/%d/mileage", created.ID), bytes.NewReader(mileagePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on mileage update, got %d: %s", rec.Code, rec.Body.String())
	}
	var patched model.Vehicle
	if err := json.NewDecoder(rec.Body).Decode(&patched); err != nil {
		t.Fatalf("failed to decode patched vehicle: %v", err)
	}
	if patched.CurrentMileage != 50000 {
		t.Errorf("expected mileage 50000, got %d", patched.CurrentMileage)
	}

	// Negative mileage in PATCH should return 400
	badMileagePayload, _ := json.Marshal(model.UpdateMileageRequest{CurrentMileage: -100})
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/vehicles/%d/mileage", created.ID), bytes.NewReader(badMileagePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on negative mileage update, got %d", rec.Code)
	}

	// 7. Full update (PUT)
	newPlate := "XYZ-9876"
	newYear := 2023
	updatePayload, _ := json.Marshal(model.UpdateVehicleRequest{
		Name:           "Corolla Altis",
		Plate:          &newPlate,
		Year:           &newYear,
		CurrentMileage: 52000,
	})
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/vehicles/%d", created.ID), bytes.NewReader(updatePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated model.Vehicle
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode updated vehicle: %v", err)
	}
	if updated.Name != "Corolla Altis" || *updated.Plate != "XYZ-9876" || updated.CurrentMileage != 52000 {
		t.Errorf("unexpected updated data: %+v", updated)
	}

	// 8. Update non-existent vehicle returns 404
	req = httptest.NewRequest(http.MethodPut, "/api/vehicles/9999", bytes.NewReader(updatePayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent vehicle update, got %d", rec.Code)
	}

	// 9. Delete vehicle
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/vehicles/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", rec.Code)
	}

	// 10. Delete non-existent vehicle returns 404
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/vehicles/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 on deleting deleted vehicle, got %d", rec.Code)
	}
}
