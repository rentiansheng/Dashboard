import { C } from '@/constants/index';
import { getValueEnumOptions } from '@/constants/options';
import {
  BetaSchemaForm,
  ProForm,
  ProFormColumnsType,
  ProFormProps,
} from '@ant-design/pro-components';
import { chain } from 'lodash';
import { ClientCondition, XAxisBy } from './types';
import { MetaSchema } from './useDataSourceMeta';

export type ChartFormValue = Pick<ClientCondition, 'cycle' | 'groupBy' | 'topN' | 'type'>;

interface IProps {
  formProps?: Omit<ProFormProps<ChartFormValue>, 'columns'>;
  isRequireDimension?: boolean;
  isReadonlyDimension?: boolean;
  metas?: MetaSchema[];
}

export const useColumns = ({
  metas,
  isRequireDimension = true,
  isReadonlyDimension = false,
}: Omit<IProps, 'formProps'>): ProFormColumnsType<ChartFormValue>[] => {
  return [
    {
      title: 'Visualization',
      dataIndex: 'type',
      valueType: 'radioButton',
      tooltip: 'The exploration method used to present your data.',
      formItemProps: {
        rules: [{ required: true }],
      },
      valueEnum: getValueEnumOptions('CHART_TYPE'),
    },
    {
      valueType: 'dependency',
      name: ['type'],
      columns: ({ type }) => {
        const isPie = type === C.CHART_TYPE.PIE.value;
        if (isPie) return [];
        return [
          {
            title: 'X Axis Type',
            dataIndex: 'xAxisBy',
            valueType: 'radio',
            formItemProps: {
              rules: [{ required: true }],
            },
            valueEnum: {
              [XAxisBy.TIME]: 'Time',
              [XAxisBy.DIMENSION]: 'Dimension',
            },
          },
        ];
      },
    },
    {
      valueType: 'dependency',
      name: ['xAxisBy', 'type'],
      columns: ({ xAxisBy, type }) => {
        const isTime = xAxisBy === XAxisBy.TIME;
        const isPie = type === C.CHART_TYPE.PIE.value;
        const subColumns: ProFormColumnsType<ChartFormValue>[] = [];
        if (isTime && !isPie) {
          subColumns.push({
            title: 'Time Cycle',
            dataIndex: 'cycle',
            valueType: 'select',
            convertValue: (value) => value.toString(),
            transform: (value) => ({ cycle: +value }),
            formItemProps: {
              rules: [{ required: true }],
            },
            valueEnum: getValueEnumOptions('CYCLE'),
          });
        }

        subColumns.push({
          title: 'Dimension',
          dataIndex: 'groupBy',
          readonly: isReadonlyDimension,
          fieldProps: {
            showSearch: true,
          },
          formItemProps:
            (!isTime || isPie) && isRequireDimension
              ? {
                  rules: [{ required: true }],
                }
              : void 0,
          valueType: 'select',
          valueEnum: chain(metas)
            .filter((meta) => !!meta.data.action.key)
            .map((meta) => [meta.data.name, meta.data.display_name])
            .fromPairs()
            .value(),
        });

        return subColumns;
      },
    },
    {
      title: 'Max Series on Dimension',
      tooltip: 'Sets the number of data series displayed in the visualization.',
      dataIndex: 'topN',
      valueType: 'digit',
    },
  ];
};

const ChartConfigForm = ({ formProps, metas, isRequireDimension }: IProps) => {
  return (
    <ProForm {...formProps}>
      <BetaSchemaForm<ChartFormValue>
        {...{
          columns: useColumns({ metas, isRequireDimension }),
          layoutType: 'Embed',
        }}
      />
    </ProForm>
  );
};

export default ChartConfigForm;
