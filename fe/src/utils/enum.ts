import { capitalize, isPlainObject, pickBy } from 'lodash';

export interface IFactoryOptions {
  capitalize?: boolean;
  overrideLabel?: Record<string, any>;
}

export interface IEnumOption {
  label: string | React.ReactNode;
  value: any;
  key?: string;
}

export interface IEnumMap {
  [key: string]: any;
}

export interface IEnumConfig {
  [key: string]: {
    value: any;
    label?: string | React.ReactNode;
    icon?: React.ReactNode;
    disabled?: boolean;
    [key: string]: any;
  };
}

// 基础枚举转化工厂函数
export const getOptionsByConstantsFactory =
  <IC extends Record<string, any>>(data: IC) =>
  <Key extends keyof IC, K extends keyof IC[Key]>(
    key: Key,
    pickKeys?: IC[Key][K][],
    options?: IFactoryOptions,
  ): IEnumOption[] => {
    const defaultOptions = { capitalize: false };
    const fOpts = Object.assign(defaultOptions, options);
    const obj = pickKeys ? pickBy(data[key], (v) => pickKeys.includes(v)) : data[key];
    return Object.keys(obj).map((propKey: string) => {
      const getPropKeyLabel = () => {
        return fOpts.overrideLabel?.[propKey] ?? fOpts.capitalize
          ? capitalize((propKey || '').replaceAll(/(\S|\d)_(\S|\d)/g, '$1 $2'))
          : propKey;
      };
      if (isPlainObject(obj[propKey])) {
        return {
          label: obj[propKey]!.label || getPropKeyLabel(),
          value: obj[propKey]!.value,
          key: propKey,
        };
      }
      return {
        label: getPropKeyLabel(),
        value: obj[propKey] as IC[Key][K],
        key: propKey,
      };
    });
  };

// 获取枚举值到标签的映射
export const getEnumValueToLabelMap = <T extends Record<string, any>>(
  enumData: T,
  options?: IFactoryOptions,
): Record<any, string | React.ReactNode> => {
  const optionsList = getOptionsByConstantsFactory({ temp: enumData })('temp', undefined, options);
  return optionsList.reduce((acc, option) => {
    acc[option.value] = option.label;
    return acc;
  }, {} as Record<any, string | React.ReactNode>);
};

// 获取枚举标签到值的映射
export const getEnumLabelToValueMap = <T extends Record<string, any>>(
  enumData: T,
  options?: IFactoryOptions,
): Record<string, any> => {
  const optionsList = getOptionsByConstantsFactory({ temp: enumData })('temp', undefined, options);
  return optionsList.reduce((acc, option) => {
    if (typeof option.label === 'string') {
      acc[option.label] = option.value;
    }
    return acc;
  }, {} as Record<string, any>);
};

// 验证枚举值是否有效
export const isValidEnumValue = <T extends Record<string, any>>(
  enumData: T,
  value: any,
): boolean => {
  const values = Object.values(enumData).map(v => isPlainObject(v) ? v.value : v);
  return values.includes(value);
};

// 获取枚举值的标签
export const getEnumLabel = <T extends Record<string, any>>(
  enumData: T,
  value: any,
  options?: IFactoryOptions,
): string | React.ReactNode | undefined => {
  const valueToLabelMap = getEnumValueToLabelMap(enumData, options);
  return valueToLabelMap[value];
};

// 获取枚举的所有值
export const getEnumValues = <T extends Record<string, any>>(enumData: T): any[] => {
  return Object.values(enumData).map(v => isPlainObject(v) ? v.value : v);
};

// 获取枚举的所有键
export const getEnumKeys = <T extends Record<string, any>>(enumData: T): string[] => {
  return Object.keys(enumData);
};

// 创建枚举配置对象
export const createEnumConfig = <T extends Record<string, any>>(
  enumData: T,
  config: Partial<Record<keyof T, { label?: string; icon?: React.ReactNode; disabled?: boolean }>>,
): IEnumConfig => {
  const result: IEnumConfig = {};
  
  Object.keys(enumData).forEach(key => {
    const value = enumData[key];
    const enumValue = isPlainObject(value) ? value.value : value;
    const defaultLabel = isPlainObject(value) ? value.label : key;
    
    result[key] = {
      value: enumValue,
      label: config[key]?.label || defaultLabel,
      icon: config[key]?.icon || (isPlainObject(value) ? value.icon : undefined),
      disabled: config[key]?.disabled || false,
    };
  });
  
  return result;
};

// 批量验证枚举值
export const validateEnumValues = <T extends Record<string, any>>(
  enumData: T,
  values: any[],
): { valid: any[]; invalid: any[] } => {
  const validValues = getEnumValues(enumData);
  const valid: any[] = [];
  const invalid: any[] = [];
  
  values.forEach(value => {
    if (validValues.includes(value)) {
      valid.push(value);
    } else {
      invalid.push(value);
    }
  });
  
  return { valid, invalid };
};

// 获取枚举的统计信息
export const getEnumStats = <T extends Record<string, any>>(enumData: T) => {
  const keys = getEnumKeys(enumData);
  const values = getEnumValues(enumData);
  const options = getOptionsByConstantsFactory({ temp: enumData })('temp');
  
  return {
    keyCount: keys.length,
    valueCount: values.length,
    optionsCount: options.length,
    hasComplexValues: values.some(v => isPlainObject(v)),
    hasIcons: options.some(opt => opt.label && typeof opt.label !== 'string'),
  };
};
