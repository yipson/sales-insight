import { useAuthStore } from "@/core/store/authStore";
import { useSyncStatus } from "@/features/sync/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card";
import { Badge } from "@/shared/components/ui/badge";
import { formatDate } from "@/shared/utils/formatters";

export default function SyncStatusCard() {
  const restaurantId = useAuthStore((state) => state.restaurantId);
  const { data: logs, isLoading } = useSyncStatus(restaurantId || "");

  const entities = ["orders", "items", "employees", "payments"];

  const getStatusForEntity = (entity: string) => {
    if (!logs) return "unknown";
    const log = logs.find((l: { entity: string }) => l.entity === entity);
    if (!log) return "unknown";
    return log.status;
  };

  const getLastSyncForEntity = (entity: string) => {
    if (!logs) return null;
    const log = logs.find((l: { entity: string }) => l.entity === entity);
    return log?.created_at;
  };

  const statusVariant = (status: string) => {
    switch (status) {
      case "success": return "default" as const;
      case "failed": return "destructive" as const;
      case "in_progress": return "secondary" as const;
      default: return "outline" as const;
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Sync Status</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="text-muted-foreground">Loading...</div>
        ) : (
          <div className="space-y-3">
            {entities.map((entity) => {
              const status = getStatusForEntity(entity);
              const lastSync = getLastSyncForEntity(entity);
              return (
                <div key={entity} className="flex items-center justify-between">
                  <span className="capitalize">{entity}</span>
                  <div className="flex items-center gap-2">
                    {lastSync && (
                      <span className="text-xs text-muted-foreground">
                        {formatDate(lastSync)}
                      </span>
                    )}
                    <Badge variant={statusVariant(status)}>
                      {status}
                    </Badge>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
