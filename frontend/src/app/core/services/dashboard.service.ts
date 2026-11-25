import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { DashboardResponse } from '../../shared/models/maintenance.model';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api';

  getDashboard(vehicleId?: number): Observable<DashboardResponse> {
    let params = new HttpParams();
    if (vehicleId !== undefined && vehicleId !== null) {
      params = params.set('vehicle_id', vehicleId.toString());
    }
    return this.http.get<DashboardResponse>(`${this.apiUrl}/dashboard`, { params });
  }
}
