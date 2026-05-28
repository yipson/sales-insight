import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card";
import { Progress } from "@/shared/components/ui/progress";
import type { TicketIdealResponse } from "../types";

interface TicketIdealCardProps {
  data?: TicketIdealResponse;
  isLoading?: boolean;
}

export default function TicketIdealCard({ data, isLoading }: TicketIdealCardProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Ticket Ideal</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-muted-foreground">Loading...</div>
        </CardContent>
      </Card>
    );
  }

  const result = data?.result;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Ticket Ideal Analysis</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {result ? (
          <>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Completion Rate</span>
              <span className="text-lg font-bold">
                {(result.completion_rate * 100).toFixed(1)}%
              </span>
            </div>
            <Progress value={result.completion_rate * 100} className="h-2" />
            <div className="grid grid-cols-3 gap-4 pt-2">
              <div className="text-center">
                <div className="text-2xl font-bold">{result.total_orders}</div>
                <div className="text-xs text-muted-foreground">Total</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">
                  {result.complete_orders}
                </div>
                <div className="text-xs text-muted-foreground">Complete</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-red-600">
                  {result.incomplete_orders}
                </div>
                <div className="text-xs text-muted-foreground">Incomplete</div>
              </div>
            </div>
          </>
        ) : (
          <div className="text-muted-foreground">No data available</div>
        )}
      </CardContent>
    </Card>
  );
}
