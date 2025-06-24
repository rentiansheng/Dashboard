


// Chart configuration types
export interface ChartConfig {
  type: string;
  title?: string;
  xField?: string;
  yField?: string;
  colorField?: string;
  data?: any[];
  [key: string]: any;
}

// Metric data configuration
export interface MetricDataConfig {
  dataSourceId: number;
  groupKeyId: number;
  dateRange: {
    start: number;
    end: number;
  };
  type: string;
  indexes: any[];
  [key: string]: any;
}

export enum MetricCycle {
  YEAR = 1,
  QUARTER = 2,
  MONTH = 3,
  WEEK = 4,
}

export enum ChartType {
  Line = 'line',
  Bar = 'bar',
  Pie = 'pie',
}


// Pre-process data function type
export type PreProcessDataFn = (data: any[]) => any[];
