import { useAuthStore } from "@/core/store/authStore";
import { useSyncLogs } from "@/features/sync/api";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/components/ui/table";
import { Badge } from "@/shared/components/ui/badge";
import { formatDate } from "@/shared/utils/formatters";

export default function SyncLogTable() {
  const restaurantId = useAuthStore((state) => state.restaurantId);
  const { data: logs, isLoading } = useSyncLogs(restaurantId || "");

  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Entity</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Records</TableHead>
            <TableHead>Triggered By</TableHead>
            <TableHead>Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-muted-foreground">
                Loading...
              </TableCell>
            </TableRow>
          ) : logs?.length ? (
            logs.map((log: { id: string; entity: string; status: string; records_processed: number; triggered_by: string; created_at: string }) => (
              <TableRow key={log.id}>
                <TableCell className="capitalize">{log.entity}</TableCell>
                <TableCell>
                  <Badge
                    variant={
                      log.status === "success"
                        ? "default"
                        : log.status === "failed"
                        ? "destructive"
                        : "outline"
                    }
                  >
                    {log.status}
                  </Badge>
                </TableCell>
                <TableCell>{log.records_processed}</TableCell>
                <TableCell>{log.triggered_by}</TableCell>
                <TableCell>{formatDate(log.created_at)}</TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-muted-foreground">
                No sync logs found
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
