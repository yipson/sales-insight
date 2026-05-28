import { create } from "zustand";
import { subDays, startOfDay, endOfDay } from "date-fns";
import type { DateRange } from "./types";

const today = new Date();

interface DashboardState {
  dateRange: DateRange;
  setDateRange: (range: DateRange) => void;
}

export const useDashboardStore = create<DashboardState>((set) => ({
  dateRange: {
    from: startOfDay(subDays(today, 30)),
    to: endOfDay(today),
  },
  setDateRange: (range) => set({ dateRange: range }),
}));
