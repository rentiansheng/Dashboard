import { C } from '@/constants';
import {
  ProForm,
  ProFormDependency,
  ProFormField,
  ProFormItem,
  ProFormList,
  ProFormProps,
} from '@ant-design/pro-components';
import styled from 'styled-components';
import { chain } from 'lodash';
import { useMemo } from 'react';
import { ConditionForm } from './ConditionForm';
import { Condition, OperatorType } from './types';
import { MetaSchema } from './useDataSourceMeta';

const ListDomContainer = styled.div`
  & > .ant-pro-form-list-container {
    display: flex;
    flex: 1 1 auto;
    gap: 16px;
    padding: 16px 16px 0;
    border: 1px #d9d9d9 dashed;
  }

  display: flex;
  gap: 16px;
  margin-bottom: 16px;

  & + & {
    margin: 16px 0;
  }
`;

export interface DefineSeriesFormValue {
  series: {
    seriesName: string;
    aggType: string;
    aggField: string;
    conditions?: Condition[];
  }[];
}

interface IProps {
  metas?: MetaSchema[];
  formProps?: ProFormProps<DefineSeriesFormValue>;
  onFormListChange?: (series: DefineSeriesFormValue['series']) => void;
}

const ChartDefineSeries = (props: IProps) => {
  const { metas, formProps, onFormListChange } = props;

  const [_form] = ProForm.useForm<DefineSeriesFormValue>();
  const form = formProps?.form || _form;

  const seriesMetas = useMemo(() => {
    return metas?.filter((meta) => !meta.data?.nested) || [];
  }, [metas]);

  const metasValueEnum = useMemo(() => {
    return chain(metas)
      .map((meta) => [meta.data.name, meta.data.display_name])
      .fromPairs()
      .value();
  }, [metas]);

  const emitListLengthChange = () => {
    const { series } = form.getFieldsValue();
    onFormListChange?.(series);
  };

  return (
    <ProForm submitter={false} form={form} {...formProps}>
      <ProFormList
        {...{
          name: 'series',
          creatorButtonProps: {
            creatorButtonText: 'Add Series',
          },
          onAfterAdd: () => {
            emitListLengthChange();
          },
          onAfterRemove: () => {
            emitListLengthChange();
          },
          creatorRecord: () => ({
            aggType: C.AGG_TYPE.COUNT.toString(),
          }),
          itemRender: ({ listDom, action }) => {
            return (
              <ListDomContainer>
                {listDom}
                {action}
              </ListDomContainer>
            );
          },
        }}
      >
        <ProFormField
          {...{
            name: 'seriesName',
            label: 'Series Name',
            valueType: 'text',
            formItemProps: {
              rules: [
                {
                  required: true,
                },
              ],
            },
          }}
        />

        <ProFormField
          {...{
            name: 'aggType',
            label: 'Agg Type',
            valueType: 'radio',
            valueEnum: {
              [C.AGG_TYPE.COUNT]: 'Count',
              [C.AGG_TYPE.SUM]: 'Sum',
              [C.AGG_TYPE.AVERAGE]: 'Average',
            },
            formItemProps: {
              rules: [
                {
                  required: true,
                },
              ],
            },
          }}
        />

        <ProFormDependency name={['aggType']}>
          {({ aggType }) => {
            if (+aggType === C.AGG_TYPE.COUNT) return null;
            return (
              <ProFormField
                {...{
                  name: 'aggField',
                  label: 'Agg Field',
                  valueType: 'select',
                  valueEnum: metasValueEnum,
                  fieldProps: {
                    showSearch: true,
                  },
                  formItemProps: {
                    rules: [
                      {
                        required: true,
                      },
                    ],
                  },
                }}
              />
            );
          }}
        </ProFormDependency>

        <ProFormItem
          {...{
            name: 'conditions',
            label: 'Agg Conditions',
          }}
        >
          <ConditionForm
            {...{
              metas: seriesMetas,
              form: false,
              enabledOperators: [OperatorType.In, OperatorType.Between],
              required: true,
              formListProps: {
                min: 1,
                name: 'conditions',
              },
            }}
          />
        </ProFormItem>
      </ProFormList>
    </ProForm>
  );
};

export default ChartDefineSeries;
