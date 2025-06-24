import { chain, invert, mapValues, startCase } from 'lodash';
import { IOperatorConfig, MetaOperatorConfig, MetaValueType, OperatorType } from './types';

const OperatorConfig: IOperatorConfig[] = [
  {
    type: [MetaValueType.STRING, MetaValueType.LINK],
    operator: [
      {
        list: [OperatorType.Like,OperatorType.Prefix,OperatorType.Suffix,OperatorType.Is],
        renderType: 'text',
      },
      {
        list: [OperatorType.In],
        renderType: 'select',
        defaultValue: [],
      },
    ],
    defaultOperator: OperatorType.Like,
    config: {
      label: {
        [OperatorType.Like]: 'Contains',
      },
    },
  },
  {
    type: [MetaValueType.INT],
    operator: [
      {
        list: [
          OperatorType.Is,
          OperatorType.LessThan,
          OperatorType.GreaterThan,
          OperatorType.LessThanOrEqualTo,
          OperatorType.GreaterThanOrEqualTo,
        ],
        renderType: 'digit',
      },
      {
        list: [OperatorType.Between],
        renderType: 'digitRange',
      },
      {
        list: [OperatorType.In],
        renderType: 'select',
        defaultValue: [],
      },
    ],
    defaultOperator: OperatorType.Is,
  },
  {
    type: [MetaValueType.TIME, MetaValueType.ES_TIME],
    operator: [
      {
        list: [OperatorType.Between],
        renderType: 'dateRange',
      },
      {
        list: [
          OperatorType.LessThan,
          OperatorType.GreaterThan,
          OperatorType.LessThanOrEqualTo,
          OperatorType.GreaterThanOrEqualTo,
        ],
        renderType: 'date',
      },
    ],
    defaultOperator: OperatorType.Between,
  },
  {
    type: [MetaValueType.CHECKBOX, MetaValueType.BOOL],
    operator: [
      {
        list: [OperatorType.Is],
        renderType: 'switch',
        defaultValue: false,
      },
    ],
    defaultOperator: OperatorType.Is,
  },
  {
    type: [MetaValueType.ENUM],
    operator: [
      {
        list: [OperatorType.Is],
        renderType: 'select',
      },
      {
        list: [OperatorType.In],
        renderType: 'select',
        defaultValue: [],
      },
    ],
    defaultOperator: OperatorType.In,
  },
];

const LabelMap = mapValues(invert(OperatorType), startCase);

const _getOperatorConfig = (type: MetaValueType) => {
  return OperatorConfig.find((config) => config.type.includes(type));
};

const _getOperatorTypeLabel = (type: MetaValueType, value: OperatorType) =>
  _getOperatorConfig(type)?.config?.label?.[value] || LabelMap[value];

export const getOperatorConfigByValueType = (type: MetaValueType, operatorValue?: OperatorType) => {
  const config = _getOperatorConfig(type);
  const result: MetaOperatorConfig = {
    operators: {},
    defaultOperator: void 0,
    valueType: void 0,
  };
  if (config) {
    const { operator, defaultOperator } = config;
    result.operators = chain(operator)
      .map('list')
      .flatten()
      .map((operator) => [operator, _getOperatorTypeLabel(type, operator)])
      .fromPairs()
      .value();
    result.defaultOperator = defaultOperator;
    const currentOperator = operatorValue || defaultOperator;
    if (currentOperator) {
      const item = operator.find((item) => item.list.includes(currentOperator));
      result.valueType = item?.renderType;
      result.defaultValue = item?.defaultValue;
    }
  }
  return result;
};

// 新增：获取操作符的字段配置
export const getValueFieldConfig = (type: MetaValueType, operator: OperatorType) => {
  const config = _getOperatorConfig(type);
  if (!config) return {};
  
  const operatorConfig = config.operator.find((item) => item.list.includes(operator));
  if (!operatorConfig) return {};
  
  const baseConfig = {
    valueType: operatorConfig.renderType,
    defaultValue: operatorConfig.defaultValue,
  };
  
  // 为支持多值的操作符添加multiple属性
  if (operator === OperatorType.In) {
    return {
      ...baseConfig,
      fieldProps: {
        mode: 'tags',
        showSearch: true,
        allowClear: true,
        tokenSeparators: [',', ' '], // 用户可以用逗号、空格分隔输入
      
      },

    };
  }
  
  return baseConfig;
};
