import { ExportOutlined } from '@ant-design/icons';
import {
  ActionType,
  ProColumnType,
  ProTable as AntdProTable,
  ProTableProps,
} from '@ant-design/pro-components';
import { useRequest } from 'ahooks';
import { SortOrder } from 'antd/es/table/interface';
import { merge } from 'lodash';
import { MutableRefObject, useMemo, useRef, useState } from 'react';
import { ToolbarButton } from './ToolbarButton';

import dayjs from 'dayjs';

export type TableParam<T extends Record<string, any>> = {
  params: {
    pageSize: number;
    current: number;
    keyword?: string;
  };
  sort?: Record<keyof T, SortOrder>;
  filter?: Record<keyof T, (string | number)[] | null>;
};

export type ProQueryTableProps<DataSource extends Record<string, any>, Params> = Omit<
  ProTableProps<DataSource, any, any>,
  'request' | 'params'
> & {
  request?: CustomRequestFn<DataSource, Params>;
  params?: (params: TableParam<DataSource>) => Promise<Params | undefined> | Params | undefined;
  indexColumn?: boolean;
  exportType?: ExportType;
};

export const ProQueryTable = <DataSource extends Record<string, any>, Param>(
  props: ProQueryTableProps<DataSource, Param>,
) => {
  const innerActionRef = useRef<ActionType>();
  const { indexColumn, columns, request, params: getParams, exportType, ...restTableProps } = props;
  const actionRef: MutableRefObject<ActionType> = (props.actionRef as any) || innerActionRef;

  const tableColumns = useMemo(() => {
    if (!columns) return;
    if (!indexColumn) return columns;
    const indexColumnDef: ProColumnType<DataSource> = {
      title: 'No.',
      dataIndex: 'No',
      render: (_text, _record, index) => {
        if (!actionRef.current?.pageInfo) return '-';
        const { current, pageSize } = actionRef.current.pageInfo;
        return (current - 1) * pageSize + index + 1;
      },
      width: 80,
    };
    return [indexColumnDef].concat(columns);
  }, [columns, indexColumn, actionRef]);

  const getReqData = async (params: TableParam<DataSource>): Promise<Param | undefined> => {
    if (!params) return;
    if (getParams) return await getParams(params);
    return params as any;
  };

  const [reqData, setReqParam] = useState<Param | undefined>();

  const [exportInfo, setExportInfo] = useState<IExportInfo>();


  const { loading: isLoading, data } = useRequest(
    () =>
      request!(reqData!).then((res) => ({
        list: res?.list || [],
        total: res?.total || 0,
      })),
    {
      ready: !!request && !!reqData,
      refreshDeps: [reqData],
      onSuccess: (res) => {
        const { path, method } = request!.requestConfig;
        const isGet = method.toLowerCase() === 'get';
        if (!res || !res?.list || !res?.total || !reqData) {
          return;
        }

        if (exportType === ExportType.DataAnalysis) {
          setExportInfo({
            request: {
              api: path,
              method,
              query: isGet
                ? new URLSearchParams(reqData as any).toString()
                : JSON.stringify(reqData),
              metric_name: `${tableProps.toolbar?.title || 'Export'} - ${dayjs().format(
                'YYYY-MM-DD_HH_mm_ss',
              )}`,
            },
            total: res.total,
            type: exportType,
          });
        }

        if (exportType === ExportType.MetricChart) {
          setExportInfo({
            request: reqData as any,
            total: res.total,
            type: exportType,
          });
        }
      },
    },
  );

  const lastRequestParamsRef = useRef<TableParam<DataSource>>();
  const lastRequestParams = lastRequestParamsRef.current;
  const defaultTableProps: ProTableProps<DataSource, any, any> = {
    rowKey: 'No',
    search: {
      labelWidth: 'auto',
    },
    cardBordered: true,
    loading: isLoading,
    pagination: {
      total: data?.total,
      pageSize: lastRequestParams?.params.pageSize,
      current: lastRequestParams?.params.current,
    },
    dataSource: data?.list ?? [],
  };

  const tableProps: ProTableProps<DataSource, any, any> = merge(
    {},
    defaultTableProps,
    restTableProps,
    {
      columns: tableColumns,
      request: async (params: any, sort: any, filter: any) => {
        let pageSize = params.pageSize || 20;
        let current = params.current || 1;
        if (
          lastRequestParams?.params.pageSize === pageSize &&
          lastRequestParams?.params.current === current
        ) {
          pageSize = 20;
          current = 1;
        }
        const reqData = {
          params: {
            ...params,
            pageSize,
            current,
          },
          sort,
          filter,
        };
        const newParams = await getReqData(reqData);
        lastRequestParamsRef.current = reqData;
        setReqParam(newParams);
      },
    },
  );

  return <AntdProTable<DataSource, any, any> {...tableProps} actionRef={actionRef} />;
};
