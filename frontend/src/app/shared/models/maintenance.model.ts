export type MaintenanceType = 'preventive' | 'corrective' | 'emergency';
export type MaintenanceStatus = 'scheduled' | 'in_progress' | 'completed' | 'cancelled';

export interface Maintenance {
  id: string;
  vehicleId: string;
  description: string;
  type: MaintenanceType;
  status: MaintenanceStatus;
  mileageAtService: number;
  cost: number;
  serviceDate: string;
  notes?: string;
}
