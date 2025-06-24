import styled from 'styled-components';
import { Tooltip } from 'antd';
import { MouseEventHandler, PropsWithChildren } from 'react';

const StyledToolbarIconButton = styled.div`
  margin-block: 0;
  margin-inline: 4px;
  color: rgba(0, 0, 0, 0.88);
  font-size: 16px;
  cursor: pointer;

  &:hover {
    color: #1890ff;
  }
`;

interface IProps {
  tooltip?: string;
  onClick?: MouseEventHandler<HTMLElement>;
}

export const ToolbarButton = (props: PropsWithChildren<IProps>) => {
  const { onClick, tooltip } = props;
  return (
    <Tooltip title={tooltip}>
      <StyledToolbarIconButton onClick={onClick}>{props.children}</StyledToolbarIconButton>
    </Tooltip>
  );
};
