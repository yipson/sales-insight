import EChart from "@/shared/components/charts/EChart";
import type { EChartsOption } from "echarts";
import type { SummaryResponse } from "../types";

interface SalesTrendChartProps {
  data?: SummaryResponse;
  isLoading?: boolean;
}

export default function SalesTrendChart({ data, isLoading }: SalesTrendChartProps) {
  if (isLoading || !data) {
    return (
      <div className="flex h-[300px] items-center justify-center rounded-lg border bg-card text-muted-foreground">
        No data available
      </div>
    );
  }

  // Since the backend doesn't provide daily breakdown in summary,
  // we show a simple gauge/metric visualization
  const option: EChartsOption = {
    tooltip: { trigger: "item" },
    series: [
      {
        type: "gauge",
        startAngle: 180,
        endAngle: 0,
        min: 0,
        max: Math.max(data.sales.total_sales * 1.5, 100000),
        splitNumber: 5,
        itemStyle: { color: "hsl(var(--primary))" },
        progress: { show: true, width: 18 },
        pointer: { show: false },
        axisLine: { lineStyle: { width: 18 } },
        axisTick: { show: false },
        splitLine: { length: 15, lineStyle: { width: 2 } },
        axisLabel: {
          distance: 25,
          formatter: (value: number) => `$${(value / 100).toFixed(0)}`,
        },
        detail: {
          valueAnimation: true,
          formatter: (value: number) => `$${(value / 100).toFixed(2)}`,
          fontSize: 24,
          offsetCenter: [0, "30%"],
        },
        data: [{ value: data.sales.total_sales }],
      },
    ],
  };

  return <EChart option={option} style={{ height: 300 }} />;
}
