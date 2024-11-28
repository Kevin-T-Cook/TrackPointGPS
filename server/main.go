package main

import (
	"log"
	"net/http"
	"os"
	"oneStepGps/handlers"
	"oneStepGps/models"
	"oneStepGps/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Upgrade, Connection")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if r.Header.Get("Upgrade") == "websocket" {
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Upgrade, Connection, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions")
        }

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable TimeZone=America/Los_Angeles"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to the database successfully")

	if !db.Migrator().HasTable(&models.Preferences{}) {
		log.Println("Creating preferences table...")
		if err := db.AutoMigrate(&models.Preferences{}); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		} else {
			log.Println("Preferences table created successfully.")
		}
	} else {
		log.Println("Preferences table already exists. Skipping migration.")
	}

	deviceHandler := handlers.DeviceHandler{
		DB:     db,
		APIKey: os.Getenv("API_KEY"),
	}
	prefsHandler := handlers.PreferencesHandler{DB: db}
	loginHandler := handlers.LoginHandler{DB: db}

	mux := http.NewServeMux()
    mux.Handle("/api/devices", utils.AuthMiddleware(http.HandlerFunc(deviceHandler.GetDevicesWithPreferences)))
	mux.Handle("/api/preferences/save", utils.AuthMiddleware(http.HandlerFunc(prefsHandler.SavePreferences)))
	mux.Handle("/api/preferences/get", utils.AuthMiddleware(http.HandlerFunc(prefsHandler.GetPreferences)))
	mux.Handle("/api/preferences/delete", utils.AuthMiddleware(http.HandlerFunc(prefsHandler.DeletePreferences)))
	mux.Handle("/api/devices/realtime", utils.WebSocketAuthMiddleware(http.HandlerFunc(deviceHandler.RealTimeUpdates)))
    mux.HandleFunc("/api/auth/login", loginHandler.Login)

	wrappedMux := enableCORS(mux)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}