import { format } from "date-fns";
import { useDashboardStore } from "../store";
import { Input } from "@/shared/components/ui/input";
import { Calendar } from "lucide-react";

export default function DateRangePicker() {
  const { dateRange, setDateRange } = useDashboardStore();

  const handleFromChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const from = e.target.valueAsDate;
    if (from) {
      setDateRange({ ...dateRange, from });
    }
  };

  const handleToChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const to = e.target.valueAsDate;
    if (to) {
      setDateRange({ ...dateRange, to });
    }
  };

  return (
    <div className="flex items-center gap-3">
      <Calendar size={18} className="text-muted-foreground" />
      <div className="flex items-center gap-2">
        <Input
          type="date"
          value={format(dateRange.from, "yyyy-MM-dd")}
          onChange={handleFromChange}
          className="w-auto"
        />
        <span className="text-muted-foreground">to</span>
        <Input
          type="date"
          value={format(dateRange.to, "yyyy-MM-dd")}
          onChange={handleToChange}
          className="w-auto"
        />
      </div>
    </div>
  );
}
