import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const login = async (credentials) => {
  try {
    const response = await apiClient.post("/auth/login", credentials);
    return response.data;
  } catch (error) {
    console.error("Error during login:", error.response?.data || error.message);
    throw error;
  }
};

export const fetchDevices = async () => {
  try {
    const response = await apiClient.get("/devices");
    return response.data;
  } catch (error) {
    console.error("Error fetching devices:", error);
    throw error;
  }
};

export const fetchPreferences = async () => {
  try {
    const response = await apiClient.get("/preferences/get");
    return response.data;
  } catch (error) {
    console.error("Error fetching preferences:", error);
    throw error;
  }
};

export const savePreferences = async (preferences) => {
  try {
    console.log("Preferences to save:", preferences);
    const response = await apiClient.post("/preferences/save", preferences);
    return response.data;
  } catch (error) {
    console.error("Error saving preferences:", error.response?.data || error.message);
    throw error;
  }
};

export const deletePreferences = async (id) => {
  try {
    const response = await apiClient.delete(`/preferences/delete?id=${id}`);
    return response.data;
  } catch (error) {
    console.error("Error deleting preferences:", error);
    throw error;
  }
};

export default {
  login,
  fetchDevices,
  fetchPreferences,
  savePreferences,
  deletePreferences,
};