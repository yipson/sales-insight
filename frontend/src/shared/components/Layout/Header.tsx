import { useAuthStore } from "@/core/store/authStore";
import { useThemeStore } from "@/core/store/themeStore";
import { useLocation } from "react-router-dom";
import { LogOut, Moon, Sun } from "lucide-react";

function getPageTitle(path: string) {
  switch (path) {
    case "/dashboard": return "Dashboard";
    case "/employees": return "Employees";
    case "/products": return "Products";
    case "/orders": return "Orders";
    case "/settings": return "Settings";
    default: return "";
  }
}

export default function Header() {
  const { pathname } = useLocation();
  const { theme, toggleTheme } = useThemeStore();
  const logout = useAuthStore((state) => state.logout);

  return (
    <header className="flex h-16 items-center justify-between border-b bg-card px-6">
      <h2 className="text-xl font-semibold text-foreground">
        {getPageTitle(pathname)}
      </h2>
      <div className="flex items-center gap-2">
        <button
          onClick={toggleTheme}
          className="rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          title="Toggle theme"
        >
          {theme === "light" ? <Moon size={18} /> : <Sun size={18} />}
        </button>
        <button
          onClick={logout}
          className="rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          title="Logout"
        >
          <LogOut size={18} />
        </button>
      </div>
    </header>
  );
}
