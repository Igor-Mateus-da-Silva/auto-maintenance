import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Vehicle } from '../../shared/models/vehicle.model';
import {
  UpdateVehicleMaintenanceRequest,
  VehicleMaintenanceDetail
} from '../../shared/models/maintenance.model';

@Injectable({
  providedIn: 'root'
})
export class VehicleService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api';

  updateMileage(vehicleId: number, currentMileage: number): Observable<Vehicle> {
    return this.http.patch<Vehicle>(`${this.apiUrl}/vehicles/${vehicleId}/mileage`, {
      current_mileage: currentMileage
    });
  }

  renewMaintenance(
    id: number,
    payload: UpdateVehicleMaintenanceRequest
  ): Observable<VehicleMaintenanceDetail> {
    return this.http.put<VehicleMaintenanceDetail>(
      `${this.apiUrl}/vehicle-maintenances/${id}`,
      payload
    );
  }
}
