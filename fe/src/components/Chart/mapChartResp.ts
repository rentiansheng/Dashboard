import { C } from '@/constants/index';
import { MetricCycle } from '@/types';
import { GroupKeyItem } from '@/types/groupKey';
import dayjs from 'dayjs';
import { chain } from 'lodash';

function matchTime(mayBeTime: string) {
  const REGEX = /^(\d{4}-\d{2}-\d{2})\n(\d{4}-\d{2}-\d{2})/;
  const matched = mayBeTime.toString().match(REGEX);
  if (matched) {
    const [, startTime, endTime] = Array.from(matched);
    return [startTime, endTime];
  }

  return null;
}

function formatTime(mayBeTime: string, cycle?: MetricCycle, crossYear = false) {
  const matched = matchTime(mayBeTime);
  if (matched) {
    const [start, end] = matched;
    let x = start;
    switch (cycle) {
      case MetricCycle.WEEK:
        x = dayjs(start).format(crossYear ? 'YYYY-MM-DD' : 'MM-DD');
        break;
      case MetricCycle.MONTH:
        x = dayjs(start).format(crossYear ? 'YYYY-MMM' : 'MMM');
        break;
      case MetricCycle.YEAR:
        x = dayjs(start).format('YYYY');
        break;
      default:
        x = dayjs(start).format('MM-DD');
        break;
    }
    return {
      startTime: dayjs(start, 'YYYY-MM-DD').startOf('day').unix(),
      endTime: dayjs(end, 'YYYY-MM-DD').endOf('day').unix(),
      x,
      _x: mayBeTime,
    };
  } else {
    // 如果不是时间格式，直接使用原始值
    return {
      x: mayBeTime || 'Unknown',
    };
  }
}

export const mapChartResp = <T = any>(
  flattenData: T[],
  options: {
    getCycle: (item: T) => number | undefined;
    getChartName?: (item: T) => string;
    getGroupKey?: (item: T) => GroupKeyItem | undefined | null;
  },
) => {
 
  
  const { getCycle, getChartName, getGroupKey } = options;
  let data: BaseChartDataItem[] = [];

  const type = (flattenData[0] as any)?.series?.items?.[0]?.type;
 
  const mapResp = (resp: any, index: number, x: string, crossYear = false) => {
    
    
    const result = resp.series?.items?.map((serie: any) => {
       
      const timeData = formatTime(x, getCycle(resp), crossYear);
       
      // 确保 value 字段有正确的值
      let value = 0;
      if (serie.data && Array.isArray(serie.data)) {
        const rawValue = serie.data[index];
         value = rawValue !== null && rawValue !== undefined ? Number(rawValue) : 0;
      } else if (typeof serie.data === 'number') {
         value = Number(serie.data);
      } else if (serie.data !== null && serie.data !== undefined) {
         value = Number(serie.data);
      }
      
      // 确保 value 是有效数字
      if (isNaN(value)) {
        console.warn('Value is NaN, setting to 0');
        value = 0;
      }
      
       
      const mappedItem = {
        ...timeData,
        name: serie.name || '',
        seriesName: serie.name || '',
        chart: getChartName?.(resp) || '',
        ...(getGroupKey
          ? {
              deptId: getGroupKey?.(resp)?.id,
              deptName: getGroupKey?.(resp)?.department_name,
            }
          : {}),
        value: value,
        // 确保 x 字段存在
        x: timeData.x || x || 'Unknown',
      };
      
       return mappedItem;
    }) || [];
    
     return result;
  };

  if (type === 'pie') {
    data = chain(flattenData)
      .map((resp) => mapResp(resp, 0, ''))
      .flatten()
      .value();
  } else {
    const crossYear =
      chain((flattenData[0] as any)?.x_axis?.data)
        .map((xValue) => matchTime(xValue)?.map((value) => dayjs(value).year()))
        .filter(Boolean)
        .flatten()
        .uniq()
        .value()?.length > 1;

    data = chain((flattenData[0] as any)?.x_axis?.data)
      .map((x, index) =>
        chain(flattenData)
          .map((resp) => mapResp(resp, index, x, crossYear))
          .value(),
      )
      .flattenDeep()
      .value();
  }
  
 
  
  return data;
};

export type BaseChartDataItem = {
  x: string;
  _x: string;
  name: string;
  chart: string;
  groupKeyId?: number;
  groupKeyName?: string;
  value: number;
  startTime?: number;
  endTime?: number;
};
