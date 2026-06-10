import EChart from "@/shared/components/charts/EChart";
import type { EChartsOption } from "echarts";
import type { CategoryCoverageResponse } from "../types";

interface CategoryCoverageChartProps {
  data?: CategoryCoverageResponse;
  isLoading?: boolean;
}

export default function CategoryCoverageChart({ data, isLoading }: CategoryCoverageChartProps) {
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
    legend: { data: ["Covered", "Total"] },
    grid: { left: "3%", right: "4%", bottom: "3%", containLabel: true },
    xAxis: {
      type: "category",
      data: data.employees.map((e: { employee_name: string }) => e.employee_name),
    },
    yAxis: { type: "value" },
    series: [
      {
        name: "Covered",
        type: "bar",
        data: data.employees.map((e: { categories_covered: number }) => e.categories_covered),
        itemStyle: { color: "hsl(var(--primary))" },
      },
      {
        name: "Total",
        type: "bar",
        data: data.employees.map((e: { total_categories: number }) => e.total_categories),
        itemStyle: { color: "hsl(var(--muted-foreground))" },
      },
    ],
  };

  return <EChart option={option} style={{ height: 300 }} />;
}
