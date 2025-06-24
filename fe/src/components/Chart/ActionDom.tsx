import { AimOutlined, ReloadOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { Button, Tooltip } from 'antd';
import 'react-splitter-layout/lib/index.css';
import { ChartAdvancedConfig, ChartAdvancedConfigProps } from './ChartAdvancedConfig';
import { IInternalChartProps } from './types';

type IActionDomProps = Pick<IInternalChartProps, 'config'> &
  Pick<ChartAdvancedConfigProps, 'onRefresh' | 'onDebug'>;

export default function ActionDom(props: IActionDomProps) {
  const { config, ...actionHandlers } = props;

  if (!config) return null;

  const { action, key } = config;

  const isClientDefinedChart = !!key;
  const shouldRenderSetting = [action?.refresh, action?.debug].some(Boolean);

  return (
    <>
      {[
        isClientDefinedChart && action?.analysis && (
          <Tooltip title="Analysis" key="detail">
            <Button type="text" style={{ padding: '4px 8px' }}>
              <Link
                target="_blank"
                to={{
                  pathname: '/chart/analysis',
                  search: new URLSearchParams({ chartKey: config.key! }).toString(),
                }}
              >
                <AimOutlined />
              </Link>
            </Button>
          </Tooltip>
        ),

        action?.refresh && (
          <Tooltip title="Refresh" key="refresh">
            <Button type="text" style={{ padding: '4px 8px' }} onClick={actionHandlers.onRefresh}>
              <ReloadOutlined />
            </Button>
          </Tooltip>
        ),

        shouldRenderSetting && (
          <ChartAdvancedConfig key="advanced" action={config.action} {...actionHandlers} />
        ),
      ].filter(Boolean)}
    </>
  );
}
