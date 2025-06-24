import { capitalize, isPlainObject, pickBy } from 'lodash';

export interface IFactoryOptions {
  capitalize?: boolean;
  overrideLabel?: Record<string, any>;
}

export const getOptionsByConstantsFactory =
  <IC extends Record<string, any>>(data: IC) =>
  <Key extends keyof IC, K extends keyof IC[Key]>(
    key: Key,
    pickKeys?: IC[Key][K][],
    options?: IFactoryOptions,
  ) => {
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
        };
      }
      return {
        label: getPropKeyLabel(),
        value: obj[propKey] as IC[Key][K],
      };
    });
  };
