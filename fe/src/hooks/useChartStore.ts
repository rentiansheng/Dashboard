import ALL_CHART, { GlobalChartKey } from '@/charts';
import { useCallback, useRef } from 'react';

export default function () {
  const chartBoard = useRef(ALL_CHART);

  const getChartConfig = useCallback(
    (key: GlobalChartKey) => {
      return chartBoard.current[key];
    },
    [chartBoard],
  );

  return {
    getChartConfig,
  };
}
