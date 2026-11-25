import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
  CreateMaintenanceTypeRequest,
  MaintenanceType,
  UpdateMaintenanceTypeRequest
} from '../../shared/models/maintenance.model';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class MaintenanceTypeService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;

  list(): Observable<MaintenanceType[]> {
    return this.http.get<MaintenanceType[]>(`${this.apiUrl}/maintenance-types`);
  }

  create(payload: CreateMaintenanceTypeRequest): Observable<MaintenanceType> {
    return this.http.post<MaintenanceType>(`${this.apiUrl}/maintenance-types`, payload);
  }

  update(id: number, payload: UpdateMaintenanceTypeRequest): Observable<MaintenanceType> {
    return this.http.put<MaintenanceType>(`${this.apiUrl}/maintenance-types/${id}`, payload);
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/maintenance-types/${id}`);
  }
}
