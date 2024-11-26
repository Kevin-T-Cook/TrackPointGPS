package storage

import (
	"encoding/json"
	"log"
	"gorm.io/gorm"
	"oneStepGps/models"
)

func SavePreferences(db *gorm.DB, prefs models.Preferences) error {
	log.Printf("Preferences before save (raw struct): %+v", prefs)

	hiddenDevicesJSON, err := json.Marshal(prefs.HiddenDevices)
	if err != nil {
		log.Printf("Error serializing HiddenDevices: %v", err)
		return err
	}

	log.Printf("Serialized HiddenDevices (JSON): %s", hiddenDevicesJSON)

	query := `
		INSERT INTO preferences (user_id, sort_order, hidden_devices)
		VALUES ($1, $2, $3::jsonb)
		ON CONFLICT (user_id) DO UPDATE 
		SET sort_order = EXCLUDED.sort_order,
		    hidden_devices = EXCLUDED.hidden_devices
		RETURNING id
	`

	err = db.Exec(query, prefs.UserID, prefs.SortOrder, string(hiddenDevicesJSON)).Error
	if err != nil {
		log.Printf("Error executing raw SQL insert/update: %v", err)
		return err
	}

	log.Println("Preferences saved successfully")
	return nil
}

func GetPreferences(db *gorm.DB) ([]models.Preferences, error) {
	var prefs []models.Preferences

	err := db.Find(&prefs).Error
	if err != nil {
		log.Printf("Error fetching preferences: %v", err)
		return nil, err
	}

	log.Printf("Fetched Preferences: %+v", prefs)
	return prefs, nil
}

func DeletePreferences(db *gorm.DB, id string) error {
	log.Printf("Deleting preferences with ID: %s", id)

	err := db.Delete(&models.Preferences{}, id).Error
	if err != nil {
		log.Printf("Error deleting preferences: %v", err)
		return err
	}

	log.Println("Preferences deleted successfully")
	return nil
}

func GetUserPreferences(db *gorm.DB, userID string) (models.Preferences, error) {
	var prefs models.Preferences

	log.Printf("Fetching preferences for UserID: %s", userID)

	err := db.Where("user_id = ?", userID).First(&prefs).Error
	if err != nil {
		log.Printf("Error fetching preferences for UserID %s: %v", userID, err)
		return prefs, err
	}

	if len(prefs.HiddenDevices) > 0 {
		var hiddenDevices []string
		if err := json.Unmarshal(prefs.HiddenDevices, &hiddenDevices); err != nil {
			log.Printf("Error unmarshalling hidden_devices for UserID %s: %v", userID, err)
		}
		log.Printf("Unmarshalled hidden_devices: %+v", hiddenDevices)
	}

	log.Printf("Fetched Preferences for UserID %s: %+v", userID, prefs)
	return prefs, nil
}

func DeletePreferencesByUserID(db *gorm.DB, userID string) error {
	log.Printf("Deleting preferences for UserID: %s", userID)
	err := db.Where("user_id = ?", userID).Delete(&models.Preferences{}).Error
	if err != nil {
		log.Printf("Error deleting preferences for UserID %s: %v", userID, err)
		return err
	}
	log.Printf("Preferences deleted successfully for UserID: %s", userID)
	return nil
}