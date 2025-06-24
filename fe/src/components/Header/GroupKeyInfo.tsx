import React from 'react';
import { Tag, Space, Typography, Tooltip, Breadcrumb } from 'antd';
import { TeamOutlined, BranchesOutlined } from '@ant-design/icons';
import { useSelectedGroupKey } from '@/hooks/useSelectedGroupKey';
import { useGroupKeyTree } from '@/hooks/useGroupKeyTree';
import { GroupKeyTree } from '@/services/groupKeyService';

const { Text } = Typography;

/**
 * 构建从根节点到指定节点的完整路径
 */
const buildNodePath = (treeData: GroupKeyTree[], targetId: number): GroupKeyTree[] => {
  const findPath = (nodes: GroupKeyTree[], path: GroupKeyTree[] = []): GroupKeyTree[] | null => {
    for (const node of nodes) {
      const currentPath = [...path, node];
      
      if (node.id === targetId) {
        return currentPath;
      }
      
      if (node.children && node.children.length > 0) {
        const foundPath = findPath(node.children, currentPath);
        if (foundPath) {
          return foundPath;
        }
      }
    }
    return null;
  };
  
  return findPath(treeData) || [];
};

/**
 * 显示当前选中部门信息的组件
 */
export const GroupKeyInfo: React.FC = () => {
  const {
    hasSelectedGroup,
    getSelectedGroupName,
    getSelectedGroupId,
    hasChildren,
    getChildrenCount,
    getSelectedGroupInfos,
    selectedGroupKeys,
  } = useSelectedGroupKey();
  
  const { rawTreeData } = useGroupKeyTree();

  if (!hasSelectedGroup) {
    return (
      <Tag color="default" icon={<TeamOutlined />}>
        未选择部门
      </Tag>
    );
  }

  // 获取选中的部门信息
  const selectedGroupInfos = getSelectedGroupInfos();
  const isMultipleSelection = selectedGroupInfos.length > 1;

  return (
    <Space direction="vertical" size="small" style={{ width: '100%' }}>
      {/* 多选情况：显示汇总信息 */}
      {isMultipleSelection ? (
        <Space direction="vertical" size="small">
          <Text type="secondary" style={{ fontSize: '12px' }}>
            已选择 {selectedGroupInfos.length} 个部门
          </Text>
          <Space wrap>
            {selectedGroupInfos.map((group) => (
              <Tooltip 
                key={group.id} 
                title={`${group.name} (ID: ${group.id})`}
              >
                <Tag color="blue" icon={<TeamOutlined />}>
                  {group.name}
                </Tag>
              </Tooltip>
            ))}
          </Space>
        </Space>
      ) : (
        /* 单选情况：显示完整路径 */
        <>
          {/* 显示完整路径 */}
          {(() => {
            const selectedGroupInfo = selectedGroupInfos[0];
            const nodePath = buildNodePath(rawTreeData, selectedGroupInfo.id);
            
            return nodePath.length > 1 ? (
              <Breadcrumb
                items={nodePath.map((node) => ({
                  title: (
                    <Tooltip title={`部门ID: ${node.id}`}>
                      <Tag color="blue"  style={{ fontSize: '12px' }}>
                        {node.name}
                      </Tag>
                    </Tooltip>
                  ),
                  key: node.id,
                }))}
                separator=">"
                style={{ fontSize: '12px' }}
              />
            ) : null;
          })()}
          
          
        </>
      )}
    </Space>
  );
}; 