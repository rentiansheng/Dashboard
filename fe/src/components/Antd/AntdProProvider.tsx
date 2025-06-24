import { ProConfigProvider } from '@ant-design/pro-components';
import dayjs from 'dayjs';
import { PropsWithChildren } from 'react';
import { ConfigProvider } from 'antd';
import enUS from 'antd/locale/en_US';
import 'dayjs/locale/en';
import { ValueTypeMap } from './ValueMap';

dayjs.locale('en');

export const AntDProProvider = (props: PropsWithChildren<unknown>) => {
  return (
    <ProConfigProvider
      {...{
        valueTypeMap: ValueTypeMap,
      }}
    >
      <ConfigProvider
        theme={{
          token: {
            borderRadius: 2,
          },
          components: {
            Table: {},
          },
        }}
        locale={enUS}
      >
        {props.children}
      </ConfigProvider>
    </ProConfigProvider>
  );
};
