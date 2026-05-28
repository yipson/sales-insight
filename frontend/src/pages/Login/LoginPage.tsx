import { useAuthStore } from "@/core/store/authStore";

export default function LoginPage() {
  const login = useAuthStore((state) => state.login);

  const handleLogin = () => {
    // Temporary bypass for dev; real login uses auth feature
    login("dev-token", "00000000-0000-0000-0000-000000000000");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-sm rounded-lg border bg-card p-8 shadow-sm">
        <h1 className="text-2xl font-bold text-primary">Sales Insight</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Sign in to your account
        </p>
        <button
          onClick={handleLogin}
          className="mt-6 w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          Login (dev)
        </button>
      </div>
    </div>
  );
}
