import { useAuthStore } from "@/core/store/authStore";
import { useDashboardStore } from "@/features/dashboard/store";
import {
  useDashboardSummary,
  useSalesByEmployee,
  useTopProducts,
  useCategoryCoverage,
  useTicketIdeal,
} from "@/features/dashboard/api";
import MetricCard from "@/features/dashboard/components/MetricCard";
import DateRangePicker from "@/features/dashboard/components/DateRangePicker";
import SalesTrendChart from "@/features/dashboard/components/SalesTrendChart";
import SalesByEmployeeChart from "@/features/dashboard/components/SalesByEmployeeChart";
import TopProductsTable from "@/features/dashboard/components/TopProductsTable";
import CategoryCoverageChart from "@/features/dashboard/components/CategoryCoverageChart";
import TicketIdealCard from "@/features/dashboard/components/TicketIdealCard";

export default function DashboardPage() {
  const restaurantId = useAuthStore((state) => state.restaurantId);
  const { dateRange } = useDashboardStore();

  const summaryQuery = useDashboardSummary(restaurantId || "", dateRange.from, dateRange.to);
  const salesByEmployeeQuery = useSalesByEmployee(restaurantId || "", dateRange.from, dateRange.to);
  const topProductsQuery = useTopProducts(restaurantId || "", dateRange.from, dateRange.to);
  const categoryCoverageQuery = useCategoryCoverage(restaurantId || "", dateRange.from, dateRange.to);
  const ticketIdealQuery = useTicketIdeal(restaurantId || "", dateRange.from, dateRange.to);

  const sales = summaryQuery.data?.sales;

  return (
    <div className="space-y-6">
      {/* Date Range Picker */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <DateRangePicker />
      </div>

      {/* Metric Cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="Total Sales"
          value={sales ? `$${(sales.total_sales / 100).toFixed(2)}` : "—"}
          isLoading={summaryQuery.isLoading}
        />
        <MetricCard
          title="Orders"
          value={sales?.order_count ?? "—"}
          isLoading={summaryQuery.isLoading}
        />
        <MetricCard
          title="Avg Ticket"
          value={sales ? `$${sales.avg_ticket.toFixed(2)}` : "—"}
          isLoading={summaryQuery.isLoading}
        />
        <MetricCard
          title="Total Tax"
          value={sales ? `$${(sales.total_tax / 100).toFixed(2)}` : "—"}
          isLoading={summaryQuery.isLoading}
        />
      </div>

      {/* Charts Row 1 */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-card p-4">
          <h3 className="mb-4 text-lg font-semibold">Sales Overview</h3>
          <SalesTrendChart data={summaryQuery.data} isLoading={summaryQuery.isLoading} />
        </div>
        <div className="rounded-lg border bg-card p-4">
          <h3 className="mb-4 text-lg font-semibold">Sales by Employee</h3>
          <SalesByEmployeeChart
            data={salesByEmployeeQuery.data}
            isLoading={salesByEmployeeQuery.isLoading}
          />
        </div>
      </div>

      {/* Charts Row 2 */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-card p-4">
          <h3 className="mb-4 text-lg font-semibold">Top Products</h3>
          <TopProductsTable
            data={topProductsQuery.data}
            isLoading={topProductsQuery.isLoading}
          />
        </div>
        <div className="rounded-lg border bg-card p-4">
          <h3 className="mb-4 text-lg font-semibold">Category Coverage</h3>
          <CategoryCoverageChart
            data={categoryCoverageQuery.data}
            isLoading={categoryCoverageQuery.isLoading}
          />
        </div>
      </div>

      {/* Ticket Ideal */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <TicketIdealCard
          data={ticketIdealQuery.data}
          isLoading={ticketIdealQuery.isLoading}
        />
      </div>
    </div>
  );
}
