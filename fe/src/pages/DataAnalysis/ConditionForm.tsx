import { PlusOutlined } from '@ant-design/icons';
import {
  FormListActionType,
  ProForm,
  ProFormInstance,
  ProFormList as AntProFormList,
  ProFormListProps,
  ProFormProps,
} from '@ant-design/pro-components';
import styled from 'styled-components';
import { Button, Divider, Dropdown, Input, Space } from 'antd';
import { chain, intersection, throttle } from 'lodash';
import { cloneElement, memo, ReactNode, useCallback, useMemo, useRef, useState } from 'react';
import { ConditionFormItem, isConditionActive } from './ConditionFormItem';
import './styles.scss';
import { Condition, OperatorType } from './types';
import { useCondition } from './useCondition';
import { MetaSchema } from './useDataSourceMeta';
import { useLatest } from 'ahooks';

const StyledDropdown = styled(Dropdown)`
  .ant-dropdown-menu {
    overflow: auto;
    max-height: 400px;
  }
`;

const ProConditionForm = styled(ProForm<{ condition: Condition[] }>)`
  .ant-pro-form-list {
    display: inline-flex;
    width: auto;
  }
`;

const ProFormList = styled(AntProFormList)<{ empty?: boolean }>`
  margin-right: ${(props) => (props.empty ? 0 : '20px')};
  margin-bottom: 0 !important;

  .ant-form-item-control-input {
    min-height: 0;
  }

  &
    > .ant-row
    > .ant-col
    > .ant-form-item-control-input
    > .ant-form-item-control-input-content
    > div {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    width: 100% !important;
  }
`;

interface IProps {
  metas: MetaSchema[];
  form?: ProFormInstance | false;
  formProps?: ProFormProps<any>;
  formListProps?: Partial<ProFormListProps<any>>;
  extraRender?: (dom: React.ReactNode) => React.ReactNode;
  onSearch?: (conditions: Required<Condition>[]) => void;
  createButtonText?: string;
  enabledOperators?: OperatorType[];
}

export const ConditionForm = memo(
  ({
    extraRender,
    metas,
    form,
    onSearch,
    formProps,
    createButtonText = 'Filter',
    formListProps,
    enabledOperators: userEnabledOperators,
  }: IProps) => {
    const actionRef = useRef<
      FormListActionType<{
        name: string;
      }>
    >();
    const [_form] = ProForm.useForm();
    const formInstance = form || _form;

    const [searchText, setSearchText] = useState('');
    const onSearchMenu = throttle((value: string) => {
      setSearchText(value);
    }, 300);

    const { menus: conditionFieldMenus, valueEnum: conditionValueEnum } = useCondition(metas);
    const [usedConditionKeys, setUsedConditionKeys] = useState<string[]>([]);
    const conditionFieldMenusWithLeft = useMemo(() => {
      return chain(conditionFieldMenus)
        .filter((menu) => !!(menu?.key && !usedConditionKeys.includes(menu.key.toString())))
        .filter((menu) => (searchText ? !!menu?.meta.data.display_name.includes(searchText) : true))
        .value();
    }, [conditionFieldMenus, searchText, usedConditionKeys]);
    const getMetaField = (metaKey: string) => {
      const res =  metas.find((item) => item.dataIndex === metaKey);
      return res;
    };

    const lastChangedRef = useRef(false);
    const onSearchLatestRef = useLatest(onSearch);
    const onQuery = useCallback(() => {
      lastChangedRef.current = false;
      let formValue = formInstance.getFieldsValue(); 

      onSearchLatestRef.current?.(formValue.condition?.filter(isConditionActive) || []);
    }, [formInstance, onSearchLatestRef]);

    const onConditionLengthChange = () => {
      setUsedConditionKeys(
        actionRef.current?.getList()?.map((item: any) => (item as Condition).field) || [],
      );
      lastChangedRef.current = true;
    };
    Object.assign(formInstance, {
      onQuery,
    });
    const onOpen = useCallback(
      (isOpen: boolean) => {
        if (isOpen) return;
        onQuery();
      },
      [onQuery],
    );
    const AddConditionButtonDom =
      conditionFieldMenusWithLeft.length || !!searchText ? (
        <StyledDropdown
          overlayClassName="condition-fields-add-overlay"
          dropdownRender={(menu) => (
            <div style={{ background: '#fff' }} className="ant-dropdown-menu">
              <Space style={{ padding: 8 }}>
                <Input
                  autoFocus
                  value={searchText}
                  allowClear
                  onChange={(e) => onSearchMenu(e.target.value)}
                />
              </Space>
              <Divider style={{ margin: 0 }} />
              {cloneElement(menu as React.ReactElement, {
                style: {
                  boxShadow: 'none',
                },
              })}
            </div>
          )}
          menu={{
            items: conditionFieldMenusWithLeft,
            onClick: (menu) => {
              const metaField = getMetaField(menu.key);
              if (!metaField) return;
              
              const {
                defaultOperator,
                defaultValue,
                operators: metaEnabledOperatorsValueEnum,
              } = metaField;
              const metaEnabledOperators = Object.keys(metaEnabledOperatorsValueEnum || {}) as OperatorType[];
              const enabledOperators = userEnabledOperators
                ? intersection(metaEnabledOperators, userEnabledOperators)
                : metaEnabledOperators;
              let operator: OperatorType | undefined = defaultOperator;
              if (operator && !enabledOperators.includes(operator)) {
                operator = enabledOperators[0];
              }
              if (!operator) return;
              actionRef.current?.add({
                field: menu.key,
                operator: operator,
                value: defaultValue,
              });
              onConditionLengthChange();
            },
          }}
        >
          <Button icon={<PlusOutlined />} type="dashed">
            {createButtonText}
          </Button>
        </StyledDropdown>
      ) : null;

    const renderAsForm = (dom: ReactNode) => {
      return (
        <ProConditionForm submitter={false} {...formProps} name="condition" form={formInstance}>
          {dom}
        </ProConditionForm>
      );
    };

    const dom = (
      <ProFormList
        name="condition"
        actionRef={actionRef}
        fieldExtraRender={() => {
          return extraRender ? extraRender(AddConditionButtonDom) : AddConditionButtonDom;
        }}
        onAfterRemove={() => {
          onConditionLengthChange();
          onQuery();
        }}
        actionRender={() => []}
        creatorButtonProps={false}
        {...formListProps}
      >
        {(meta, index, action, count) => {
          return (
            <ConditionFormItem
              {...{
                key: meta.name,
                formListProps: {
                  meta,
                  index,
                  action,
                  count,
                },
                getMetaField,
                enabledOperators: userEnabledOperators,
                onOpen,
                conditionValueEnum,
              }}
            />
          );
        }}
      </ProFormList>
    );

    return form === false ? dom : renderAsForm(dom);
  },
);
