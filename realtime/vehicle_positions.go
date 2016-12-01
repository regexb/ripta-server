package realtime

type Vehicle struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	LicencePlate string `json:"license_plate"`
}

type VehiclePositionsService service
