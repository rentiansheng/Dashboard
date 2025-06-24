import { ChartConfig, MetricChartConfig } from '@/types';
import { useCallback, useEffect, useState } from 'react';

export function useChartConfig<T extends Record<string, any>>(
  data: T[],
  config?: MetricChartConfig | ChartConfig,
) {
  const [_chartData, _setChartData] = useState<T[]>([]);
  const [_chartConfig, _setChartConfig] = useState<ChartConfig>(config as any);
  const xField = config?.type !== 'pie' ? config?.xField : void 0;
  const shouldConvertArrayXField = ['column', 'line'].includes(config?.type || '');

  const convertConfig = useCallback(
    (config?: MetricChartConfig | ChartConfig): ChartConfig | null => {
      if (!config) return null;
      if (!shouldConvertArrayXField || !xField) return config as ChartConfig;
      return {
        ...config,
        xField: Array.isArray(xField) ? xField.join('_') : xField,
      } as ChartConfig;
    },
    [xField, shouldConvertArrayXField],
  );

  useEffect(() => {
    if (Array.isArray(xField)) {
      const newXField = xField.join('_');
      const _data = data.map((item) => ({
        ...item,
        [newXField]: xField.map((field) => (item as any)[field]).join('\n'),
      }));
      _setChartData(_data);
      _setChartConfig((config) => ({
        ...config,
        xField: newXField,
      }));
    } else {
      _setChartData(data);
    }
  }, [data, xField]);

  useEffect(() => {
    const newConfig = convertConfig(config);
    if (newConfig) _setChartConfig(newConfig);
  }, [config, convertConfig]);

  return {
    data: _chartData,
    config: _chartConfig,
    setConfig: _setChartConfig,
  };
}
