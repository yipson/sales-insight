import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/core/api/client";

export interface Product {
  id: string;
  name: string;
  price: number;
  category_id?: string;
  category_name?: string;
  clover_item_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  name: string;
  clover_category_id: string;
  created_at: string;
}

export function useProducts() {
  return useQuery<Product[]>({
    queryKey: ["products"],
    queryFn: async () => {
      const { data } = await apiClient.get("/products");
      return data;
    },
  });
}

export function useProduct(id: string) {
  return useQuery<Product>({
    queryKey: ["products", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/products/${id}`);
      return data;
    },
    enabled: !!id,
  });
}

export function useCategories() {
  return useQuery<Category[]>({
    queryKey: ["categories"],
    queryFn: async () => {
      const { data } = await apiClient.get("/products/categories");
      return data;
    },
  });
}
