import { useRef, useEffect } from "react";
import * as echarts from "echarts";
import type { EChartsOption } from "echarts";

interface UseEChartOptions {
  option: EChartsOption;
  onEvents?: Record<string, (params: unknown) => void>;
}

export function useEChart({ option, onEvents }: UseEChartOptions) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartInstance.current = echarts.init(chartRef.current);

    if (onEvents) {
      Object.entries(onEvents).forEach(([eventName, handler]) => {
        chartInstance.current?.on(eventName, handler);
      });
    }

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      if (onEvents) {
        Object.entries(onEvents).forEach(([eventName, handler]) => {
          chartInstance.current?.off(eventName, handler);
        });
      }
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, [onEvents]);

  useEffect(() => {
    if (chartInstance.current) {
      chartInstance.current.setOption(option, true);
    }
  }, [option]);

  return chartRef;
}
