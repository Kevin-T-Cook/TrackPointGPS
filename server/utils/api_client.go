package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
    "time"
	"net/http"
	"oneStepGps/models"
)

func FetchDevices(apiKey string) ([]models.Device, error) {
	url := fmt.Sprintf("https://track.onestepgps.com/v3/api/public/device?latest_point=true&api-key=%s", apiKey)
	
	log.Println("Fetching devices from OneStepGPS API...")

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making HTTP GET request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %v", resp.Status)
		return nil, fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	log.Println("Raw API Response:", string(body))

	var apiResponse struct {
		ResultList []models.Device `json:"result_list"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling API response: %v", err)
		return nil, err
	}

	for i, device := range apiResponse.ResultList {
		if device.LatestDevicePoint.DtTracker == "" {
			apiResponse.ResultList[i].LatestDevicePoint.DtTracker = time.Now().Format(time.RFC3339)
			log.Printf("Populated DtTracker for device ID %s: %s", device.ID, apiResponse.ResultList[i].LatestDevicePoint.DtTracker)
		}
	}

	log.Printf("Processed Devices: %+v\n", apiResponse.ResultList)

	return apiResponse.ResultList, nil
}

func ProcessDevices(devices []models.Device) []models.Device {
	processedDevices := make([]models.Device, 0, len(devices))

	for _, device := range devices {
		if device.LatestDevicePoint.Latitude != 0 || device.LatestDevicePoint.Longitude != 0 {
			device.LatestDevicePoint.Speed = kmToMph(device.LatestDevicePoint.Speed)
		}
        if device.LatestDevicePoint.DtTracker == "" {
			device.LatestDevicePoint.DtTracker = time.Now().Format(time.RFC3339)
		}
		processedDevices = append(processedDevices, device)
	}

	return processedDevices
}

func kmToMph(speedKm float64) float64 {
	return speedKm * 0.621371
}