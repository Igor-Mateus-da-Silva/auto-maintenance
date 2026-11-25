package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"car-maintenance/internal/handler"
	"car-maintenance/internal/model"
)

func TestDashboard_FullMaintenanceLifecycle(t *testing.T) {
	db := setupTestDB(t)

	vh := handler.NewVehicleHandler(db)
	mth := handler.NewMaintenanceTypeHandler(db)
	vmh := handler.NewVehicleMaintenanceHandler(db)
	dh := handler.NewDashboardHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/vehicles", vh.Create)
	mux.HandleFunc("PATCH /api/vehicles/{id}/mileage", vh.UpdateMileage)
	mux.HandleFunc("POST /api/maintenance-types", mth.Create)
	mux.HandleFunc("GET /api/vehicles/{id}/maintenances", vmh.ListByVehicle)
	mux.HandleFunc("POST /api/vehicles/{id}/maintenances", vmh.Create)
	mux.HandleFunc("PUT /api/vehicle-maintenances/{id}", vmh.Update)
	mux.HandleFunc("DELETE /api/vehicle-maintenances/{id}", vmh.Delete)
	mux.HandleFunc("GET /api/dashboard", dh.Get)

	// 1. Create vehicle at 10,000 km
	vPayload, _ := json.Marshal(model.CreateVehicleRequest{
		Name:           "Civic Touring",
		CurrentMileage: 10000,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewReader(vPayload))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create vehicle failed: %d %s", rec.Code, rec.Body.String())
	}
	var vehicle model.Vehicle
	_ = json.NewDecoder(rec.Body).Decode(&vehicle)

	// 2. Create maintenance type (every 5,000 km)
	km5000 := 5000
	mtPayload, _ := json.Marshal(model.CreateMaintenanceTypeRequest{
		Name:           "Troca de Óleo",
		ExpirationType: "km",
		IntervalKM:     &km5000,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/maintenance-types", bytes.NewReader(mtPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create maintenance type failed: %d %s", rec.Code, rec.Body.String())
	}
	var mType model.MaintenanceType
	_ = json.NewDecoder(rec.Body).Decode(&mType)

	// 3. Link maintenance to vehicle (last performed at 10,000 km -> due at 15,000 km)
	today := time.Now().Format("2006-01-02")
	vmPayload, _ := json.Marshal(model.CreateVehicleMaintenanceRequest{
		MaintenanceTypeID:    mType.ID,
		LastPerformedDate:    today,
		LastPerformedMileage: 10000,
	})
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/vehicles/%d/maintenances", vehicle.ID), bytes.NewReader(vmPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("link maintenance failed: %d %s", rec.Code, rec.Body.String())
	}
	var vmDetail model.VehicleMaintenanceDetail
	_ = json.NewDecoder(rec.Body).Decode(&vmDetail)
	if vmDetail.Status != "ok" {
		t.Errorf("expected initial status 'ok', got %q", vmDetail.Status)
	}

	// 4. Verify Dashboard shows 1 OK maintenance
	req = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get dashboard failed: %d", rec.Code)
	}
	var dash model.DashboardResponse
	_ = json.NewDecoder(rec.Body).Decode(&dash)
	if dash.Summary.Total != 1 || dash.Summary.OK != 1 || dash.Summary.Warning != 0 || dash.Summary.Overdue != 0 {
		t.Errorf("unexpected summary at 10000km: %+v", dash.Summary)
	}

	// 5. Update vehicle KM to 14,600 (400 km remaining <= 500 -> warning)
	kmPatch, _ := json.Marshal(model.UpdateMileageRequest{CurrentMileage: 14600})
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/vehicles/%d/mileage", vehicle.ID), bytes.NewReader(kmPatch))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch mileage failed: %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/dashboard?vehicle_id=%d", vehicle.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	_ = json.NewDecoder(rec.Body).Decode(&dash)
	if dash.Summary.Warning != 1 || dash.Maintenances[0].Status != "warning" {
		t.Errorf("expected warning status at 14600km, got summary=%+v status=%s", dash.Summary, dash.Maintenances[0].Status)
	}

	// 6. Update vehicle KM to 15,200 (-200 km remaining -> overdue)
	kmPatch, _ = json.Marshal(model.UpdateMileageRequest{CurrentMileage: 15200})
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/vehicles/%d/mileage", vehicle.ID), bytes.NewReader(kmPatch))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch mileage failed: %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	_ = json.NewDecoder(rec.Body).Decode(&dash)
	if dash.Summary.Overdue != 1 || dash.Maintenances[0].Status != "overdue" {
		t.Errorf("expected overdue status at 15200km, got summary=%+v status=%s", dash.Summary, dash.Maintenances[0].Status)
	}

	// 7. Renew maintenance at 15,200 km -> due at 20,200 km -> status returns to ok
	renewPayload, _ := json.Marshal(model.UpdateVehicleMaintenanceRequest{
		LastPerformedDate:    today,
		LastPerformedMileage: 15200,
	})
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/vehicle-maintenances/%d", vmDetail.ID), bytes.NewReader(renewPayload))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("renew maintenance failed: %d %s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	_ = json.NewDecoder(rec.Body).Decode(&dash)
	if dash.Summary.OK != 1 || dash.Maintenances[0].Status != "ok" {
		t.Errorf("expected ok status after renewal, got summary=%+v", dash.Summary)
	}

	// 8. Delete maintenance link
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/vehicle-maintenances/%d", vmDetail.ID), nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", rec.Code)
	}
}
