import Chart from '@/components/Chart/Chart';
import { ProQueryTable, TableParam } from '@/components/Table';
import ViewPortHeight from '@/components/ViewPortHeight';
import { C } from '@/constants/index';
import { RangePresets } from '@/constants/variables';
import { ApiDataSourceQueryTableRequest, DataSourceService , FilterGroup,FilterRule,TableOutput} from '@/services';
import { ChartType, MetricCycle } from '@/types';
import {
  ActionType,
  BetaSchemaForm,
  ProCard,
  ProForm,
  ProFormColumnsType,
  ProFormField,
} from '@ant-design/pro-components';
import { Button, ConfigProvider, FormInstance, message, Space, Alert } from 'antd';
import dayjs, { Dayjs } from 'dayjs';
import { chain, each, isEmpty, last } from 'lodash';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import ChartConfigForm from './ChartConfigForm';
import ChartDefineSeries, { DefineSeriesFormValue } from './ChartDefineSeries';
import { ConditionForm } from './ConditionForm';
import { isConditionActive } from './ConditionFormItem';
import { ClientCondition, Condition, XAxisBy } from './types';
import useDataSource from './useDataSource';
import styled from 'styled-components';
import { FullscreenExitOutlined, FullscreenOutlined, ReloadOutlined } from '@ant-design/icons';
import { convertToTableQueryRequest } from './convert';
import { useFullscreen } from 'ahooks';
import { useSelectedGroupKey } from '@/hooks/useSelectedGroupKey';

interface ChartFormValue {
  cycle: MetricCycle;
  type: ChartType;
  groupBy?: string;
  topN?: number;
  xAxisBy?: XAxisBy;
}

interface SegmentFormValue {
  dataSourceId: string;
  timeField: string;
  range: [Dayjs, Dayjs];
}

interface ConditionFormValue {
  condition: Condition[];
}

export interface FormValue {
  segment: SegmentFormValue;
  condition: ConditionFormValue;
  chart: ChartFormValue;
  series: DefineSeriesFormValue;
}

enum TabType {
  Table = 'Table',
  Chart = 'Chart',
}

const DataAnalysis = () => {
  const [conditionForm] = ProForm.useForm<ConditionFormValue>();
  const [segmentForm] = ProForm.useForm<SegmentFormValue>();
  const [chartForm] = ProForm.useForm<ChartFormValue>();
  const [seriesForm] = ProForm.useForm<DefineSeriesFormValue>();
  const forms = useRef<Record<string, FormInstance>>({
    segment: segmentForm,
    condition: conditionForm,
    chart: chartForm,
    series: seriesForm,
  });

  // 获取选中的部门信息
  const { hasSelectedGroup, getSelectedGroupId, getSelectedGroupName } = useSelectedGroupKey();

  const getFormValues = async (): Promise<FormValue | null> => {
    try {
      const validationResults = await Promise.all(
        Object.values(forms.current).map((form) => {
          return form.validateFields();
        }),
      );
      
      const result = chain(forms.current)
        .mapValues((form) => form.getFieldsValue())
        .merge()
        .value();
      
      return result;
    } catch (error) {
      console.error('getFormValues: Form validation failed', error);
      return null;
    }
  };

 
  const dataSourceId = ProForm.useWatch('dataSourceId', segmentForm);
  const [activeTab, setActiveTab] = useState<TabType>(TabType.Table);

  const {
    findDataSource,
    dataSourceValueEnum,
    setDataSourceId,
    loadingDataSource,
    loadingDataSourceMeta,
    timeFieldsValueEnum,
    sortFieldsValueEnum,
    syncBaseState,
    metas,
    columns: tableColumns,
    getTableParams,
    sortValue,
    defaultColumnsState,
    setDefaultColumnsState,
    backendColumnState,
    setSortValue,

    requestChart,
    chartConfig,
    chartData,
    chartLoading,
  } = useDataSource();

  const chartRange = useMemo(() => {
    if (!chartData) return null;
    const startTime = chartData[0]?.startTime;
    const endTime = last(chartData)?.endTime;
    if (startTime && endTime) {
      return [startTime, endTime];
    }
    return null;
  }, [chartData]);

  const resetForm = useCallback(
    (id?: string) => {
      each(forms.current, (form) => form.resetFields());
      forms.current.segment.setFieldValue('dataSourceId', id || dataSourceId);
      setSortValue(undefined);
    },
    [dataSourceId, setSortValue],
  );

  useEffect(() => {
    resetForm(dataSourceId);
  }, [dataSourceId, resetForm]);

const convertFormValueToClientCondition = (formValue: FormValue): ClientCondition | null => {
    if (!formValue?.segment?.range || formValue.segment.range.length < 2) {
      console.warn('Date range is not properly set in convertFormValueToClientCondition');
      return null;
    }
    
    return {
      type: formValue.chart.type as ChartType,
      groupKeyId: getSelectedGroupId() || 0, // 从全局状态获取选中的部门ID
      dataSourceId: +formValue.segment.dataSourceId,
      name: findDataSource(+formValue.segment.dataSourceId)?.name || '',
      timeField: formValue.segment.timeField,
      cycle: formValue.chart.cycle, // 直接使用 MetricCycle 类型
      groupBy: formValue.chart.groupBy,
      topN: formValue.chart.topN,
      startTime: formValue.segment.range[0].unix(),
      endTime: formValue.segment.range[1].unix(),
      filters: formValue.condition?.condition?.filter(isConditionActive),
      xAxisBy: formValue.chart.xAxisBy,
      series: formValue.series?.series?.map((series) => ({
        ...series,
        aggType: +series.aggType,
        conditions: series.conditions?.filter(isConditionActive),
      })),
    };
  };

  const actionRef = useRef<ActionType>();


  const  buildConds = (formValue: FormValue) => {
    const items =  formValue.condition?.condition?.filter(isConditionActive);

    return items?.map((item) => ({
      field: item.field,
      operator: item.operator!,
      values: item.value,
      pause: false,
 
   
    } as FilterRule)) as FilterRule[];
  }

  const buildTableOutput = (formValue: FormValue) => {
    return {
      time_field: formValue.segment.timeField,
      fields: [],
      sorts: [],
    } as TableOutput;
  }

 

  const tableParams = async (params: any) => {
    const formValue = await getFormValues();
    if (!formValue || !formValue.segment || !formValue.segment.range) {
      console.warn('Form values are not properly set in tableParams');
      return;
    }
    
    // 获取当前数据源
    const dataSource = findDataSource(+formValue.segment.dataSourceId);
    if (!dataSource) {
      console.warn('DataSource not found');
      return;
    }

    // 获取基础查询参数
    const baseRequest = convertToTableQueryRequest({
      formValue,
      dataSource,
      groupKeyId: getSelectedGroupId() || 0, // 从全局状态获取选中的部门ID
    });

    if (!baseRequest) {
      console.warn('Failed to convert form values to table request');
      return;
    }

    // 合并分页参数
    return {
      ...baseRequest,
      page: {
        page_num: params.current || 1,
        page_size: params.pageSize || 20,
      },
    };
  };

  const columns: ProFormColumnsType[] = [
    {
      title: 'Data Source',
      valueType: 'select',
      dataIndex: 'dataSourceId',
      valueEnum: dataSourceValueEnum,
      formItemProps: {
        rules: [
          {
            required: true,
          },
        ],
      },
      fieldProps: {
        loading: loadingDataSource,
        onChange: (value) => {
          setDataSourceId(+value);
          conditionForm.resetFields();
        },
      },
    },
    {
      valueType: 'group',
      className: 'range-group',
      columns: [
        {
          dataIndex: 'timeField',
          title: 'Range',
          valueType: 'select',
          valueEnum: timeFieldsValueEnum,
          width: 'sm',
          formItemProps: {
            rules: [
              {
                required: true,
              },
            ],
          },
          fieldProps: {
            loading: loadingDataSourceMeta,
          },
        },
        {
          dataIndex: 'range',
          valueType: 'dateRange',
          formItemProps: {
            rules: [
              {
                required: true,
              },
            ],
          },
          fieldProps: {
            presets: RangePresets[C.CYCLE.MONTH],
          },
        },
      ],
    },
  ];

  const [lastCondition, setLastCondition] = useState<ClientCondition>();

  const query = async (tab: string = activeTab) => {
    const formValue = await getFormValues();
    if (!formValue || !formValue.segment || !formValue.segment.range) {
      console.warn('Form values are not properly set');
      return;
    }

    syncBaseState({
      date: {
        start: formValue.segment.range[0].unix(),
        end: formValue.segment.range[1].unix(),
      },
      conditionFields:
        chain(formValue.condition?.condition || [])
          .map((condition) => condition.field)
          .uniq()
          .value() || [],
    });
    if (tab === 'Table') {
      actionRef.current?.reload();
    } else {
      const condition = convertFormValueToClientCondition(formValue);
      if (!condition) return;
      setLastCondition(condition);
      requestChart(condition);
    }
  };



  const fullscreenRef = useRef<HTMLElement>();
  const [isFullscreen, fullscreenAction] = useFullscreen(fullscreenRef);

  const queryDom = (
    <Space>
      <Button
        onClick={(e) => {
          e.stopPropagation();
          resetForm();
        }}
      >
        Reset
      </Button>
      <Button
        type="primary"
        // loading={/* dataSourceTable.isLoading ||  */ dataSourceChart.isLoading}
        onClick={(e) => {
          e.stopPropagation();
          query(activeTab);
        }}
      >
        Query
      </Button>
    </Space>
  );

  const [enabledAdvancedConfig, setEnabledAdvancedConfig] = useState(false);

  // 包装 DataSourceService.getDataSourceQueryTable 以适配 ProQueryTable 期望的格式
  const handleTableRequest = async (params: any) => {
    try {
      const response = await DataSourceService.getDataSourceQueryTable(params);
      
      // 确保返回正确的格式：{ list: [], total: number }
      return {
        list: response.data?.list || [],
        total: response.data?.total || 0,
      };
    } catch (error) {
      console.error('Table request failed:', error);
      return {
        list: [],
        total: 0,
      };
    }
  };

  return (
    <ViewPortHeight>
      <StyledSpace size="large" direction="vertical">
        <ProCard bordered>
          <BetaSchemaForm<SegmentFormValue>
            columns={columns}
            form={segmentForm}
            submitter={false}
            name="segment"
            labelWidth="auto"
            style={{ padding: 0 }}
            layoutType="QueryFilter"
          />

          {!!dataSourceId && (
            <ConditionForm
              {...{
                metas,
                formProps: {
                  style: {
                    marginBottom: 16,
                  },
                },
                form: conditionForm,
                onSearch: () => {
                  query();
                },
              }}
            />
          )}

          {queryDom}
        </ProCard>

        <ConfigProvider getPopupContainer={() => fullscreenRef.current ?? document.body}>
          <ProCard
            colSpan={18}
            bordered
            key={dataSourceId}
            style={{ flexGrow: 1 }}
            ref={fullscreenRef as any}
            tabs={{
              type: 'card',
              activeKey: activeTab,
              tabBarExtraContent: (
                <Space style={{ marginRight: 20 }}>
                  <Button
                    onClick={() => fullscreenAction.toggleFullscreen()}
                    icon={isFullscreen ? <FullscreenExitOutlined /> : <FullscreenOutlined />}
                  >
                    {isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}
                  </Button>

                  {/* {activeTab === TabType.Chart && !!chartData && !isEmpty(chartData) && (
                    <Button>Save Chart</Button>
                  )} */}
                  <Button icon={<ReloadOutlined />} onClick={() => query(activeTab)} type="primary">
                    Query
                  </Button>
                </Space>
              ),
              onChange: (key) => {
                setActiveTab(key as TabType);
              },
            }}
          >
            <ProCard.TabPane key={TabType.Table} tab={TabType.Table}>
              <div style={{ margin: '0 -20px' }}>
                {/* 
                  ValueTypeMap 数据渲染说明：
                  1. 时间字段 (date, datetime) -> 自动格式化为 'YYYY-MM-DD HH:mm:ss'
                  2. 布尔字段 (bool, boolean) -> 显示为 Switch 组件
                  3. 邮箱字段 (包含 'email') -> 显示为可点击的邮件链接
                  4. 链接字段 (包含 'link', 'url') -> 显示为可点击的外部链接
                  5. 持续时间字段 (包含 'duration') -> 自动转换为可读格式
                  6. 其他字段 -> 使用 autoTransformData 进行智能转换
                */}
                <ProQueryTable
                  {...{
                    indexColumn: true,
                    actionRef,
                    request: handleTableRequest,
                    manualRequest: true,
                    columns: tableColumns, // 这里使用了 ValueTypeMap 渲染的列配置
                    cardBordered: false,
                    rowKey: 'id',
                    columnsState: {
                      value: defaultColumnsState ?? backendColumnState,
                      defaultValue: backendColumnState,
                      onChange: setDefaultColumnsState,
                    },
                    toolbar: {
                      subTitle: [
                        <ProFormField
                          key="sort"
                          {...{
                            valueType: 'select',
                            valueEnum: sortFieldsValueEnum,
                            formItemProps: {
                              style: {
                                marginBottom: 0,
                              },
                            },
                            fieldProps: {
                              value: sortValue,
                              onChange: (value: string) => {
                                setSortValue(value);
                              },
                            },
                            label: 'Sort by',
                          }}
                        />,
                      ],
                    },
                    params: tableParams,
                    search: false,
                    options: {
                      reload: false,
                    },
                    scroll: {
                      y: 500,
                    },
                  }}
                />
              </div>
            </ProCard.TabPane>
            <ProCard.TabPane key={TabType.Chart} tab={TabType.Chart}>
              <ProCard split="vertical" style={{ margin: -16 }}>
                <ProCard title="Chart" colSpan={5}>
                  <ChartConfigForm
                    {...{
                      isRequireDimension: !enabledAdvancedConfig,
                      formProps: {
                        name: 'chart',
                        form: chartForm,
                        submitter: false,
                        onValuesChange: (changedValues) => {
                          if (changedValues.groupBy && !chartForm.getFieldValue('topN')) {
                            chartForm.setFieldValue('topN', 3);
                          }
                        },
                        initialValues: {
                          cycle: C.CYCLE.WEEK,
                          type: C.CHART_TYPE.COLUMN.value,
                          xAxisBy: XAxisBy.TIME,
                        },
                      },
                      metas,
                    }}
                  />
                </ProCard>
                <ChartProCard colSpan={18} headerBordered>
                  {chartConfig ? (
                    <>
                      <ProCard
                        title="Define Series"
                        style={{ padding: 0, margin: '-16px -16px 16px' }}
                        collapsible
                        defaultCollapsed
                      >
                        <ChartDefineSeries
                          {...{
                            metas,
                            onFormListChange: (series) => {
                              setEnabledAdvancedConfig(!!series.length);
                            },
                            formProps: {
                              form: seriesForm,
                            },
                          }}
                        />
                      </ProCard>
                      <div>
                        {!!(chartRange?.length === 2) && (
                          <div style={{ marginBottom: 34 }}>
                            Current chart time is from{' '}
                            <b>{dayjs.unix(chartRange[0]).format('YYYY-MM-DD')}</b> to{' '}
                            <b>{dayjs.unix(chartRange[1]).format('YYYY-MM-DD')}</b>
                          </div>
                        )}
                        <Chart
                          {...{
                            config: chartConfig,
                            data: chartData,
                            clickable: false,
                            isLoading: chartLoading,
                          }}
                        />
                      </div>
                    </>
                  ) : (
                    <div style={{ 
                      display: 'flex', 
                      justifyContent: 'center', 
                      alignItems: 'center', 
                      height: '300px',
                      flexDirection: 'column',
                      color: '#999'
                    }}>
                      <p>No chart configuration available</p>
                      <p>Please configure chart settings and click Query to load data</p>
                      {chartLoading && <p>Loading...</p>}
                    </div>
                  )}
                </ChartProCard>
              </ProCard>
            </ProCard.TabPane>
          </ProCard>
        </ConfigProvider>
      </StyledSpace>
    </ViewPortHeight>
  );
};

export default DataAnalysis;

const StyledSpace = styled(Space)`
  display: flex;
  width: 100%;
  height: 100%;

  .ant-space-item:last-child {
    flex-grow: 1;
    min-height: 0;

    .ant-pro-card-tabs,
    .ant-tabs,
    .ant-tabs-content,
    .ant-tabs-tabpane {
      height: 100%;

      & > .ant-pro-card {
        height: 100%;
        overflow: auto;
      }
    }

    & > .ant-pro-card,
    .ant-tabs-content > .ant-pro-card,
    .ant-pro-card-col > .ant-pro-card {
      height: 100%;
    }
  }
`;

const ChartProCard = styled(ProCard)`
  height: 100%;
  overflow: auto;
`;
