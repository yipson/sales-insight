import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "@/core/store/authStore";

const SKIP_AUTH = import.meta.env.VITE_SKIP_AUTH === "true";

export default function ProtectedRoute() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated && !SKIP_AUTH) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
