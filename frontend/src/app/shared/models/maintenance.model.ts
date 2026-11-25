import { Vehicle } from './vehicle.model';

export type ExpirationType = 'time' | 'km' | 'both';
export type MaintenanceStatus = 'ok' | 'warning' | 'overdue';

export interface MaintenanceType {
  id: number;
  name: string;
  description?: string | null;
  expiration_type: ExpirationType;
  interval_days?: number | null;
  interval_km?: number | null;
  created_at: string;
}

export interface CreateMaintenanceTypeRequest {
  name: string;
  description?: string | null;
  expiration_type: ExpirationType;
  interval_days?: number | null;
  interval_km?: number | null;
}

export interface UpdateMaintenanceTypeRequest {
  name: string;
  description?: string | null;
  expiration_type: ExpirationType;
  interval_days?: number | null;
  interval_km?: number | null;
}

export interface VehicleMaintenanceDetail {
  id: number;
  vehicle_id: number;
  vehicle_name: string;
  vehicle_current_mileage: number;
  maintenance_type_id: number;
  maintenance_type_name: string;
  maintenance_type_description?: string | null;
  expiration_type: ExpirationType;
  interval_days?: number | null;
  interval_km?: number | null;
  last_performed_date: string;
  last_performed_mileage: number;
  notes?: string | null;
  status: MaintenanceStatus;
  next_due_date?: string | null;
  days_remaining?: number | null;
  next_due_mileage?: number | null;
  km_remaining?: number | null;
}

export interface CreateVehicleMaintenanceRequest {
  maintenance_type_id: number;
  last_performed_date: string;
  last_performed_mileage: number;
  notes?: string | null;
}

export interface UpdateVehicleMaintenanceRequest {
  last_performed_date: string;
  last_performed_mileage: number;
  notes?: string | null;
}

export interface DashboardSummary {
  total: number;
  ok: number;
  warning: number;
  overdue: number;
}

export interface DashboardResponse {
  summary: DashboardSummary;
  vehicles: Vehicle[];
  maintenances: VehicleMaintenanceDetail[];
}
