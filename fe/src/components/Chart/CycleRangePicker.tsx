import { C } from '@/constants/index';
import { ProFormField } from '@ant-design/pro-components';
import { useControllableValue } from 'ahooks';
import dayjs, { Dayjs } from 'dayjs';
import { Required } from 'utility-types';

const Cycle = C.CYCLE;

const CycleConfig = {
  [Cycle.WEEK]: {
    presets: [
      {
        label: 'Last Month',
        value: [dayjs().add(-1, 'months'), dayjs().endOf('day')],
      },
      {
        label: 'Last 3 Months',
        value: [dayjs().add(-3, 'months'), dayjs().endOf('day')],
      },
    ],
    defaultValue: [dayjs().add(-3, 'months'), dayjs().endOf('day')],
    label: 'Week',
  },
  [Cycle.MONTH]: {
    presets: [
      {
        label: 'Last 3 Months',
        value: [dayjs().add(-3, 'months'), dayjs().endOf('day')],
      },
      {
        label: 'Last 6 Months',
        value: [dayjs().add(-6, 'months'), dayjs().endOf('day')],
      },
      {
        label: 'Last Year',
        value: [dayjs().add(-1, 'years'), dayjs().endOf('day')],
      },
    ],
    defaultValue: [dayjs().add(-6, 'months'), dayjs().endOf('day')],
    label: 'Month',
  },
  [Cycle.QUARTER]: {
    presets: [
      {
        label: 'Last Year',
        value: [dayjs().add(-1, 'years'), dayjs().endOf('day')],
      },
    ],
    defaultValue: [dayjs().add(-1, 'years'), dayjs().endOf('day')],
    label: 'Quarter',
  },
  [Cycle.YEAR]: {
    presets: [
      {
        label: 'Last Year',
        value: [dayjs().add(-1, 'years'), dayjs().endOf('day')],
      },
    ],
    defaultValue: [dayjs().add(-1, 'years'), dayjs().endOf('day')],
    label: 'Year',
  },
};

export const getDefaultCycleRange = (cycle: any): [Dayjs, Dayjs] => {
  return (CycleConfig[cycle]?.defaultValue ?? CycleConfig[DefaultValue.cycle].defaultValue) as [
    Dayjs,
    Dayjs,
  ];
};

type CycleRangePickerValue = [dayjs.Dayjs, dayjs.Dayjs];

interface IProps {
  value?: CycleRangePickerValue;
  cycle?: number;
  readonly?: boolean;
  onChange?: (value: CycleRangePickerValue) => void;
}

const DefaultValue: Required<IProps, 'cycle' | 'value'> = {
  cycle: Cycle.MONTH,
  value: [dayjs().add(-3, 'months'), dayjs().endOf('day')],
};

export const CycleRangePicker = (props: IProps) => {
  const { cycle } = props;
  const [value, onChange] = useControllableValue(props);
  const cycleValue = cycle;
  return (
    <ProFormField
      {...{
        valueType: 'dateRange',
        fieldProps: {
          picker: cycleValue
            ? {
                [Cycle.WEEK]: 'week',
                [Cycle.MONTH]: 'month',
                [Cycle.QUARTER]: 'quarter',
                [Cycle.YEAR]: 'year',
              }[cycleValue]
            : 'date',
          format: cycleValue
            ? {
                [Cycle.WEEK]: 'YYYY-[W]w',
                [Cycle.MONTH]: 'YYYY-MMM',
                [Cycle.QUARTER]: 'YYYY-[Q]Q',
                [Cycle.YEAR]: 'YYYY',
              }[cycleValue]
            : 'YYYY-MM-DD',
          value: value,
          presets: CycleConfig[cycle || DefaultValue.cycle]?.presets,
          onChange,
        },
      }}
    />
  );
};
