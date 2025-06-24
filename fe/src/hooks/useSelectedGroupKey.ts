import { useGlobalContext } from '@/context/GlobalContext';
import { GroupKeyTree } from '@/services/groupKeyService';

/**
 * 自定义Hook，用于获取和操作选中的部门信息
 * @returns 选中的部门信息和操作方法
 */
export const useSelectedGroupKey = () => {
  const { 
    selectedGroupKeys, 
    setSelectedGroupKeys, 
    selectedGroupKeyIds, 
    setSelectedGroupKeyIds,
    // 保持向后兼容
    selectedGroupKey, 
    setSelectedGroupKey, 
    selectedGroupKeyId, 
    setSelectedGroupKeyId 
  } = useGlobalContext();

  /**
   * 检查是否有选中的部门
   */
  const hasSelectedGroup = selectedGroupKeys.length > 0;

  /**
   * 获取选中部门的名称（多个部门用逗号分隔）
   */
  const getSelectedGroupNames = (): string => {
    if (selectedGroupKeys.length === 0) {
      return '未选择部门';
    }
    return selectedGroupKeys.map(key => key.name).join(', ');
  };

  /**
   * 获取第一个选中部门的名称（向后兼容）
   */
  const getSelectedGroupName = (): string => {
    return selectedGroupKey?.name || '未选择部门';
  };

  /**
   * 获取选中部门的ID数组
   */
  const getSelectedGroupIds = (): number[] => {
    return selectedGroupKeyIds;
  };

  /**
   * 获取第一个选中部门的ID（向后兼容）
   */
  const getSelectedGroupId = (): number | null => {
    return selectedGroupKeyId;
  };

  /**
   * 获取选中部门的完整信息数组
   */
  const getSelectedGroupInfos = (): GroupKeyTree[] => {
    return selectedGroupKeys;
  };

  /**
   * 获取第一个选中部门的完整信息（向后兼容）
   */
  const getSelectedGroupInfo = (): GroupKeyTree | null => {
    return selectedGroupKey;
  };

  /**
   * 检查选中的部门是否有子部门
   */
  const hasChildren = (): boolean => {
    return selectedGroupKeys.some(key => key.children?.length > 0);
  };

  /**
   * 获取所有子部门数量
   */
  const getChildrenCount = (): number => {
    return selectedGroupKeys.reduce((total, key) => total + (key.children?.length || 0), 0);
  };

  /**
   * 清空选中的部门
   */
  const clearSelection = (): void => {
    setSelectedGroupKeys([]);
  };

  /**
   * 添加部门到选中列表
   */
  const addGroupKey = (groupKey: GroupKeyTree): void => {
    if (!selectedGroupKeys.find(key => key.id === groupKey.id)) {
      setSelectedGroupKeys([...selectedGroupKeys, groupKey]);
    }
  };

  /**
   * 从选中列表中移除部门
   */
  const removeGroupKey = (id: number): void => {
    setSelectedGroupKeys(selectedGroupKeys.filter(key => key.id !== id));
  };

  /**
   * 检查指定部门是否被选中
   */
  const isGroupKeySelected = (id: number): boolean => {
    return selectedGroupKeys.some(key => key.id === id);
  };

  return {
    // 多选状态
    selectedGroupKeys,
    selectedGroupKeyIds,
    setSelectedGroupKeys,
    setSelectedGroupKeyIds,
    
    // 向后兼容的状态
    selectedGroupKey,
    selectedGroupKeyId,
    setSelectedGroupKey,
    setSelectedGroupKeyId,
    
    // 便捷方法
    hasSelectedGroup,
    getSelectedGroupNames,
    getSelectedGroupName,
    getSelectedGroupIds,
    getSelectedGroupId,
    getSelectedGroupInfos,
    getSelectedGroupInfo,
    hasChildren,
    getChildrenCount,
    clearSelection,
    addGroupKey,
    removeGroupKey,
    isGroupKeySelected,
  };
}; 