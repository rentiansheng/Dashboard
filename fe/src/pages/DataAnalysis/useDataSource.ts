import { DataSourceService, ApiDataSourceMetaField } from '@/services';
import { ProFieldValueEnumType, ProSchemaValueEnumType } from '@ant-design/pro-components';
import { chain } from 'lodash';
import { useState, useCallback, useMemo, useRef } from 'react';
import { useDataSourceChart } from './useDataSourceChart';
import { useDataSourceMeta, MetaSchema } from './useDataSourceMeta';
import { useRequest } from 'ahooks';
import { ClientCondition } from './types';
import { ValueTypeMap } from '@/components/Antd/ValueMap';

export default function useDataSource() {
  const [dataSourceId, _setDataSourceId] = useState<number>();
  const { data: dataSources, loading: loadingDataSource } = useRequest(
    DataSourceService.getDataSourceList,
  );

  const list = dataSources?.data?.list;

  const valueEnum= (meta: MetaSchema) : ProFieldValueEnumType => {
    // TODO: metaEnum.api 存在调用api逻辑，后续处理
    if(!meta || !meta.data)  {
      return {};
    };
    // 当前只有配置enum值
    const enumInfo = meta.data?.enum.values
    if(!enumInfo || !enumInfo.type) {
      return  {} ;
    }
    switch(enumInfo.type) {
      case 'kv':
        const res = chain(enumInfo.values)
          .map((key, value) => [key, value])
          .fromPairs()
          .value();
        return res;
        case 'array':
            return chain(enumInfo.values)
            .map((value) => [value, value])
            .fromPairs()
            .value();
    }

    return {};
  };

  const find = useCallback(
    (id: number) => {
      return list?.find((item: any) => item.id === id);
    },
    [list],
  );

  const getMetaSortersByName = useCallback(
    (id: number): string[] => {
      const fields = find(id)?.sort_fields;
      if (!fields) return [];
      return fields.split(',').filter((item: string) => !!item);
    },
    [find],
  );

  const { metas, loading: loadingDataSourceMeta } = useDataSourceMeta(dataSourceId);
 


  const dataSourceValueEnum = useMemo(() => {
    return chain(list)
      .map((item: any) => [item.id, item.name])
      .fromPairs()
      .value();
  }, [list]);

  const timeFieldsValueEnum = useMemo(() => {
    if (!metas?.length) return undefined;
    const result = chain(metas)
      .filter((meta: MetaSchema) => Boolean(meta.data.action?.time))
      .map((meta: MetaSchema) => [meta.data.name, meta.data.display_name])
      .fromPairs()
      .value();
    return Object.keys(result).length > 0 ? result : undefined;
  }, [metas]);

  const sortFieldsValueEnum = useMemo(() => {
    if (!metas?.length) return undefined;
    const result = chain(metas)
      .filter((meta: MetaSchema) => Boolean(meta.data.action?.sort))
      .map((meta: MetaSchema) => [meta.data.name, meta.data.display_name])
      .fromPairs()
      .value();
    return Object.keys(result).length > 0 ? result : undefined;
  }, [metas]);

  const findDataSource = (id: number) => {
    return dataSources?.data?.list?.find((item: any) => item.id === id);
  };

  const [sortValue, setSortValue] = useState<string>();
  const [defaultColumnsState, setDefaultColumnsState] = useState<any>();
  const [backendColumnState, setBackendColumnState] = useState<any>();

 

  const columns = useMemo(() => {
    return metas?.map((meta: MetaSchema) => {
      const fieldType = meta.data.data_type || 'string';
   
      return {
        title: meta.data.display_name || meta.data.name || 'Unknown',
        dataIndex: meta.dataIndex,
        valueType: fieldType,
        valueEnum: valueEnum(meta),
        render: (value: any, record: any) => {
          // 调试信息
          
          // 检查 ValueTypeMap 中是否有对应的渲染器
          if (ValueTypeMap[fieldType] && ValueTypeMap[fieldType].render) {
            try {
              const result = ValueTypeMap[fieldType].render!(value, { text: value, mode: 'read' }, {} as any);
              return result;
            } catch (error) {
              console.warn(`${fieldType} render failed:`, error);
              return value;
            }
          } else {
            return value;
          }
        },
      };
    });
  }, [metas]);

  const getTableParams = ({ formValue, dataSource, groupKeyId }: any) => {
    return async (params: any) => {
      return {
        ...params,
        data_source_id: dataSource?.id,
        group_key_id: groupKeyId,
        sort_field: sortValue,
      };
    };
  };


  type BaseState = {
    date: {
      start: number;
      end: number;
    };
    conditionFields: string[];
  };

  const baseState = useRef<BaseState>();

 
  const syncBaseState = ({ date, conditionFields }: any) => {
    baseState.current = {
      date,
      conditionFields,
    };
  };

  const setDataSourceId = (id: number) => {
    const dataSource = findDataSource(id);
    if (!dataSource) return;
    _setDataSourceId(id);
  };

  const {
    config: chartConfig,
    data: chartData,
    request: requestChart,
    loading: chartLoading,
  } = useDataSourceChart();

  return {
    findDataSource,
    dataSourceValueEnum,
    setDataSourceId,
    loadingDataSource,
    loadingDataSourceMeta,
    timeFieldsValueEnum,
    sortFieldsValueEnum,
    syncBaseState,
    metas,
    columns,
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
    valueEnum,
  };
}
