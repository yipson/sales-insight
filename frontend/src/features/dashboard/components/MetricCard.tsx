import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card";
import { cn } from "@/shared/utils/cn";

interface MetricCardProps {
  title: string;
  value: string | number;
  suffix?: string;
  trend?: number;
  className?: string;
  isLoading?: boolean;
}

export default function MetricCard({ title, value, suffix, trend, className, isLoading }: MetricCardProps) {
  const trendPositive = trend && trend > 0;
  const trendNegative = trend && trend < 0;

  return (
    <Card className={cn("", className)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-baseline gap-1">
          <span className="text-2xl font-bold text-foreground">
            {isLoading ? "—" : value}
          </span>
          {suffix && <span className="text-sm text-muted-foreground">{suffix}</span>}
        </div>
        {trend !== undefined && (
          <span
            className={cn(
              "mt-1 text-xs font-medium",
              trendPositive && "text-green-600",
              trendNegative && "text-red-600",
              !trendPositive && !trendNegative && "text-muted-foreground"
            )}
          >
            {trend > 0 ? "+" : ""}
            {trend.toFixed(1)}%
          </span>
        )}
      </CardContent>
    </Card>
  );
}
