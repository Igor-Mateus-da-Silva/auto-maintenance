package model

type Vehicle struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Plate          *string `json:"plate"`
	Year           *int    `json:"year"`
	CurrentMileage int     `json:"current_mileage"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type CreateVehicleRequest struct {
	Name           string  `json:"name"`
	Plate          *string `json:"plate"`
	Year           *int    `json:"year"`
	CurrentMileage int     `json:"current_mileage"`
}

type UpdateVehicleRequest struct {
	Name           string  `json:"name"`
	Plate          *string `json:"plate"`
	Year           *int    `json:"year"`
	CurrentMileage int     `json:"current_mileage"`
}

type UpdateMileageRequest struct {
	CurrentMileage int `json:"current_mileage"`
}
