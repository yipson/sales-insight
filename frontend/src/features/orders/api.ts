import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/core/api/client";

export interface Order {
  id: string;
  clover_order_id: string;
  restaurant_id: string;
  employee_id?: string;
  employee_name?: string;
  total: number;
  tax: number;
  tip?: number;
  discount?: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
}

export interface OrderDetail extends Order {
  items: OrderItem[];
}

export function useOrders() {
  return useQuery<Order[]>({
    queryKey: ["orders"],
    queryFn: async () => {
      const { data } = await apiClient.get("/orders");
      return data;
    },
  });
}

export function useOrder(id: string) {
  return useQuery<OrderDetail>({
    queryKey: ["orders", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/orders/${id}`);
      return data;
    },
    enabled: !!id,
  });
}
