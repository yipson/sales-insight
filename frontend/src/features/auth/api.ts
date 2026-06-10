import { useQuery, useMutation } from "@tanstack/react-query";
import { apiClient } from "@/core/api/client";
import type { AuthStatus, BootstrapRequest, BootstrapResponse, RevokeRequest } from "./types";

const AUTH_KEY = "auth";

export function useAuthStatus() {
  return useQuery<AuthStatus>({
    queryKey: [AUTH_KEY, "status"],
    queryFn: async () => {
      const { data } = await apiClient.get("/auth/status");
      const merchants = data.merchants || [];
      const connected = merchants.length > 0;
      return {
        connected,
        merchantId: merchants[0]?.id,
        name: merchants[0]?.name,
        cloverEnv: merchants[0]?.clover_env,
      };
    },
    retry: false,
  });
}

export function useBootstrap() {
  return useMutation<BootstrapResponse, Error, BootstrapRequest>({
    mutationFn: async (payload) => {
      const { data } = await apiClient.post("/auth/bootstrap", payload);
      return data;
    },
  });
}

export function useRevoke() {
  return useMutation<void, Error, RevokeRequest>({
    mutationFn: async (payload) => {
      await apiClient.post("/auth/revoke", payload);
    },
  });
}
