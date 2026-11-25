export interface Vehicle {
  id: number;
  name: string;
  plate?: string | null;
  year?: number | null;
  current_mileage: number;
  created_at: string;
  updated_at: string;
}

export interface CreateVehicleRequest {
  name: string;
  plate?: string | null;
  year?: number | null;
  current_mileage: number;
}

export interface UpdateVehicleRequest {
  name: string;
  plate?: string | null;
  year?: number | null;
  current_mileage: number;
}

export interface UpdateMileageRequest {
  current_mileage: number;
}
