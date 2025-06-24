import { getOptionsByConstantsFactory } from '../utils/enum';
import { chain } from 'lodash';
import C from './variables';

type IC = typeof C;

export const getOptions = getOptionsByConstantsFactory(C);
export const getValueEnumOptions = (...args: Parameters<typeof getOptions>) => {
  const options = getOptions(...args);
  return chain(options)
    .keyBy('value')
    .mapValues((v) => ({ text: v.label, value: v.value }))
    .value();
};

export type ConstantType<T extends keyof IC> = IC[T][keyof IC[T]];
