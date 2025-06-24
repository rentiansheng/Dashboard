import React from 'react';
import { ConfigContextPropsType } from '@ant-design/pro-components';
import dayjs from 'dayjs';
import { isPlainObject, mapValues } from 'lodash';
import { CycleRangePicker } from '@/components/Chart/CycleRangePicker';
import { formatDuration } from '@/utils/dayjs';
import { Button, Switch } from 'antd';
import 'dayjs/locale/zh-cn';
import { ArrayRenderGuard } from './ArrayRenderGuard';
import { RenderEmail } from './RenderEmail';
import styled from 'styled-components';

const _valueMap: ConfigContextPropsType['valueTypeMap'] = {
  string: {
    render: (text: string) => {
      return <>{text}</>;
    },
  },
  duration: {
    render: (text: string | number) => {
      return formatDuration(text);
    },
  },
  'es time': {
    render: (text: string | React.ReactElement) => {
      let num: number;
      
      if (typeof text === 'string') {
        num = +text;
      } else if (text.props?.text) {
        num = +text.props.text;
      } else {
        return <>-</>;
      }
      
      return <>{dayjs(num).format('YYYY-MM-DD HH:mm:ss')}</>;
    },
  },
  time: {
    render: (text: string) => {
      let num: number;
      
      if (typeof text === 'string') {
        num = +text;
      } else if (text.props?.text) {
        num = +text.props.text;
      } else {
        return <>-</>;
      }
      return <>{dayjs.unix(+text).format('YYYY-MM-DD HH:mm:ss')}</>;
    },
  },
  bool: {
    render: (value: boolean) => {
      return <Switch disabled value={value} />;
    },
  },
  link: {
    render: (obj) => {
      if (!isPlainObject(obj)) return <>-</>;
      if (!obj.link) return <>{obj.name || ''}</>;
      return (
        <TableLink type="link" target="_blank" href={obj.link} rel="noreferrer">
          {obj.name}
        </TableLink>
      );
    },
  },
  cycleRange: {
    renderFormItem: (_, props) => {
      return <CycleRangePicker {...props.fieldProps} />;
    },
  },
  email: {
    render: (value, props) => {
      const isArrayType = props.fieldProps?.is_array;
      if (isArrayType) {
        if (!Array.isArray(value)) {
          return value;
        } else {
          return value.map((item) => <RenderEmail value={item} />);
        }
      } else {
        return <RenderEmail value={value} />;
      }
    },
  },
};

export const ValueTypeMap = mapValues(_valueMap, (value, valueType) => {
  const newValue = { ...value };
  if (newValue.render) {
    newValue.render = ArrayRenderGuard(valueType, _valueMap);
  }
  return newValue;
});

const TableLink = styled(Button)`
  padding: 0;
  text-overflow: ellipsis;
  white-space: normal;
  width: 100%;
  display: inline-flex;
  height: auto;
  text-align: left;
  overflow: hidden;
  align-items: center;

  & > .ant-btn-icon {
    flex: 0;
    min-width: fit-content;
  }

  & > span {
    width: 100%;
    min-width: 0;
    flex: 1 1 auto;
  }
`;
