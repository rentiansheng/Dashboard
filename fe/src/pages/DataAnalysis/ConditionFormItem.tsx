import { DeleteOutlined, FilterFilled, FilterOutlined } from '@ant-design/icons';
import {
  ProFieldValueEnumType,
  ProFormDependency,
  ProFormField,
  ProSchemaValueEnumObj,
} from '@ant-design/pro-components';
import styled from 'styled-components';
import { Button, Popover, Space, Tooltip } from 'antd';
import { pickBy } from 'lodash';
import { useEffect, useRef, useState } from 'react';
import { Condition, OperatorType } from './types';
import { MetaSchema } from './useDataSourceMeta';
import useDataSource from './useDataSource';

interface IProps {
  formListProps: any;
  conditionValueEnum: ProSchemaValueEnumObj;
  getMetaField: (key: string) => MetaSchema | undefined;
  onOpen?: (isOpen: boolean) => void;
  enabledOperators?: OperatorType[];
}

export const isConditionActive = (condition: Condition) => {
  return !!(
    condition.operator &&
    condition.value !== void 0 &&
    condition.value !== null &&
    (Array.isArray(condition.value) ? condition.value.length > 0 : true)
  );
};

const ConditionItemContainer = styled.div<{
  readonly?: boolean;
}>`
  .ant-form-item {
    margin-bottom: ${(props) => (props.readonly ? 0 : '20px')};
  }

  .ant-form-item-control-input {
    min-height: 0;
  }

  & > .ant-btn {
    background-color: rgba(0, 0, 0, 0.06);
  }
`;

export const ConditionFormItem = (props: IProps) => {
  const { enabledOperators } = props;
  const { index, action } = props.formListProps;
  const { conditionValueEnum, getMetaField } = props;

  const [open, setOpen] = useState(true);

  const { onOpen } = props;
  useEffect(() => {
    onOpen?.(open);
  }, [onOpen, open]);

  const data = action.getCurrentRowData();

  const op = {
    ...action,
    remove: () => action.remove(index),
  };

 
  const {
    valueEnum
  } = useDataSource()

  const getListDom = (data: Condition, action: any, readonly?: boolean) => {
    const isActive = isConditionActive(data);
    const operatorValueEnum = enabledOperators
      ? (pickBy(getMetaField(data?.field)?.operators, (_value, key) =>
          enabledOperators.includes(key as OperatorType),
        ) as ProFieldValueEnumType)
      : getMetaField(data?.field)?.operators;
    const dom = [
      <ProFormField
        key="field"
        {...{
          valueType: 'select',
          name: 'field',
          size: 'small',
          readonly: true,
          valueEnum: conditionValueEnum,
        }}
      />,
      <ProFormField
        key="operator"
        {...{
          valueType: 'select',
          size: 'small',
          name: 'operator',
          readonly,
          fieldProps: {
            onChange: () => {
              action.setCurrentRowData({ ...action.getCurrentRowData(), value: void 0 });
            },
          },
          valueEnum: operatorValueEnum,
        }}
      />,
      <ProFormDependency key="value"  name={['operator']}>
        {({ operator }) => {
          let valueEnumMap = getMetaField(data.field) ? valueEnum(getMetaField(data.field)!) : {};
          // 如果valueEnumMap 有值，则将operator 设置为 select
          if (valueEnumMap && Object.keys(valueEnumMap).length > 0) {
            operator = "select"
          }
          return (
            <ProFormField
              key="value"
              {...{
                ...getMetaField(data.field)?.getValueField(operator),
                size: 'small',
                name: 'value',
                readonly,
                valueEnum: getMetaField(data.field) ? valueEnum(getMetaField(data.field)!) : {},
                fieldProps: {
                  showSearch: true,
                  allowClear: true,
                  placeholder: '请输入或选择',
                  mode: "tags",
                  tokenSeparators: [','," "]
                },
              }}
            />
          );
        }}
      </ProFormDependency>,
    ];

    return (
      <ConditionItemContainer readonly={readonly}>
        {readonly ? (
          <Button type="text">
            <Space>
              {isActive ? <FilterFilled style={{ color: '#1890ff' }} /> : <FilterOutlined />}
              {dom}
              <Tooltip title="Delete">
                <Button
                  type="text"
                  size="small"
                  onClick={() => action.remove()}
                  icon={<DeleteOutlined />}
                  onMouseOver={(e) => {
                    e.stopPropagation();
                  }}
                ></Button>
              </Tooltip>
            </Space>
          </Button>
        ) : (
          <div>
            <div>
              <Space style={{ display: 'flex', justifyContent: 'space-between' }}>
                {dom.slice(0, 2)}
              </Space>
            </div>
            <div>{dom.slice(2)}</div>
          </div>
        )}
      </ConditionItemContainer>
    );
  };

  const openRef = useRef(0);
  const setPopoverOpen = (open: boolean) => {
    if (!openRef.current && !open) {
      openRef.current++;
      return;
    }
    setOpen(open);
  };

  return (
    <Popover
      open={open}
      overlayStyle={{ position: 'fixed' }}
      placement="bottomLeft"
      trigger={['click']}
      arrow={{
        pointAtCenter: false,
      }}
      onOpenChange={(open) => {
        setPopoverOpen(open);
      }}
      content={<div style={{ width: 300 }}>{getListDom(data, op, false)}</div>}
    >
      {getListDom(data, op, true)}
    </Popover>
  );
};
