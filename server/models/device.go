package models

type LatestDevicePoint struct {
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
	Speed      float64 `json:"speed"`
	Unit       string  `json:"unit"`
	DtTracker  string  `json:"dt_tracker"`
}

type DeviceState struct {
	DriveStatus      string    `json:"drive_status"`
	FuelPercent      *int      `json:"fuel_percent"`
	SoftwareOdometer *Odometer `json:"software_odometer"`
}

type Odometer struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Device struct {
    ID                string             `json:"device_id"`
    Name              string             `json:"display_name"`
    LatestDevicePoint LatestDevicePoint  `json:"latest_device_point"`
    State             DeviceState        `json:"device_state"`
    ActiveState       string             `json:"active_state"`
    UpdatedAt         string             `json:"updated_at"`
}