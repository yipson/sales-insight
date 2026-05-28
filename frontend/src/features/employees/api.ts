import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/core/api/client";

export interface Employee {
  id: string;
  name: string;
  role: string;
  clover_employee_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export function useEmployees() {
  return useQuery<Employee[]>({
    queryKey: ["employees"],
    queryFn: async () => {
      const { data } = await apiClient.get("/employees");
      return data;
    },
  });
}

export function useEmployee(id: string) {
  return useQuery<Employee>({
    queryKey: ["employees", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/employees/${id}`);
      return data;
    },
    enabled: !!id,
  });
}
