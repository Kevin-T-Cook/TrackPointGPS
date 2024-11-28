<template>
  <div class="modal-overlay" @click.self="closeModal">
    <div class="modal">
      <h2>Preferences</h2>

      <!-- Sort Order -->
      <div class="preference-section sort-order">
        <label for="sortOrder" class="sort-order-label">Sort Order:</label>
        <select id="sortOrder" v-model="localPreferences.sort_order" @change="updateSortOrder"
          class="sort-order-select">
          <option value="asc">Ascending</option>
          <option value="desc">Descending</option>
        </select>
      </div>

      <!-- Hidden Devices -->
      <div class="preference-section">
        <label>Hidden Devices:</label>
        <div v-for="device in devices" :key="device.device_id">
          <input type="checkbox" :value="device.device_id"
            :checked="localPreferences.hiddenDevices?.includes(device.device_id)"
            @change="toggleHiddenDevice(device.device_id)" />
          {{ device.display_name }}
        </div>
      </div>
      <div class="modal-actions">
        <button class="save-button" @click="savePreferences">Save</button>
        <button class="clear-button" @click="clearPreferences">Clear</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    showPreferences: {
      type: Boolean,
      required: true,
    },
    devices: {
      type: Array,
      required: true,
    },
    preferences: {
      type: Object,
      required: true,
      default: () => ({
        hiddenDevices: [],
        sort_order: "asc",
      }),
    },
  },
  data() {
    return {
      localPreferences: { ...this.preferences },
    };
  },
  mounted() {
    console.log("Modal Props:", {
      devices: this.devices,
      preferences: this.preferences,
    });
  },
  methods: {
    updateSortOrder(event) {
      const selectedValue = event.target.value;
      this.localPreferences.sort_order = selectedValue;
    },

    toggleHiddenDevice(deviceId) {
      if (!Array.isArray(this.localPreferences.hiddenDevices)) {
        this.localPreferences.hiddenDevices = [];
      }
      const hiddenDevices = [...this.localPreferences.hiddenDevices];
      const index = hiddenDevices.indexOf(deviceId);
      if (index > -1) {
        hiddenDevices.splice(index, 1);
      } else {
        hiddenDevices.push(deviceId);
      }
      this.localPreferences.hiddenDevices = hiddenDevices;
    },

    closeModal() {
      this.$emit("close-modal");
    },

    savePreferences() {
      this.$emit("update-preferences", this.localPreferences);
      this.$emit("save-preferences");
    },

    clearPreferences() {
      const defaultPreferences = {
        sort_order: "asc",
        hiddenDevices: [],
      };
      this.localPreferences = defaultPreferences;
      this.$emit("update-preferences", defaultPreferences);
      this.$emit("clear-preferences");
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
}

.preference-section {
  margin-bottom: 15px;
}

.sort-order {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 15px;
}

.sort-order-label {
  font-weight: bold;
  font-size: 16px;
}

.sort-order-select {
  padding: 8px 10px;
  font-size: 14px;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: #f9f9f9;
  transition: border-color 0.3s, background-color 0.3s;
  outline: none;
  cursor: pointer;
}

.sort-order-select:hover {
  border-color: #888;
}

.sort-order-select:focus {
  border-color: #4caf50;
  background-color: #fff;
}

.modal-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

.save-button,
.clear-button {
  padding: 10px 15px;
  background-color: #f57c00;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}

.clear-button {
  background-color: #d32f2f;
}

.save-button:hover {
  background-color: #e65100;
}

.clear-button:hover {
  background-color: #b71c1c;
}

@media (max-width: 600px) {
  .modal {
    padding: 15px;
    width: 95%;
  }

  .modal-actions {
    flex-direction: column;
  }
}
</style>