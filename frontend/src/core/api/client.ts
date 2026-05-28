import axios from "axios";
import { useAuthStore } from "@/core/store/authStore";
import type { ApiError } from "./types";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000,
});

// Request interceptor: inject JWT token
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor: handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;

      if (status === 401) {
        // Token expired or invalid: clear auth and reload
        useAuthStore.getState().logout();
        window.location.href = "/login";
        return Promise.reject(error);
      }

      if (status >= 500) {
        console.error("Server error:", data);
      }

      const apiError: ApiError = {
        message: data?.message || data?.error || "An unexpected error occurred",
        status,
        details: data?.details,
      };

      return Promise.reject(apiError);
    }

    // Network error
    const networkError: ApiError = {
      message: error.message || "Network error. Please check your connection.",
    };
    return Promise.reject(networkError);
  }
);
