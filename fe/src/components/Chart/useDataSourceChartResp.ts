import { C } from '@/constants';
import { ClientCondition } from '@/pages/DataAnalysis/types';
import { ChartConfig } from '@/types';
import { useEffect, useState } from 'react';
import { BaseChartDataItem, mapChartResp } from './mapChartResp';

export type UseDataSourceChartRespProps = {
  resp: any; // 使用 any 类型来避免导入问题
  condition?: ClientCondition;
};

export type UseDataSourceChartRespOptions = {
  title?: string;
};

export default function useDataSourceChartResp(
  { resp, condition }: UseDataSourceChartRespProps,
  options: UseDataSourceChartRespOptions = {},
) {
  const [config, setConfig] = useState<ChartConfig>();
  const [data, setData] = useState<BaseChartDataItem[]>([]);

  useEffect(() => {
    if (!condition) {
      return;
    }
    
    // 检查 resp 是否存在且有数据
    if (!resp || !resp.data) {
      console.log('No response data available');
      return;
    }
    
   
    const isGroupBy = !!condition.groupBy;
    const cycle = condition.cycle;
    const chartType = condition.type;
    
    // 映射图表类型
    let mappedType: string = chartType;
    if (chartType === 'bar') {
      mappedType = 'column';
    }
    
    const newConfig = {
      title: options?.title,
      type: mappedType,
      xField: 'x',
      yField: 'value',
      seriesField: 'name',
      isStack: !!(
        mappedType === 'column' &&
        (isGroupBy || (condition.series?.length ?? 0) > 1)
      ),
      legend:
        isGroupBy || (condition.series?.length ?? 0) > 1
          ? {
              position: 'right',
            }
          : false,
      action: {
        cycleRange: true,
        setting: true,
      },
    };
    
     setConfig(newConfig);
    
    // 传递 resp.data 而不是 resp，因为 mapChartResp 期望的是实际的数据对象
    const mappedData = mapChartResp([resp.data], {
      getCycle: () => cycle,
    });
 
    
    setData(mappedData);
  }, [resp, condition, options?.title]);

  return {
    config,
    data,
  };
}
