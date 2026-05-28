import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useBootstrap } from "../api";
import { useAuthStore } from "@/core/store/authStore";
import { Input } from "@/shared/components/ui/input";
import { Button } from "@/shared/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/shared/components/ui/card";

const schema = z.object({
  name: z.string().min(1, "Name is required"),
  cloverMerchantId: z.string().min(1, "Clover Merchant ID is required"),
  accessToken: z.string().min(1, "Access token is required"),
  refreshToken: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export default function LoginForm() {
  const login = useAuthStore((state) => state.login);
  const bootstrap = useBootstrap();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    try {
      const response = await bootstrap.mutateAsync({
        name: data.name,
        clover_merchant_id: data.cloverMerchantId,
        access_token: data.accessToken,
        refresh_token: data.refreshToken,
      });
      // Use the returned merchant ID or a dummy token for now
      login(data.accessToken, response.id);
    } catch (error) {
      console.error("Bootstrap failed:", error);
    }
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Sales Insight</CardTitle>
        <CardDescription>Connect your Clover account</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <Input placeholder="Restaurant Name" {...register("name")} />
            {errors.name && <p className="mt-1 text-xs text-destructive">{errors.name.message}</p>}
          </div>
          <div>
            <Input placeholder="Clover Merchant ID" {...register("cloverMerchantId")} />
            {errors.cloverMerchantId && <p className="mt-1 text-xs text-destructive">{errors.cloverMerchantId.message}</p>}
          </div>
          <div>
            <Input placeholder="Access Token" type="password" {...register("accessToken")} />
            {errors.accessToken && <p className="mt-1 text-xs text-destructive">{errors.accessToken.message}</p>}
          </div>
          <div>
            <Input placeholder="Refresh Token (optional)" type="password" {...register("refreshToken")} />
          </div>
          <Button type="submit" className="w-full" disabled={isSubmitting || bootstrap.isPending}>
            {bootstrap.isPending ? "Connecting..." : "Connect"}
          </Button>
          {bootstrap.isError && (
            <p className="text-xs text-destructive">Failed to connect. Please check your credentials.</p>
          )}
        </form>
      </CardContent>
    </Card>
  );
}
