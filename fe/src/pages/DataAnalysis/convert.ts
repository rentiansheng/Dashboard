import { C } from '@/constants';
import {
  ApiDataSourceQueryMetricRequest,
  ApiDataSourceQueryTableRequest,
} from '@/services';
import { ChartType, DataSource, MetricCycle } from '@/types';
import { isConditionActive } from './ConditionFormItem';
import { FormValue } from './index';
import { ClientCondition, Condition, ISaveChartItem, OperatorType, XAxisBy } from './types';
import dayjs from 'dayjs';

export enum AggType {
  TIME = 0,
  TOP_N = 1,
  FILTER = 2,
  TOP_N_TIME = 3,
  FILTER_TIME = 4,
}

export function getAggType(isPie: boolean, isGroupBy: boolean) {
  const aggType = isGroupBy ? (isPie ? AggType.TOP_N : AggType.TOP_N_TIME) : AggType.TIME;
  return aggType;
}

export interface IConvertToTableQueryRequestParams {
  formValue: FormValue;
  dataSource: DataSource;
  groupKeyId: number;
}

const convertCondition = (condition: Condition) => {
  const convertValue = (value: unknown) => {
    if (dayjs.isDayjs(value)) {
      return value.unix();
    } else {
      return value;
    }
  };
  const convertValues = (value: unknown) => {
    return Array.isArray(value) ? value?.map(convertValue) : [convertValue(value)];
  };
  return {
    field: condition.field,
    pause: false,
    operator: condition.operator!,
    values: convertValues(condition.value),
  };
};

export const convertToTableQueryRequest = (
  params: IConvertToTableQueryRequestParams,
): Omit<ApiDataSourceQueryTableRequest, 'page'> | undefined => {
  const { formValue, groupKeyId, dataSource } = params;
  if (
    !formValue.segment?.dataSourceId ||
    !formValue.segment?.timeField ||
    !formValue.segment?.range
  ) {
    return undefined;
  }
  const dataSourceId = +formValue.segment.dataSourceId;
  const dataSourceName = dataSource?.name || '';
  return {
    data_source_id: dataSourceId,
    date: {
      start: formValue.segment.range[0].startOf('day').unix(),
      end: formValue.segment.range[1].endOf('day').unix(),
    },
    group_key_id: groupKeyId,
    indexes: [
      {
        name: dataSourceName,
        output: {
          time_field: formValue.segment.timeField,
          sorts: [formValue.segment.timeField],
        },
        filters: [
          {
            rules:
              formValue.condition?.condition?.filter(isConditionActive).map(convertCondition) || [],
            name: dataSourceName,
          },
        ],
      },
    ],
  };
};

export interface ChartFormValue {
  cycle: number;
  type: string;
  groupBy?: string;
  top?: number;
}

export interface IConvertFormValueToChartQueryRequestParams {
  formValue: FormValue & {
    chart: ChartFormValue;
  };
  dataSource: DataSource;
  deptId: number;
}

export const convertFormValueToChartQueryRequestV2 = (
  condition: ClientCondition,
): ApiDataSourceQueryMetricRequest | undefined => {
  const {
    dataSourceId,
    groupBy,
    topN,
    cycle,
    timeField,
    name,
    type,
    startTime,
    endTime,
    groupKeyId,
    filters,
    xAxisBy,
    series,
  } = condition;
  const isPie = type === 'pie';
  const isXAxisTime = xAxisBy === C.XAXIS_BY.TIME;
  return {
    data_source_id: dataSourceId,
    type,
    date: {
      start: startTime,
      end: endTime,
    },
    group_key_id: groupKeyId,
    indexes: [
      {
        name,
        output: {
          time_field: timeField,
          cycle: (!isPie && isXAxisTime ? +cycle : void 0) || C.CYCLE.WEEK,
          aggregator: {
            is_top_n: !!groupBy,
            is_date_histogram: !isPie && (xAxisBy === XAxisBy.TIME || !xAxisBy),
            v2: true,
            ...(groupBy
              ? {
                  top_n: {
                    field: groupBy,
                    num: topN,
                  },
                }
              : {}),
            filters: series?.map((item) => {
              const termsCond = item.conditions
                ?.filter((cond) => cond.operator === OperatorType.In)
                .map((cond) => ({
                  field: cond.field,
                  value: cond.value,
                }));
              const rangeCond = item.conditions
                ?.filter((cond) => cond.operator === OperatorType.Between)
                .map((cond) => ({
                  field: cond.field,
                  value: [cond.value[0], cond.value[1]],
                }));
              return {
                name: item.seriesName,
                agg_type: +item.aggType,
                extra_agg_field: item.aggField,
                terms: termsCond,
                range: rangeCond,
              };
            }),
          },
        },
        filters: [
          {
            rules:
            filters,
            name: name,
          },
        ],
      },
    ],
  };
};

export const convertSaveChartItemToServer = (params: ISaveChartItem) => {
  const { condition, groupKeyId, groupName, name, tags, permission } = params;
  const queryCondition = convertFormValueToChartQueryRequestV2(condition);
  if (!queryCondition || !name || !permission) return;
  return {
    query_condition: queryCondition,
    group_name: groupName,
    name: name,
    rule_type: permission,
    group_key_id: groupKeyId,
    labels: tags,
  };
};

export const convertQueryConditionToClientCondition = (
  condition: ApiDataSourceQueryMetricRequest,
): ClientCondition => {
  return {
    type: condition.type as ChartType,
    groupKeyId: condition.group_key_id,
    dataSourceId: condition.data_source_id,
    name: condition.indexes[0]?.name,
    timeField: condition.indexes[0]?.output?.time_field,
    groupBy: condition.indexes[0]?.output?.aggregator?.top_n?.field,
    topN: condition.indexes[0]?.output?.aggregator?.top_n?.num,
    cycle: condition.indexes[0]?.output?.cycle || MetricCycle.WEEK,
    startTime: condition.date?.start,
    endTime: condition.date?.end,
    xAxisBy: condition.indexes[0]?.output?.aggregator?.is_date_histogram
      ? XAxisBy.TIME
      : XAxisBy.DIMENSION,

    filters: condition.indexes[0]?.filters?.[0]?.rules?.map((rule) => ({
      field: rule.field,
      operator: rule.operator as OperatorType,
      value: rule.values,
    })),
    series: condition.indexes[0]?.output?.aggregator?.filters
      ?.filter((agg) => !!agg.name)
      .map((agg) => {
        return {
          aggType: agg.agg_type || C.AGG_TYPE.COUNT,
          aggField: agg.extra_agg_field!,
          seriesName: agg.name!,
          conditions: ([] as Condition[])
            .concat(
              agg?.terms
                ?.filter((term) => term.field && term.value)
                .map((term) => ({
                  field: term.field!,
                  operator: OperatorType.In,
                  value: term.value!,
                })) || [],
            )
            .concat(
              agg.range?.map((range) => ({
                field: range.field!,
                operator: OperatorType.Between,
                value: [range.from, range.to]!,
              })) || [],
            ),
        };
      }),
  };
};
