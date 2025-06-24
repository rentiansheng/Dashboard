import { DataSourceService, ApiDataSourceMetaField } from '@/services';
import { ApiDataSourceMetaFieldAction } from '@/services/dataSourceService';
import { useRequest } from 'ahooks';
import { useMemo } from 'react';
import { getOperatorConfigByValueType, getValueFieldConfig } from './operator';
import { MetaValueType, OperatorType } from './types';

import {
  ProColumnType,
  ProFieldValueEnumType,
  ProFormFieldProps,
  ProSchema,
  ProSchemaValueEnumType,
} from '@ant-design/pro-components';

export interface MetaSchema  extends ProSchema{
  dataIndex: string;
  data: {
    name: string;
    display_name: string;
    nested?: boolean;
    action?: ApiDataSourceMetaFieldAction;
    data_type: string;
    enum: MetaSchemaValueEnumType;
    operators:  Partial<Record<OperatorType, string>>,
    defaultValue?: any;
    defaultOperator: OperatorType | undefined;

    getValueField?: (operator: string) => any;
  };
  operators: Partial<Record<OperatorType, string>>;
  defaultOperator: OperatorType | undefined;
  defaultValue?: any;
  getValueField: (operator: string) => any;
}

export interface MetaSchemaValueEnumType {
  api: {
    path: string;
    dynamic: boolean;
  };
  values: {
    type: string;
    values: any;
  };
}
 

export function useDataSourceMeta(dataSourceId?: number) {
  const { data, loading } = useRequest(
    async () => {
      if (!dataSourceId) return {};
      const resp = await DataSourceService.getDataSourceMeta(dataSourceId);
      return resp.data || {};
    },
    {
      refreshDeps: [dataSourceId],
    }
  );

  // 将data.fields 转换为 MetaSchema[]
  const metas = useMemo(() => {
    return data?.fields?.map((item: ApiDataSourceMetaField) => {
      // 通过 getOperatorConfigByValueType 获取 operators 和 defaultValue
      const operatorConfig = getOperatorConfigByValueType(item.data_type as MetaValueType);
      
      return {
        // 将item.name 转换为 dataIndex
        dataIndex: item.name,
        data: {
          ...item,
          operators: operatorConfig.operators,
          defaultValue: operatorConfig.defaultValue,
          defaultOperator: operatorConfig.defaultOperator,
        },
        operators: operatorConfig.operators, 
        defaultOperator: operatorConfig.defaultOperator,
        defaultValue: operatorConfig.defaultValue,
        getValueField: (operator: string) => {
          return getValueFieldConfig(item.data_type as MetaValueType, operator as OperatorType);
        },
      };
    }) || [] as MetaSchema[];
  }, [data?.fields]);

  return {
    metas,
    loading,
  };
} 