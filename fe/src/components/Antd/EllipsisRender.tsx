import { Typography } from 'antd';
import { BlockProps } from 'antd/es/typography/Base';
import { useMemo } from 'react';

export interface EllipsisRenderProps extends React.PropsWithChildren<Pick<BlockProps, 'ellipsis'>> {
  children: React.ReactNode;
}

export const EllipsisRender = (props: EllipsisRenderProps) => {
  const { children, ellipsis } = props;
  const ellipsisProps = useMemo(() => {
    if (!ellipsis || typeof ellipsis === 'boolean') {
      return ellipsis;
    }
    return {
      expandable: true,
      tooltip: true,
      ...ellipsis,
    };
  }, [ellipsis]);
  return (
    <Typography.Paragraph ellipsis={ellipsisProps} style={{ marginBottom: 0 }}>
      {children}
    </Typography.Paragraph>
  );
};
