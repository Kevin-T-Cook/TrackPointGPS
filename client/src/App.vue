<template>
  <div id="app">
    <div v-if="isAuthenticated">
      <header class="header">
        <div class="logo-container">
          <img src="./assets/TrackPoint.png" alt="Logo" class="logo" />
        </div>
        <button class="preferences-button" @click="togglePreferences">
          Preferences
        </button>
        <button class="zoom-out-button" @click="zoomOutToAllDevices">
          Show All Devices
        </button>
      </header>
      <div class="main-container">
        <div class="sidebar">
          <DeviceList :devices="devices" :preferences="preferences" @zoom-to-device="zoomToDevice"
            @zoom-to-region="zoomToRegion" />
        </div>
        <div class="map-container">
          <DeviceMap :devices="devices" ref="deviceMap" />
        </div>
      </div>
      <PreferencesModal v-if="showPreferences" :showPreferences="showPreferences" :preferences="preferences"
        :devices="devices" @update-preferences="updatePreferences" @close-modal="closePreferencesModal"
        @save-preferences="savePreferences" @clear-preferences="clearPreferences" />
    </div>
    <div v-else>
      <router-view />
    </div>
  </div>
</template>

<script>
import DeviceList from "./components/DeviceList.vue";
import DeviceMap from "./components/DeviceMap.vue";
import PreferencesModal from "./components/PreferencesModal.vue";
import api from "./api/api";

export default {
  name: "App",
  components: {
    DeviceList,
    DeviceMap,
    PreferencesModal,
  },
  data() {
    return {
      isAuthenticated: false,
      showPreferences: false,
      devices: [],
      preferences: {
        sort_order: "asc",
        hiddenDevices: [],
      },
      websocket: null,
      isLoading: true,
      error: null,
    };
  },
  async created() {
    this.isAuthenticated = !!localStorage.getItem("token");

    if (this.isAuthenticated) {
      try {
        await this.initializeData();
      } catch (err) {
        console.error("Error during initialization:", err);
        this.error = "Failed to initialize application data";
      } finally {
        this.isLoading = false;
      }
    }
  },
  beforeUnmount() {
    if (this.websocket) {
      this.websocket.close();
    }
  },
  methods: {
    async initializeData() {
      try {
        this.preferences = await api.fetchPreferences().catch(() => ({
          sort_order: "asc",
          hiddenDevices: [],
        }));

        const devices = await api.fetchDevices();
        this.devices = this.processDevices(devices);

        this.initWebSocket();
      } catch (error) {
        console.error("Error initializing data:", error);
        throw error;
      }
    },

    processDevices(devices) {
      return devices.map((device) => ({
        ...device,
        speed: device.latest_device_point?.speed
          ? Math.round(parseFloat(device.latest_device_point.speed) * 0.621371)
          : "0",
        latitude: device.latest_device_point?.lat,
        longitude: device.latest_device_point?.lng,
        dt_tracker: device.latest_device_point?.dt_tracker,
      }));
    },

    applyPreferences() {
      console.log("Applying preferences:", this.preferences);
      let filteredDevices = [...this.devices];

      if (this.preferences.hiddenDevices?.length) {
        filteredDevices = filteredDevices.filter(
          (device) => !this.preferences.hiddenDevices.includes(device.device_id)
        );
      }

      filteredDevices.sort((a, b) => {
        const comparison = a.display_name.localeCompare(b.display_name);
        return this.preferences.sort_order === "asc" ? comparison : -comparison;
      });
      this.devices = filteredDevices;
    },

    async updatePreferences(updatedPreferences) {
      console.log("Updating preferences:", updatedPreferences);
      this.preferences = { ...updatedPreferences };
    },

    togglePreferences() {
      this.showPreferences = true;
    },

    closePreferencesModal() {
      this.showPreferences = false;
    },

    async savePreferences() {
      try {
        console.log("Saving preferences:", this.preferences);
        const transformedPreferences = {
          sort_order: this.preferences.sort_order,
          hidden_devices: this.preferences.hiddenDevices,
        };

        await api.savePreferences(transformedPreferences);
        this.showPreferences = false;
      } catch (err) {
        console.error("Error saving preferences:", err);
        alert("Failed to save preferences. Please try again.");
      }
    },

    async clearPreferences() {
      try {
        const defaultPreferences = {
          sort_order: "asc",
          hiddenDevices: [],
        };

        await api.savePreferences(defaultPreferences);
        this.preferences = { ...defaultPreferences };
      } catch (err) {
        console.error("Error clearing preferences:", err);
        alert("Failed to clear preferences. Please try again.");
      }
    },

    initWebSocket() {
      const token = localStorage.getItem("token");
      if (!token) {
        console.error("Token is missing. Unable to initialize WebSocket.");
        return;
      }

      const websocketUrl = `ws://localhost:8080/api/devices/realtime?token=${token}`;
      console.log("Connecting to WebSocket:", websocketUrl);

      this.websocket = new WebSocket(websocketUrl);

      this.websocket.onopen = () => {
        console.log("WebSocket connection established.");
      };

      this.websocket.onmessage = (event) => {
        try {
          const updatedDevices = JSON.parse(event.data);
          console.log("WebSocket Data Received:", updatedDevices);
          this.devices = this.processDevices(updatedDevices);
        } catch (error) {
          console.error("Error processing WebSocket message:", error);
        }
      };

      this.websocket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      this.websocket.onclose = () => {
        console.warn("WebSocket connection closed. Reconnecting...");
        setTimeout(this.initWebSocket, 5000);
      };
    },

    zoomToDevice(device) {
      this.$refs.deviceMap.zoomToDevice(device);
    },

    zoomToRegion(regionDevices) {
      this.$refs.deviceMap.zoomToRegion(regionDevices);
    },

    zoomOutToAllDevices() {
      this.$refs.deviceMap.zoomOutToAllDevices();
    },
  },
};
</script>

<style>
* {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow-x: hidden;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
}

.main-container {
  display: flex;
  flex: 1;
  height: calc(100vh - 80px);
  overflow: hidden;
}

.sidebar {
  width: 29%;
  max-width: 29%;
  padding: 10px;
  background-color: #f4f4f4;
  overflow-y: auto;
}

.map-container {
  flex: 1;
  position: relative;
  height: 100%;
  overflow: hidden;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #ddd;
  height: 80px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-sizing: border-box;
}

.logo-container {
  display: flex;
  align-items: center;
}

.logo {
  height: 60px;
}

.preferences-button,
.zoom-out-button {
  padding: 10px 20px;
  background-color: #f57c00;
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 14px;
  cursor: pointer;
}

.preferences-button:hover,
.zoom-out-button:hover {
  background-color: #e65100;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    height: auto;
    padding: 10px;
  }

  .logo {
    height: 40px;
    margin-bottom: 10px;
  }

  .preferences-button,
  .zoom-out-button {
    width: 100%;
    margin-bottom: 10px;
    font-size: 12px;
  }

  .main-container {
    flex-direction: column;
    height: auto;
  }

  .sidebar {
    width: 100%;
    max-width: 100%;
    height: auto;
    margin-bottom: 10px;
  }

  .map-container {
    height: 400px;
  }
}

@media (max-width: 480px) {
  .logo {
    height: 30px;
  }

  .preferences-button,
  .zoom-out-button {
    font-size: 10px;
    padding: 8px 15px;
  }

  .map-container {
    height: 300px;
  }
}
</style>