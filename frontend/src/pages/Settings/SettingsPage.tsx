import SyncStatusCard from "@/features/sync/components/SyncStatusCard";
import SyncTriggerButton from "@/features/sync/components/SyncTriggerButton";
import SyncLogTable from "@/features/sync/components/SyncLogTable";
import { useRevoke } from "@/features/auth/api";
import { Button } from "@/shared/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card";

export default function SettingsPage() {
  const revoke = useRevoke();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <SyncStatusCard />
        <Card>
          <CardHeader>
            <CardTitle>Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="mb-2 text-sm text-muted-foreground">Manual sync</p>
              <SyncTriggerButton />
            </div>
            <div>
              <p className="mb-2 text-sm text-muted-foreground">Connection</p>
              <Button
                variant="destructive"
                onClick={() => revoke.mutate({ merchant_id: "dummy" })}
                disabled={revoke.isPending}
              >
                {revoke.isPending ? "Revoking..." : "Revoke Connection"}
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      <div>
        <h2 className="mb-4 text-lg font-semibold">Sync Logs</h2>
        <SyncLogTable />
      </div>
    </div>
  );
}
