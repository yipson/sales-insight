import { useState } from "react";
import { useAuthStore } from "@/core/store/authStore";
import { useSyncTrigger } from "@/features/sync/api";
import { Button } from "@/shared/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu";
import { Loader2, RefreshCw } from "lucide-react";

const entities = [
  { value: "orders", label: "Orders" },
  { value: "items", label: "Items" },
  { value: "employees", label: "Employees" },
  { value: "payments", label: "Payments" },
];

export default function SyncTriggerButton() {
  const restaurantId = useAuthStore((state) => state.restaurantId);
  const [selectedEntity, setSelectedEntity] = useState("orders");
  const trigger = useSyncTrigger();

  const handleSync = () => {
    if (!restaurantId) return;
    trigger.mutate({ restaurant_id: restaurantId, entity: selectedEntity });
  };

  return (
    <div className="flex items-center gap-2">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" className="w-32">
            {entities.find((e) => e.value === selectedEntity)?.label}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          {entities.map((entity) => (
            <DropdownMenuItem
              key={entity.value}
              onClick={() => setSelectedEntity(entity.value)}
            >
              {entity.label}
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
      <Button onClick={handleSync} disabled={trigger.isPending || !restaurantId}>
        {trigger.isPending ? (
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        ) : (
          <RefreshCw className="mr-2 h-4 w-4" />
        )}
        Sync Now
      </Button>
    </div>
  );
}
