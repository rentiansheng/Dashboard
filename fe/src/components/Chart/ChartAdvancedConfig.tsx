import { ChartConfig } from '@/types';
import useUrlState from '@ahooksjs/use-url-state';
import { SettingOutlined } from '@ant-design/icons';
import { Button, Dropdown, MenuProps } from 'antd';

export interface ChartAdvancedConfigProps {
  onRefresh?: () => void;
  onDebug?: () => void;
  action?: ChartConfig['action'];
}

export const ChartAdvancedConfig = (props: ChartAdvancedConfigProps) => {
  const { onRefresh, onDebug, action } = props;
  const [state] = useUrlState();

  const hasUrlDebug = Object.keys(state).includes('debug');

  const items = [
    (hasUrlDebug || action?.debug) && {
      label: 'Config',
      key: ChartAdvancedConfig.Action.DEBUG,
    },
  ].filter(Boolean) as MenuProps['items'];

  if (!items?.length) {
    return null;
  }

  return (
    <Dropdown
      menu={{
        items,
        onClick: ({ key }) => {
          switch (key) {
            case ChartAdvancedConfig.Action.REFRESH:
              onRefresh?.();
              break;
            case ChartAdvancedConfig.Action.DEBUG:
              onDebug?.();
              break;
          }
        },
      }}
    >
      <Button
        type="text"
        style={{
          padding: '4px 8px',
        }}
      >
        <SettingOutlined />
      </Button>
    </Dropdown>
  );
};

ChartAdvancedConfig.Action = {
  REFRESH: 'refresh',
  DEBUG: 'show_config',
};
