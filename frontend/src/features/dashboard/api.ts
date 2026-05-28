import { useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { apiClient } from "@/core/api/client";
import type {
  SummaryResponse,
  SalesByEmployeeResponse,
  TopProductsResponse,
  CategoryCoverageResponse,
  TicketIdealResponse,
} from "./types";

const DASHBOARD_KEY = "dashboard";

function formatDateParam(date: Date): string {
  return format(date, "yyyy-MM-dd");
}

export function useDashboardSummary(restaurantId: string, from: Date, to: Date) {
  return useQuery<SummaryResponse>({
    queryKey: [DASHBOARD_KEY, "summary", restaurantId, from, to],
    queryFn: async () => {
      const { data } = await apiClient.get("/dashboard/summary", {
        params: {
          restaurant_id: restaurantId,
          from: formatDateParam(from),
          to: formatDateParam(to),
        },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useSalesByEmployee(restaurantId: string, from: Date, to: Date) {
  return useQuery<SalesByEmployeeResponse>({
    queryKey: [DASHBOARD_KEY, "sales-by-employee", restaurantId, from, to],
    queryFn: async () => {
      const { data } = await apiClient.get("/dashboard/sales-by-employee", {
        params: {
          restaurant_id: restaurantId,
          from: formatDateParam(from),
          to: formatDateParam(to),
        },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useTopProducts(restaurantId: string, from: Date, to: Date, limit = 10) {
  return useQuery<TopProductsResponse>({
    queryKey: [DASHBOARD_KEY, "top-products", restaurantId, from, to, limit],
    queryFn: async () => {
      const { data } = await apiClient.get("/dashboard/top-products", {
        params: {
          restaurant_id: restaurantId,
          from: formatDateParam(from),
          to: formatDateParam(to),
          limit,
        },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useCategoryCoverage(restaurantId: string, from: Date, to: Date) {
  return useQuery<CategoryCoverageResponse>({
    queryKey: [DASHBOARD_KEY, "category-coverage", restaurantId, from, to],
    queryFn: async () => {
      const { data } = await apiClient.get("/dashboard/category-coverage", {
        params: {
          restaurant_id: restaurantId,
          from: formatDateParam(from),
          to: formatDateParam(to),
        },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}

export function useTicketIdeal(restaurantId: string, from: Date, to: Date) {
  return useQuery<TicketIdealResponse>({
    queryKey: [DASHBOARD_KEY, "ticket-ideal", restaurantId, from, to],
    queryFn: async () => {
      const { data } = await apiClient.get("/dashboard/ticket-ideal", {
        params: {
          restaurant_id: restaurantId,
          from: formatDateParam(from),
          to: formatDateParam(to),
        },
      });
      return data;
    },
    enabled: !!restaurantId,
  });
}
