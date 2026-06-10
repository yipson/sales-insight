import EChart from "@/shared/components/charts/EChart";
import type { EChartsOption } from "echarts";
import type { SalesByEmployeeResponse } from "../types";

interface SalesByEmployeeChartProps {
  data?: SalesByEmployeeResponse;
  isLoading?: boolean;
}

export default function SalesByEmployeeChart({ data, isLoading }: SalesByEmployeeChartProps) {
  if (isLoading || !data?.employees?.length) {
    return (
      <div className="flex h-[300px] items-center justify-center rounded-lg border bg-card text-muted-foreground">
        No data available
      </div>
    );
  }

  const option: EChartsOption = {
    tooltip: {
      trigger: "axis",
      axisPointer: { type: "shadow" },
    },
    grid: { left: "3%", right: "4%", bottom: "3%", containLabel: true },
    xAxis: {
      type: "value",
      axisLabel: { formatter: (value: number) => `$${(value / 100).toFixed(0)}` },
    },
    yAxis: {
      type: "category",
      data: data.employees.map((e: { employee_name: string }) => e.employee_name).reverse(),
    },
    series: [
      {
        name: "Total Sales",
        type: "bar",
        data: data.employees.map((e: { total_sales: number }) => e.total_sales / 100).reverse(),
        itemStyle: { color: "hsl(var(--primary))" },
      },
    ],
  };

  return <EChart option={option} style={{ height: 300 }} />;
}
