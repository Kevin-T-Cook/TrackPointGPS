<template>
  <div>
    <h1>Devices</h1>

    <div v-if="loading">Loading devices...</div>
    <div v-else>

      <!-- NorCal Section -->
      <div>
        <h2 @click="zoomToRegion('NorCal')">
          NorCal <span>{{ collapsedRegions.NorCal ? '+' : '-' }}</span>
        </h2>
        <ul v-if="!collapsedRegions.NorCal && filteredDevices?.NorCal.length">
          <li v-for="device in filteredDevices.NorCal" :key="device.device_id" class="device-card"
            @click="zoomToDevice(device)">
            <div class="device-icon">
              <img :src="getDeviceImage(device.device_id)" alt="Device Icon" />
            </div>
            <div class="device-info">
              <strong>Device Name:</strong> {{ device.display_name }}<br />
              <strong>Status:</strong> <span>{{ deviceStatus(device) }}</span><br />
              <strong>Speed:</strong> {{ device.speed }} mph<br />
              <strong>Last Updated:</strong> {{ formatDateTime(device.dt_tracker) }}
            </div>
          </li>
        </ul>
        <p v-if="!filteredDevices?.NorCal.length">No devices in NorCal.</p>
      </div>

      <!-- SoCal Section -->
      <div>
        <h2 @click="zoomToRegion('SoCal')">
          SoCal <span>{{ collapsedRegions.SoCal ? '+' : '-' }}</span>
        </h2>
        <ul v-if="!collapsedRegions.SoCal && filteredDevices?.SoCal.length">
          <li v-for="device in filteredDevices.SoCal" :key="device.device_id" class="device-card"
            @click="zoomToDevice(device)">
            <div class="device-icon">
              <img :src="getDeviceImage()" alt="Device Icon" />
            </div>
            <div class="device-info">
              <strong>Device Name:</strong> {{ device.display_name }}<br />
              <strong>Status:</strong> <span>{{ deviceStatus(device) }}</span><br />
              <strong>Speed:</strong> {{ device.speed }} mph<br />
              <strong>Last Updated:</strong> {{ formatDateTime(device.dt_tracker) }}
            </div>
          </li>
        </ul>
        <p v-if="!filteredDevices?.SoCal.length">No devices in SoCal.</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "DeviceList",
  props: {
    devices: {
      type: Array,
      required: true,
      default: () => [],
    },
    preferences: {
      type: Object,
      required: true,
      default: () => ({
        sort_order: "asc",
        hiddenDevices: [],
      }),
    },
  },
  data() {
    return {
      loading: true,
      collapsedRegions: {
        NorCal: false,
        SoCal: false,
      },
      defaultIconPath: "./assets/TrackPointFavicon.png",
    };
  },
  computed: {
    filteredDevices() {
      if (!this.devices || !Array.isArray(this.devices)) {
        return { NorCal: [], SoCal: [] };
      }

      const norCalDevices = this.devices.filter(
        (device) =>
          device.latitude >= 37.0 &&
          !this.preferences.hiddenDevices?.includes(device.device_id)
      );
      const soCalDevices = this.devices.filter(
        (device) =>
          device.latitude < 37.0 &&
          !this.preferences.hiddenDevices?.includes(device.device_id)
      );

      const sortOrder = this.preferences.sort_order === "asc" ? 1 : -1;
      norCalDevices.sort((a, b) => a.display_name.localeCompare(b.display_name) * sortOrder);
      soCalDevices.sort((a, b) => a.display_name.localeCompare(b.display_name) * sortOrder);

      return {
        NorCal: norCalDevices,
        SoCal: soCalDevices,
      };
    },
  },

  watch: {
    devices: {
      handler(newDevices) {
        console.log("Updated Devices in DeviceList:", newDevices);
        this.loading = false;
      },
      immediate: true,
    },
  },

  methods: {
    getDeviceImage() {
      return require("@/assets/TrackPointFavicon.png");
    },

    toggleRegion(region) {
      this.collapsedRegions[region] = !this.collapsedRegions[region];
    },

    deviceStatus(device) {
      return parseFloat(device.speed) > 0 ? "Active" : "Inactive";
    },

    formatDateTime(timestamp) {
      if (!timestamp) return "N/A";
      const date = new Date(timestamp);
      const options = { year: "numeric", month: "short", day: "numeric" };
      const dateString = date.toLocaleDateString("en-US", options);
      const timeString = date.toLocaleTimeString("en-US", {
        hour: "numeric",
        minute: "2-digit",
      });
      return `${dateString}, ${timeString}`;
    },

    zoomToRegion(region) {
      const devices = this.filteredDevices[region];
      this.$emit("zoom-to-region", devices);
    },

    zoomToDevice(device) {
      this.$emit("zoom-to-device", device);
    },

    zoomOutToAllDevices() {
      this.$emit("zoom-out-to-all");
    },
  },
};
</script>


<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}

li.device-card {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  background-color: #faf5eb;
  width: 400px;
}

li.device-card:hover {
  background-color: #f5e6d8;
}

.device-icon {
  flex-shrink: 0;
  width: 50px;
  height: 50px;
  margin-right: 15px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
}

.device-icon img {
  max-width: 100%;
  max-height: 100%;
}

.device-info {
  flex: 1;
  line-height: 1.5;
}

strong {
  font-weight: bold;
}

h2 {
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>