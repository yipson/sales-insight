import { useQuery, useMutation } from "@tanstack/react-query";
import { apiClient } from "@/core/api/client";

export interface SyncLog {
  id: string;
  restaurant_id: string;
  entity: string;
  status: string;
  records_processed: number;
  cursor_from?: string;
  cursor_to?: string;
  triggered_by: string;
  created_at: string;
}

export interface SyncError {
  id: string;
  restaurant_id: string;
  entity: string;
  message: string;
  resolved: boolean;
  created_at: string;
}

export function useSyncStatus(restaurantId: string) {
  return useQuery<SyncLog[]>({
    queryKey: ["sync", "status", restaurantId],
    queryFn: async () => {
      const { data } = await apiClient.get("/sync/status", {
        params: { restaurant_id: restaurantId },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useSyncLogs(restaurantId: string, limit = 50) {
  return useQuery<SyncLog[]>({
    queryKey: ["sync", "logs", restaurantId, limit],
    queryFn: async () => {
      const { data } = await apiClient.get("/sync/logs", {
        params: { restaurant_id: restaurantId, limit },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useSyncTrigger() {
  return useMutation<void, Error, { restaurant_id: string; entity: string }>({
    mutationFn: async (payload) => {
      await apiClient.post("/sync/trigger", payload);
    },
  });
}

export function useSyncBackfill() {
  return useMutation<void, Error, { restaurant_id: string; entity: string; from: string }>({
    mutationFn: async (payload) => {
      await apiClient.post("/sync/backfill", payload);
    },
  });
}
