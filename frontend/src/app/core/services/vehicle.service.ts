import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
  CreateVehicleRequest,
  UpdateVehicleRequest,
  Vehicle,
} from '../../shared/models/vehicle.model';
import {
  CreateVehicleMaintenanceRequest,
  UpdateVehicleMaintenanceRequest,
  VehicleMaintenanceDetail,
} from '../../shared/models/maintenance.model';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class VehicleService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;

  list(): Observable<Vehicle[]> {
    return this.http.get<Vehicle[]>(`${this.apiUrl}/vehicles`);
  }

  create(payload: CreateVehicleRequest): Observable<Vehicle> {
    return this.http.post<Vehicle>(`${this.apiUrl}/vehicles`, payload);
  }

  update(id: number, payload: UpdateVehicleRequest): Observable<Vehicle> {
    return this.http.put<Vehicle>(`${this.apiUrl}/vehicles/${id}`, payload);
  }

  updateMileage(vehicleId: number, currentMileage: number): Observable<Vehicle> {
    return this.http.patch<Vehicle>(`${this.apiUrl}/vehicles/${vehicleId}/mileage`, {
      current_mileage: currentMileage,
    });
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/vehicles/${id}`);
  }

  listMaintenances(vehicleId: number): Observable<VehicleMaintenanceDetail[]> {
    return this.http.get<VehicleMaintenanceDetail[]>(
      `${this.apiUrl}/vehicles/${vehicleId}/maintenances`,
    );
  }

  attachMaintenance(
    vehicleId: number,
    payload: CreateVehicleMaintenanceRequest,
  ): Observable<VehicleMaintenanceDetail> {
    return this.http.post<VehicleMaintenanceDetail>(
      `${this.apiUrl}/vehicles/${vehicleId}/maintenances`,
      payload,
    );
  }

  renewMaintenance(
    id: number,
    payload: UpdateVehicleMaintenanceRequest,
  ): Observable<VehicleMaintenanceDetail> {
    return this.http.put<VehicleMaintenanceDetail>(
      `${this.apiUrl}/vehicle-maintenances/${id}`,
      payload,
    );
  }

  deleteMaintenance(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/vehicle-maintenances/${id}`);
  }
}
