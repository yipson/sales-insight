import { create } from "zustand";

interface AuthState {
  token: string | null;
  restaurantId: string | null;
  isAuthenticated: boolean;
  login: (token: string, restaurantId: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  restaurantId: null,
  isAuthenticated: false,
  login: (token, restaurantId) => set({ token, restaurantId, isAuthenticated: true }),
  logout: () => set({ token: null, restaurantId: null, isAuthenticated: false }),
}));
