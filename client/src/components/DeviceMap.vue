/* global google */

<template>
  <div id="map" class="map-container"></div>
</template>

<script>
export default {
  name: "DeviceMap",
  props: {
    devices: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      map: null,
      markers: new Map(),
      defaultCenter: { lat: 34.0522, lng: -118.2437 },
      defaultZoom: 8,
    };
  },
  mounted() {
    this.loadGoogleMapsAPI()
      .then(() => {
        this.initMap();
      })
      .catch((error) => {
        console.error("Error loading Google Maps API:", error);
      });
  },
  watch: {
    devices: {
      handler: "updateMarkers",
      deep: true,
    },
  },
  methods: {
    loadGoogleMapsAPI() {
      return new Promise((resolve, reject) => {
        if (typeof google !== "undefined") {
          resolve();
        } else {
          const checkGoogle = setInterval(() => {
            if (typeof google !== "undefined") {
              clearInterval(checkGoogle);
              resolve();
            }
          }, 100);
          setTimeout(() => {
            clearInterval(checkGoogle);
            reject(new Error("Google Maps API failed to load."));
          }, 10000);
        }
      });
    },

    initMap() {
      this.map = new google.maps.Map(document.getElementById("map"), {
        center: this.defaultCenter,
        zoom: this.defaultZoom,
      });
    },

    updateMarkers() {
      this.devices.forEach((device) => {
        const deviceId = device.device_id;
        const latestPoint = device.latest_device_point;

        if (!latestPoint) return;

        const position = {
          lat: parseFloat(latestPoint.lat),
          lng: parseFloat(latestPoint.lng),
        };

        if (this.markers.has(deviceId)) {
          const marker = this.markers.get(deviceId);
          marker.setPosition(position);
        } else {
          const marker = new google.maps.Marker({
            position,
            map: this.map,
            title: device.display_name,
          });

          this.markers.set(deviceId, marker);

          const infoWindow = new google.maps.InfoWindow({
            content: `<strong>${device.display_name}</strong><br>Speed: ${
              latestPoint.speed || "N/A"
            } mph`,
          });

          marker.addListener("click", () => {
            infoWindow.open(this.map, marker);
          });
        }
      });

      const currentDeviceIds = new Set(this.devices.map((d) => d.device_id));
      this.markers.forEach((marker, deviceId) => {
        if (!currentDeviceIds.has(deviceId)) {
          marker.setMap(null);
          this.markers.delete(deviceId);
        }
      });
    },

    zoomToDevice(device) {
      if (device?.latest_device_point) {
        const { lat, lng } = device.latest_device_point;
        this.map.setCenter({ lat: parseFloat(lat), lng: parseFloat(lng) });
        this.map.setZoom(14);
      }
    },

    zoomToRegion(regionDevices) {
      if (!regionDevices.length) return;
      const bounds = new google.maps.LatLngBounds();
      regionDevices.forEach((device) => {
        if (device?.latest_device_point) {
          const { lat, lng } = device.latest_device_point;
          bounds.extend({ lat: parseFloat(lat), lng: parseFloat(lng) });
        }
      });
      this.map.fitBounds(bounds);
    },

    zoomOutToAllDevices() {
      if (!this.devices.length) return;
      const bounds = new google.maps.LatLngBounds();
      this.devices.forEach((device) => {
        if (device?.latest_device_point) {
          const { lat, lng } = device.latest_device_point;
          bounds.extend({ lat: parseFloat(lat), lng: parseFloat(lng) });
        }
      });
      this.map.fitBounds(bounds);
    },
  },
};
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 100%;
}

@media (max-width: 600px) {
  .map-container {
    height: 50vh;
  }
}
</style>