package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"oneStepGps/models"
	"oneStepGps/storage"
	"gorm.io/gorm"
)

type PreferencesHandler struct {
	DB *gorm.DB
}

func (h *PreferencesHandler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	var prefs models.Preferences

	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	prefs.UserID = userID

	log.Printf("Received Preferences for user %s: %+v\n", userID, prefs)

	if err := storage.SavePreferences(h.DB, prefs); err != nil {
		log.Printf("Error saving preferences for user %s: %v", userID, err)
		http.Error(w, "Failed to save preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Preferences saved successfully"))
}

func (h *PreferencesHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	prefs, err := storage.GetUserPreferences(h.DB, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No preferences found for user %s", userID)
			http.Error(w, "No preferences found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching preferences for user %s: %v", userID, err)
		http.Error(w, "Failed to fetch preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prefs)
}

func (h *PreferencesHandler) DeletePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	if err := storage.DeletePreferencesByUserID(h.DB, userID); err != nil {
		log.Printf("Error deleting preferences for user %s: %v", userID, err)
		http.Error(w, "Failed to delete preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Preferences cleared successfully"))
}