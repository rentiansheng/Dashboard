import { ChartTableConfig } from '@/types';
import { ProColumnType, ProTable } from '@ant-design/pro-components';
import { chain, groupBy, uniqueId } from 'lodash';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { BaseChartDataItem } from './mapChartResp';
import { Typography, Tag, Space } from 'antd';
import styled from 'styled-components';

type ChartToTableProps = ChartTableConfig & {
  data: BaseChartDataItem[];
  onItemClick?: (item: BaseChartDataItem) => void;
};

interface TableDataItem extends BaseChartDataItem {
  __serie_key__: string;
  __rowkey__: string;
  _loaded?: boolean;
  children?: TableDataItem[];
  rootType?: string;
  nodeLevel?: number;
  isRoot?: boolean;
  nodeType?: 'root' | 'branch' | 'leaf';
}

export const SeriesKey = '__serie_key__';

const uniqueKey = () => uniqueId('rowkey');
const RowKey = '__rowkey__';

const StyledProTable = styled(ProTable)`
  .ant-pro-card {
    border-bottom: 0;
    margin-bottom: 20px;
  }
  .ant-pro-card-body {
    padding: 0;
  }
  
  .root-node {
    background-color: #f0f9ff;
    font-weight: 600;
  }
  
  .branch-node {
    background-color: #fef3c7;
  }
  
  .leaf-node {
    background-color: #f3f4f6;
  }
  
  .node-type-tag {
    margin-right: 8px;
  }
  
  .tree-indent {
    display: inline-block;
    width: 20px;
    height: 1px;
  }
`;

const RootTypeRenderer = ({ record }: { record: any }) => {
  const { nodeType, nodeLevel = 0, rootType, isRoot } = record;
  
  const getTypeColor = (type: string) => {
    switch (type) {
      case 'root': return 'blue';
      case 'branch': return 'orange';
      case 'leaf': return 'green';
      default: return 'default';
    }
  };

  const getTypeText = (type: string) => {
    switch (type) {
      case 'root': return '根';
      case 'branch': return '分支';
      case 'leaf': return '叶子';
      default: return '未知';
    }
  };

  return (
    <Space>
      {Array.from({ length: nodeLevel }).map((_, index) => (
        <span key={index} className="tree-indent" />
      ))}
      
      {nodeType && (
        <Tag 
          color={getTypeColor(nodeType)} 
          className="node-type-tag"
        >
          {getTypeText(nodeType)}
        </Tag>
      )}
      
      {rootType && (
        <Tag color="purple" className="node-type-tag">
          {rootType}
        </Tag>
      )}
      
      {isRoot && (
        <Tag color="red" className="node-type-tag">
          根节点
        </Tag>
      )}
      
      <span>{record[SeriesKey]}</span>
    </Space>
  );
};

export const ChartToTable = (props: ChartToTableProps) => {
  const { xField, yField, data, seriesField, meta, table } = props;
  const [dataSource, setDataSource] = useState<TableDataItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [expandedRowKeys, setExpandedRowKeys] = useState<string[]>([]);

  const analyzeRootTypes = useCallback((data: BaseChartDataItem[]) => {
    const rootTypes = new Set<string>();
    
    data.forEach(item => {
      if (item.groupKeyId) rootTypes.add('groupKey');
      if (item.name) rootTypes.add('category');
      if ((item as any).deptId) rootTypes.add('department');
    });
    
    return Array.from(rootTypes);
  }, []);

  const mapToTableDataItem = useCallback(
    (data: BaseChartDataItem[], level = 0, parentRootType?: string) => {
      const getSeriesData = (series: string): Record<string, any> => {
        const items = data.filter((item) => item[seriesField] === series);
        return chain(items)
          .map((item) => {
            const yValue = meta?.[yField]?.formatter?.(item[yField]) || item[yField];
            return [item[xField], { ...item, yValue }];
          })
          .fromPairs()
          .merge({ [RowKey]: uniqueKey() }, ...items)
          .value();
      };

      const rootTypes = analyzeRootTypes(data);

      return chain(data)
        .map((item) => item[seriesField])
        .uniq()
        .map((series) => {
          const seriesData = getSeriesData(series);
          const hasChildren = table?.getSubDeptIds && table.getSubDeptIds((seriesData as any).deptId)?.length;
          
          const tableItem = {
            [SeriesKey]: meta?.[seriesField]?.formatter?.(series) || series,
            ...seriesData,
            nodeLevel: level,
            nodeType: level === 0 ? 'root' : (hasChildren ? 'branch' : 'leaf'),
            isRoot: level === 0,
            rootType: parentRootType || (rootTypes.length > 0 ? rootTypes[0] : undefined),
          } as TableDataItem;
          
          if (hasChildren) {
            tableItem.children = [];
          }
          return tableItem;
        })
        .value() as any[] as TableDataItem[];
    },
    [meta, seriesField, table, xField, yField, analyzeRootTypes],
  );

  const tableDataProps = useMemo(() => {
    const xFieldValues = chain(data)
      .map((item) => item[xField])
      .filter((v) => v !== undefined)
      .uniq()
      .value();
    
    const columns: any[] = xFieldValues.map(
      (key) => ({
        key,
        dataIndex: key,
        title: meta?.[xField]?.formatter?.(key) ?? key,
        renderText: (record: any) => {
          if (record?.yValue === null || record?.yValue === void 0) {
            return 'N/A';
          } else {
            return (
              <Typography.Link onClick={() => props.onItemClick?.(record)}>
                {record.yValue}
              </Typography.Link>
            );
          }
        },
        ...table?.columnConfig?.(key),
      }),
    );

    columns.unshift({
      key: SeriesKey,
      dataIndex: SeriesKey,
      title: meta?.[props.seriesField]?.alias ?? props.seriesField,
      render: (_: any, record: any) => <RootTypeRenderer record={record} />,
      width: 300,
    });

    const hasExpandable = !!table?.props?.expandable?.rowExpandable;

    if (hasExpandable) {
      Object.assign(table?.props?.expandable ?? {}, {
        expandedRowKeys,
        onExpand: async (expanded: boolean, record: TableDataItem) => {
          setExpandedRowKeys((list) => {
            if (expanded) {
              return [...list, record[RowKey]];
            } else {
              return list.filter((key) => key !== record[RowKey]);
            }
          });
          const xField0 = xFieldValues[0];
          const onRequestNextDeptLevel = props.table?.onRequestNextDeptLevel;
          const deptId = (record as any)[xField0]?.deptId;
          const setLeaf = () => {
            record.children = void 0;
            setDataSource((list) => list.slice());
          };
          if (!deptId) {
            setLeaf();
            return;
          }
          if (
            expanded &&
            seriesField === 'deptName' &&
            xField0 &&
            !record._loaded &&
            deptId &&
            onRequestNextDeptLevel
          ) {
            setLoading(true);
            try {
              const { data } = await onRequestNextDeptLevel(deptId);
              if (!data?.length) {
                record.children = void 0;
              } else {
                const dataByDeptId = groupBy(data, 'deptId');
                record.children = Object.values(dataByDeptId)
                  .map((list) => {
                    return mapToTableDataItem(list, (record.nodeLevel || 0) + 1, record.rootType);
                  })
                  .flat();
              }
              record._loaded = true;
              setDataSource((list) => list.slice());
            } finally {
              setLoading(false);
            }
          }
        },
      });
    }

    return {
      ...props.table?.props,
      columns,
      dataSource,
      pagination: false,
      onRow: (record: TableDataItem) => ({
        className: record.nodeType ? `${record.nodeType}-node` : '',
      }),
    };
  }, [
    data,
    meta,
    props,
    table,
    dataSource,
    xField,
    expandedRowKeys,
    seriesField,
    mapToTableDataItem,
  ]);

  useEffect(() => {
    const dataSource = mapToTableDataItem(data);
    setDataSource(dataSource);
    if (dataSource.length === 1) {
      tableDataProps?.expandable?.onExpand?.(true, dataSource[0]);
    }
  }, [data, mapToTableDataItem]);

  return (
    <StyledProTable
      {...tableDataProps}
      loading={props.loading || loading}
      dataSource={dataSource}
      rowKey={RowKey}
      cardBordered
      search={false}
      options={{ reload: false, setting: false, density: false }}
    />
  );
};
