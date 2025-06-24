import { ChartConfig } from '@/types';
import { Column, Line, Pie, Plot } from '@ant-design/charts';
import { useLatest } from 'ahooks';
import { Spin } from 'antd';
import { chain, each, merge } from 'lodash';
import { useEffect, useMemo, useState } from 'react';
import { ChartToTable } from './ChartToTable';
import { IInternalChartProps } from './types';
import { useChartConfig } from './useChartConfig';

const EventMap = {
  onClick: 'plot:click',
  onElementClick: 'element:click',
  onLineClick: 'line:click',
};

export default function InternalChart(props: IInternalChartProps) {
  const { clickable, config, data: _data, isLoading } = props;

  const { config: chartConfig, data } = useChartConfig(_data, config);
  const type = chartConfig?.type;
  

  let chartDom = null;
  const [plot, setPlot] = useState<Plot>();

  const cursor = clickable ? 'pointer' : 'default';

  const finalConfig = useMemo(() => {
    // 确保数据格式正确
    const processedData = data?.map(item => ({
      ...item,
      value: item.value !== null && item.value !== undefined ? Number(item.value) : 0,
      name: item.name || '',
    })) || [];

    // 如果是饼图，使用专门的扇形配置
    if (type === 'pie') {
      return {
        data: processedData,
        angleField: 'value',
        colorField: 'name',
        radius: 0.8,
        autoFit: true,
        label: {
          type: 'outer',
        },
        tooltip: {
          "title":""
        },
        legend: {
          color: {
            title: false,
            position: 'right',
            rowPadding: 5,
          },
        },
        onReady: (plot: any) => {
          setPlot(plot);
        },
      };
    }

    // 其他图表的基础配置
    return {
      data: processedData,
      xField: chartConfig?.xField || 'x',
      yField: chartConfig?.yField || 'value',
      seriesField: chartConfig?.seriesField || 'name',
      colorField: 'name',
      autoFit: true,
  
      tooltip: {
        shared: false,
        showCrosshairs: true,
        crosshairs: { type: 'xy' },
        showMarkers: true,
        showContent: true,
        enterable: true,
        title: '',
      },
      onReady: (plot: any) => {
        setPlot(plot);
      },
    };
  }, [config, data, chartConfig, type]);

  const propsRef = useLatest(props);
  const onItemClick = (...args: any[]) => {
    if (!clickable || chartConfig.transformMode === 'accumulate') {
      return;
    }
    propsRef.current?.onItemClick?.(...args);
  };

  useEffect(() => {
    if (!plot) return;
    const bindEventMap: Record<string, any[]> = chain(EventMap)
      .pickBy((_, event) => (propsRef.current as any)?.[event])
      .map((plotEvent, propEvent) => {
        const propEventHandler = (propsRef.current as any)?.[propEvent];
        return [plotEvent, [propEventHandler]];
      })
      .fromPairs()
      .value();

    bindEventMap['plot:click'] = [
      (event: any) => {
        if (event.data?.data) onItemClick?.(event.data.data);
      },
    ];

    const bind = (on = true) => {
      each(bindEventMap, (handlers, event) => {
        handlers.forEach((handler) => {
          plot.chart[on ? 'on' : 'off'](event, handler);
        });
      });
    };

    bind();

    return () => {
      bind(false);
    };
  }, [plot]);
 

  if (isLoading) {
    chartDom = (
      <Spin>
        <div style={{ height: chartConfig?.height || 200 }} />
      </Spin>
    );
  } else {
    try {
      // 使用真实的 API 数据
      if (data && data.length > 0) {
        
        switch (type) {
          case 'column':
            chartDom = <Column {...(finalConfig as any)} />;
            break;
          case 'line':
            chartDom = <Line {...(finalConfig as any)} />;
            break;
          case 'area':
            chartDom = <Line {...(finalConfig as any)} />;
            break;
          case 'table':
            chartDom = <ChartToTable {...(finalConfig as any)} onItemClick={onItemClick} />;
            break;
          case 'pie':
            chartDom = <Pie {...(finalConfig as any)} />;
            break;
          default:
            console.warn('Unknown chart type:', type);
            chartDom = <div>Unsupported chart type: {type}</div>;
            break;
        }
      } else {
        chartDom = <div>No data available for chart</div>;
      }
    } catch (error) {
      console.error('Error rendering chart:', error);
      chartDom = (
        <div style={{ padding: '20px', textAlign: 'center', color: 'red' }}>
          <h3>Chart Rendering Error</h3>
          <p>Type: {type}</p>
          <p>Error: {(error as Error).message}</p>
          <p>Data length: {data?.length || 0}</p>
        </div>
      );
    }
  }

  return (
    <div style={{ height: chartConfig.type === 'table' ? 'auto' : chartConfig?.height || 300 }}>
      {chartDom}
    </div>
  );
}
