package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"oneStepGps/models"
	"oneStepGps/storage"
	"oneStepGps/utils"
)

type DeviceHandler struct {
	DB     *gorm.DB
	APIKey string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *DeviceHandler) RealTimeUpdates(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		devices, err := utils.FetchDevices(h.APIKey)
		if err != nil {
			log.Printf("Error fetching devices: %v", err)
			break
		}

		updatedDevices := h.prepareDeviceResponse(devices)

		if err := conn.WriteJSON(updatedDevices); err != nil {
			log.Printf("Error writing WebSocket data: %v", err)
			break
		}

		log.Println("Sent updated devices to client")
		time.Sleep(5 * time.Second)
	}

	log.Println("WebSocket connection closed")
}

func (h *DeviceHandler) GetDevicesWithPreferences(w http.ResponseWriter, r *http.Request) {
    devices, err := utils.FetchDevices(h.APIKey)
    if err != nil {
        log.Printf("Error fetching devices: %v", err)
        http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
        return
    }

    userId, ok := r.Context().Value("userID").(string)
    if !ok || userId == "" {
        log.Printf("User ID not found in context")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    prefs, err := storage.GetUserPreferences(h.DB, userId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            log.Printf("No preferences found for user %s, returning all devices", userId)
            response := h.prepareDeviceResponse(devices)
            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(response); err != nil {
                log.Printf("Error encoding devices to JSON: %v", err)
                http.Error(w, "Failed to encode devices", http.StatusInternalServerError)
            }
            return
        }
        log.Printf("Error fetching preferences for user %s: %v", userId, err)
        http.Error(w, "Failed to fetch preferences", http.StatusInternalServerError)
        return
    }

    filteredDevices := h.applyPreferences(devices, prefs)

    response := h.prepareDeviceResponse(filteredDevices)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding filtered devices to JSON: %v", err)
        http.Error(w, "Failed to encode devices", http.StatusInternalServerError)
    }
}

func (h *DeviceHandler) prepareDeviceResponse(devices []models.Device) []map[string]interface{} {
	var response []map[string]interface{}

	for _, device := range devices {
		deviceData := map[string]interface{}{
			"device_id":    device.ID,
			"display_name": device.Name,
			"latest_device_point": map[string]interface{}{
				"lat":   device.LatestDevicePoint.Latitude,
				"lng":   device.LatestDevicePoint.Longitude,
				"speed": formatSpeedInMph(device.LatestDevicePoint.Speed),
				"dt_tracker": device.LatestDevicePoint.DtTracker,
			},
			"device_state": map[string]interface{}{
				"drive_status": device.State.DriveStatus,
			},
			"fuel_level":   device.State.FuelPercent,
			"last_updated": device.UpdatedAt,
		}

		response = append(response, deviceData)
	}
	return response
}

func formatSpeedInMph(speed float64) string {
	return fmt.Sprintf("%.2f mph", speed)
}

func (h *DeviceHandler) applyPreferences(devices []models.Device, prefs models.Preferences) []models.Device {
	var hiddenDevices []string
	if err := json.Unmarshal(prefs.HiddenDevices, &hiddenDevices); err != nil {
		log.Printf("Error unmarshalling hidden_devices: %v", err)
		return devices
	}

	hiddenDevicesMap := make(map[string]bool)
	for _, id := range hiddenDevices {
		hiddenDevicesMap[id] = true
	}

	filteredDevices := []models.Device{}
	for _, device := range devices {
		if !hiddenDevicesMap[device.ID] {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return filteredDevices
}