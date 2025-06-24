import useDataSourceChartResp from '@/components/Chart/useDataSourceChartResp';
import { DataSourceService } from '@/services';
import { useRequest } from 'ahooks';
import { useState } from 'react';
import { convertFormValueToChartQueryRequestV2 } from './convert';
import { ClientCondition } from './types';

export function useDataSourceChart() {
  const [condition, setCondition] = useState<ClientCondition>();
  const [resData, setResData] = useState<any>();
  const query = useRequest(
    (params) => DataSourceService.getDataSourceQueryChart(params),
    {
      manual: true,
      onSuccess: (data) => {
 
        setResData(data);
      },
    },
  );

  const { config, data } = useDataSourceChartResp({
    resp: resData,
    condition,
  });

  const request = (condition: ClientCondition) => {
    const param = convertFormValueToChartQueryRequestV2(condition);
    if (!param) return;
    setCondition(condition);
    query.run(param);
  };

  const reset = () => {
    setResData(undefined);
    setCondition(undefined);
  };

  return {
    config,
    data,
    request,
    errorMessage: query.error?.message,
    refresh: query.refresh,
    reset,
    loading: query.loading,
  };
}
