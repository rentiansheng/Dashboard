import { ChartConfig, MetricDataConfig, PreProcessDataFn } from '@/types';
import { PlotEvent } from '@ant-design/charts';
import { BaseChartDataItem } from './mapChartResp';

export interface IInternalChartProps {
  config?: ChartConfig;
  data: BaseChartDataItem[];
  clickable?: boolean;
  isLoading?: boolean;
  errorMessage?: string;
  onClick?: (event: PlotEvent) => void;
  onElementClick?: (event: PlotEvent) => void;
  onPointClick?: (item: any, event: PlotEvent) => void;
  onItemClick?: (item: BaseChartDataItem) => void;
  onRefresh?: () => void;
  onDebug?: () => void;
}

type OverrideConfigType = ChartConfig | ((originalConfig: ChartConfig) => ChartConfig);

export type MetricChartProps = {
  overrideConfig?: OverrideConfigType;
  preprocessData?: PreProcessDataFn;
  renderer?: (props: {
    formDom: JSX.Element;
    chartDom: JSX.Element;
    detailDom: JSX.Element | null;
    detailVisible: boolean;
    config: ChartConfig | undefined;
  }) => JSX.Element;
  params?: Partial<MetricDataConfig>;
  onChange?: (params: MetricDataConfig) => void;
};
