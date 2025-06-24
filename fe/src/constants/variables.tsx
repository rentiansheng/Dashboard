import { BarChartOutlined, LineChartOutlined, PieChartOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';

const variables = {
  CYCLE: {
    YEAR: 1,
    QUARTER: 2,
    MONTH: 3,
    WEEK: 4,
  },
  CHART_TYPE: {
    COLUMN: {
      value: 'column',
      icon: <BarChartOutlined />,
      label: (
        <span>
          <BarChartOutlined />
        </span>
      ),
    },
    LINE: {
      value: 'line',
      icon: <LineChartOutlined />,
      label: (
        <span>
          <LineChartOutlined />
        </span>
      ),
    },
    PIE: {
      value: 'pie',
      icon: <PieChartOutlined />,
      label: (
        <span>
          <PieChartOutlined />
        </span>
      ),
    },
  },
  AGG_TYPE: {
    COUNT: 1,
    AVERAGE: 2,
    SUM: 3,
  },
  XAXIS_BY: {
    TIME: 'time',
    DIMENSION: 'dimension',
  },
};

export default {
  ...variables,
};

export const RangePresets = {
  [variables.CYCLE.MONTH]: [
    { label: 'Last 3 Months', value: [dayjs().add(-3, 'month'), dayjs().endOf('day')] },
    { label: 'Last 6 Months', value: [dayjs().add(-6, 'month'), dayjs().endOf('day')] },
    { label: 'Last Year', value: [dayjs().add(-1, 'year'), dayjs().endOf('day')] },
  ],
}; 