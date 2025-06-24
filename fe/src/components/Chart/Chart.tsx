import { ProCard } from '@ant-design/pro-components';
import styled from 'styled-components';
import { Alert } from 'antd';
import ActionDom from './ActionDom';
import InternalChart from './InternalChart';
import { IInternalChartProps } from './types';

type IChartProps = IInternalChartProps & {
  children?: (dom: JSX.Element) => JSX.Element;
  cardBodySlot?: JSX.Element | null;
};

const StyledProCard = styled(ProCard)`
  .ant-btn-text {
    padding: 4px 8px;
  }
`;

export default function Chart(props: IChartProps) {
  const chartProps = props;
  const { config, errorMessage, children, cardBodySlot } = chartProps;

  if (!config) return null;
  const { title = '-', layoutType } = config;
  const isLayoutCard = layoutType === 'card';

  const handleDebug = () => {
   };

  let dom = null;
  if (isLayoutCard) {
    dom = (
      <StyledProCard
        headerBordered
        bordered
        className="inner-chart-card"
        title={title}
        extra={<ActionDom config={config} onDebug={handleDebug} onRefresh={props.onRefresh} />}
      >
        {errorMessage ? (
          <Alert message="Error" description={errorMessage} type="error" showIcon />
        ) : (
          <InternalChart {...chartProps} />
        )}
        {cardBodySlot}
      </StyledProCard>
    );
  } else {
    dom = (
      <>
        <InternalChart {...chartProps} />
        {cardBodySlot}
      </>
    );
  }

  if (children) {
    return children(dom);
  } else {
    return dom;
  }
}
