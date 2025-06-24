import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { GroupKeyTree } from '@/services/groupKeyService';

interface GlobalState {
  selectedGroupKeys: GroupKeyTree[];
  setSelectedGroupKeys: (groupKeys: GroupKeyTree[]) => void;
  selectedGroupKeyIds: number[];
  setSelectedGroupKeyIds: (ids: number[]) => void;
  // 保持向后兼容
  selectedGroupKey: GroupKeyTree | null;
  setSelectedGroupKey: (groupKey: GroupKeyTree | null) => void;
  selectedGroupKeyId: number | null;
  setSelectedGroupKeyId: (id: number | null) => void;
}

const GlobalContext = createContext<GlobalState | undefined>(undefined);

const STORAGE_KEY = 'selected_group_key_ids';

interface GlobalProviderProps {
  children: ReactNode;
}

export const GlobalProvider: React.FC<GlobalProviderProps> = ({ children }) => {
  const [selectedGroupKeys, setSelectedGroupKeys] = useState<GroupKeyTree[]>([]);
  const [selectedGroupKeyIds, setSelectedGroupKeyIds] = useState<number[]>([]);

  // 从localStorage初始化选中的部门ID
  useEffect(() => {
    const savedIds = localStorage.getItem(STORAGE_KEY);
    if (savedIds) {
      try {
        const ids = JSON.parse(savedIds);
        if (Array.isArray(ids)) {
          setSelectedGroupKeyIds(ids);
        }
      } catch (error) {
        console.error('Failed to parse saved group key ids:', error);
      }
    }
  }, []);

  // 当选中部门变化时，同步到localStorage
  const handleSetSelectedGroupKeys = (groupKeys: GroupKeyTree[]) => {
    setSelectedGroupKeys(groupKeys);
    const ids = groupKeys.map(key => key.id);
    setSelectedGroupKeyIds(ids);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(ids));
  };

  // 当选中部门ID变化时，同步到localStorage
  const handleSetSelectedGroupKeyIds = (ids: number[]) => {
    setSelectedGroupKeyIds(ids);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(ids));
  };

  // 向后兼容的方法 - 获取第一个选中的部门
  const selectedGroupKey = selectedGroupKeys.length > 0 ? selectedGroupKeys[0] : null;
  const selectedGroupKeyId = selectedGroupKeyIds.length > 0 ? selectedGroupKeyIds[0] : null;

  const handleSetSelectedGroupKey = (groupKey: GroupKeyTree | null) => {
    if (groupKey) {
      handleSetSelectedGroupKeys([groupKey]);
    } else {
      handleSetSelectedGroupKeys([]);
    }
  };

  const handleSetSelectedGroupKeyId = (id: number | null) => {
    if (id !== null) {
      handleSetSelectedGroupKeyIds([id]);
    } else {
      handleSetSelectedGroupKeyIds([]);
    }
  };

  const value: GlobalState = {
    selectedGroupKeys,
    setSelectedGroupKeys: handleSetSelectedGroupKeys,
    selectedGroupKeyIds,
    setSelectedGroupKeyIds: handleSetSelectedGroupKeyIds,
    // 向后兼容
    selectedGroupKey,
    setSelectedGroupKey: handleSetSelectedGroupKey,
    selectedGroupKeyId,
    setSelectedGroupKeyId: handleSetSelectedGroupKeyId,
  };

  return (
    <GlobalContext.Provider value={value}>
      {children}
    </GlobalContext.Provider>
  );
};

export const useGlobalContext = (): GlobalState => {
  const context = useContext(GlobalContext);
  if (context === undefined) {
    throw new Error('useGlobalContext must be used within a GlobalProvider');
  }
  return context;
}; 