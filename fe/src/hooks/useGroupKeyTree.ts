import { useState, useEffect, useCallback } from 'react';
import { GroupKeyService, GroupKeyTree } from '@/services/groupKeyService';
import type { DataNode } from 'antd/es/tree';

/**
 * 自定义Hook，用于获取和管理部门树数据
 * 避免重复请求API
 */
export const useGroupKeyTree = () => {
  const [treeData, setTreeData] = useState<DataNode[]>([]);
  const [rawTreeData, setRawTreeData] = useState<GroupKeyTree[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isInitialized, setIsInitialized] = useState(false);

  // 根据id 找到对应的节点
  const findNodeById = useCallback((data: GroupKeyTree[], id: number): GroupKeyTree | null => {
    for (const item of data) {
      if (item.id === id) {
        return item;
      }
      if (item.children) {
        const found = findNodeById(item.children, id);
        if (found) {
          return found;
        }
      }
    }
    return null;
  }, []);

  // 转换树形数据为 Ant Design Tree 所需格式
  const mapTreeData = useCallback((data: GroupKeyTree[]): DataNode[] => {
    if (!data || !Array.isArray(data)) {
      return [];
    }
    
    return data.map((item) => {
      const mappedItem: DataNode = {
        key: item.id,
        title: item.name,
        children: item.children && item.children.length > 0 
          ? mapTreeData(item.children) 
          : undefined
      };
      return mappedItem;
    });
  }, []);

  // 获取部门树数据
  const fetchGroupKeyTree = useCallback(async () => {
    if (isInitialized) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const res = await GroupKeyService.getGroupKeyTree(0);
      const data = res.data || [];
      
      setRawTreeData(data);
      setTreeData(mapTreeData(data));
      setIsInitialized(true);
    } catch (err) {
      console.error('Failed to fetch group key tree:', err);
      setError(err instanceof Error ? err.message : '获取部门树失败');
    } finally {
      setLoading(false);
    }
  }, [isInitialized, mapTreeData]);

  // 初始化时获取数据
  useEffect(() => {
    fetchGroupKeyTree();
  }, [fetchGroupKeyTree]);

  return {
    treeData,
    rawTreeData,
    loading,
    error,
    isInitialized,
    findNodeById,
    refetch: fetchGroupKeyTree,
  };
}; 