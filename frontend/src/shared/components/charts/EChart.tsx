import ReactECharts from "echarts-for-react";
import type { EChartsOption } from "echarts";

interface EChartProps {
  option: EChartsOption;
  style?: React.CSSProperties;
  className?: string;
  onEvents?: Record<string, (params: unknown) => void>;
}

export default function EChart({ option, style, className, onEvents }: EChartProps) {
  return (
    <ReactECharts
      option={option}
      style={{ height: 300, ...style }}
      className={className}
      onEvents={onEvents}
      opts={{ renderer: "canvas" }}
    />
  );
}
