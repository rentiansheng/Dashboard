import { ProSchemaValueEnumObj } from '@ant-design/pro-components';
import { chain } from 'lodash';
import { useMemo } from 'react';
import { MetaSchema } from './useDataSourceMeta';

type MenuItemType = {
  description?: string;
  meta: MetaSchema;
  label: React.ReactNode;
  key: string;
};

export function useCondition(metas?: MetaSchema[]) {
 

 


  const menus = useMemo(() => {
    return metas?.map((meta) => ({
      key: meta.data.name,
      label: meta.data.display_name,
      meta,
    }));
  }, [metas]);



  const valueEnum = useMemo(() => {
    return chain(metas)
      .filter((meta) => !meta.data?.action?.time && !!meta.data?.action?.filter)
      .map((meta) => [meta.data.name, meta.data.display_name])
      .fromPairs()
      .value();
  }, [metas]);

  return {
    menus,
    valueEnum,
  };
} 