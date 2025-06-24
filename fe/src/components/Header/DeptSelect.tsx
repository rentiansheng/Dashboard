import { ClusterOutlined, DownOutlined } from '@ant-design/icons';
import { Button as AntButton, Popover, Tree } from 'antd';
import styled from 'styled-components';
import { useEffect } from 'react';
import { useGlobalContext } from '@/context/GlobalContext';
import { useGroupKeyTree } from '@/hooks/useGroupKeyTree';

const Button = styled(AntButton)`
  color: inherit !important;
`;

export const DeptSelect: React.FC = () => {
  const { selectedGroupKey, setSelectedGroupKey, selectedGroupKeyId } = useGlobalContext();
  const { treeData, rawTreeData, loading, error, isInitialized, findNodeById } = useGroupKeyTree();

  // 初始化时设置选中的部门
  useEffect(() => {
    if (isInitialized && rawTreeData.length > 0 && !selectedGroupKey) {
      if (selectedGroupKeyId) {
        const found = findNodeById(rawTreeData, selectedGroupKeyId);
        if (found) {
          setSelectedGroupKey(found);
        }
      } else {
        // 如果没有保存的 ID，选择第一个节点
        const firstNode = rawTreeData[0];
        setSelectedGroupKey(firstNode);
      }
    }
  }, [isInitialized, rawTreeData, selectedGroupKeyId, selectedGroupKey, findNodeById, setSelectedGroupKey]);

  const handleSelect = (deptIds: any[]) => {
    if (deptIds.length > 0 && deptIds[0]) {
      const selected = findNodeById(rawTreeData, deptIds[0]);
      if (selected) {
        setSelectedGroupKey(selected);
      }
    }
  };

  if (error) {
    return (
      <Button type="text" icon={<ClusterOutlined />} disabled>
        加载失败
      </Button>
    );
  }

  return (
    <Popover
      content={
        <Tree
          {...{
            showSearch: true,
            dropdownStyle: { maxHeight: 400, overflow: 'auto', minWidth: 300 },
            treeDefaultExpandAll: true,
            treeNodeFilterProp: 'title',
            autoFocus: false,
            multiple: false,
            defaultValue: selectedGroupKey?.id,
            treeData: treeData,
            loading,
            onSelect: handleSelect,
          }}
        />
      }
    >
      <Button type="text" icon={<ClusterOutlined />}>
        {selectedGroupKey ? `${selectedGroupKey.name}(${selectedGroupKey.id})` : 'Select '}
        <DownOutlined />
      </Button>
    </Popover>
  );
};
