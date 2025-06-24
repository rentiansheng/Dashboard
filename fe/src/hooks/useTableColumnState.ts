import { ColumnsState } from '@ant-design/pro-components';
import { useCallback } from 'react';
import { useLocalStorage } from 'react-use';

export function useTableColumnState(props: {
  tableKey: string;
  initialColumnsState?: Record<string, ColumnsState>;
}) {
  const { tableKey, initialColumnsState } = props;
  const [defaultColumnsState, setDefaultColumnsState] = useLocalStorage<
    Record<string, ColumnsState>
  >(tableKey, void 0, {
    raw: false,
    serializer: (v) => {
      return JSON.stringify(v);
    },
    deserializer: (v) => {
      return JSON.parse(v);
    },
  });

  const resetColumnsState = useCallback(() => {
    setDefaultColumnsState(initialColumnsState);
  }, [initialColumnsState, setDefaultColumnsState]);

  return {
    defaultColumnsState,
    resetColumnsState,
    setDefaultColumnsState,
  };
}
