import { ChartType, MetricCycle } from '@/types';
import { ProFieldValueType } from '@ant-design/pro-components';
import useDataSource from './useDataSource';

export enum OperatorType {
  Is = 'eq',
  In = 'in',
  GreaterThan = 'gt',
  GreaterThanOrEqualTo = 'gte',
  LessThan = 'lt',
  LessThanOrEqualTo = 'lte',
  Between = 'between',
  Like = 'like',
  Prefix = 'prefix',
  Suffix = 'suffix',
}

export enum MetaValueType {
  ENUM = 'select',
  ES_TIME = 'es time',
  TIME = 'time',
  INT = 'int',
  LINK = 'link',
  STRING = 'string',
  CHECKBOX = 'checkbox',
  BOOL = 'bool',
}

export const MetaTypeIsNumber = (value: any) => {
  return value === MetaValueType.ES_TIME || value === MetaValueType.TIME || value === MetaValueType.INT;
}

export enum XAxisBy {
  TIME = 'time',
  DIMENSION = 'dimension',
}

export interface Condition<Value = any> {
  field: string;
  operator?: OperatorType;
  value?: Value;
}

export interface IOperatorConfig {
  type: MetaValueType[];
  operator: {
    list: OperatorType[];
    renderType: ValueType;
    defaultValue?: any;
  }[];
  defaultOperator?: OperatorType;
  config?: {
    label?: Record<string, string>;
  };
}

export type ValueType = ProFieldValueType;

export type MetaOperatorConfig = {
  operators: Partial<Record<OperatorType, string>>;
  defaultOperator: OperatorType | undefined;
  defaultValue?: any;
  valueType: ValueType | undefined;
};

export type UseDataSourceReturn = ReturnType<typeof useDataSource>;

export enum Permission {
  USER = 1,
  DEPT = 2,
}

export interface ISaveChartItem {
  permission: Permission;
  groupKeyId?: number;
  groupName: string;
  name: string;
  tags?: string[];
  condition: ClientCondition;
}

export interface ClientCondition {
  type: ChartType;
  groupKeyId: number;
  dataSourceId: number;
  name: string;
  timeField: string;
  cycle: MetricCycle;
  groupBy?: string;
  topN?: number;
  startTime: number;
  endTime: number;
  filters?: Condition[];
  xAxisBy?: XAxisBy;
  series?: UserDefineSeries[];
}

export interface UserDefineSeries {
  seriesName: string;
  aggType: number;
  aggField: string;
  conditions?: Condition[];
}
