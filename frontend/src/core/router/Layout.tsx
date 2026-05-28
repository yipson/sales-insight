import { Outlet } from "react-router-dom";
import Sidebar from "@/shared/components/Layout/Sidebar";
import Header from "@/shared/components/Layout/Header";

export default function Layout() {
  return (
    <div className="flex min-h-screen bg-background">
      <Sidebar />
      <div className="flex flex-1 flex-col">
        <Header />
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
