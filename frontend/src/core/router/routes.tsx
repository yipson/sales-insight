import { type RouteObject } from "react-router-dom";
import AppRouter from "./AppRouter";
import ProtectedRoute from "./ProtectedRoute";
import Layout from "./Layout";
import LoginPage from "@/pages/Login/LoginPage";
import DashboardPage from "@/pages/Dashboard/DashboardPage";
import EmployeesPage from "@/pages/Employees/EmployeesPage";
import ProductsPage from "@/pages/Products/ProductsPage";
import OrdersPage from "@/pages/Orders/OrdersPage";
import SettingsPage from "@/pages/Settings/SettingsPage";

export const routes: RouteObject[] = [
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <Layout />,
        children: [
          { path: "/dashboard", element: <DashboardPage /> },
          { path: "/employees", element: <EmployeesPage /> },
          { path: "/products", element: <ProductsPage /> },
          { path: "/orders", element: <OrdersPage /> },
          { path: "/settings", element: <SettingsPage /> },
        ],
      },
    ],
  },
];

export { AppRouter };
