import http, { ApiResponse } from '../utils/http';

// GroupKey 接口，对应后端的 GroupKey 结构体
export interface GroupKey {
  id: number;
  name: string;
  parent_id: number;
  order_index: number;
  remark: string;
  mtime: number;
  ctime: number;
}

// GroupKeyTree 接口，对应后端的 GroupKeyTree 结构体
export interface GroupKeyTree {
  id: number;
  name: string;
  order_index: number;
  children: GroupKeyTree[];
}

// GroupKey 服务类
export class GroupKeyService {
  /**
   * 获取分组根节点列表
   */
  static async getGroupKeyRoots(): Promise<ApiResponse<GroupKey[]>> {
    return http.get<GroupKey[]>('/api/group/key/roots');
  }

  /**
   * 获取分组树形结构
   * @param rootId 根节点ID
   * @returns 树形结构数据（按 order_index 排序）
   */
  static async getGroupKeyTree(rootId: number): Promise<ApiResponse<GroupKeyTree[]>> {
    const response = await http.get<GroupKeyTree[]>('/api/group/key/tree', {
      params: {
        root_id: rootId
      }
    });
    
    // 对返回的数据进行排序
    if (response.data && Array.isArray(response.data)) {
      response.data = this.sortGroupKeyTreeByOrder(response.data);
    }
    
    return response;
  }

  /**
   * 对 GroupKeyTree 数组按 order_index 进行递归排序
   * @param trees GroupKeyTree 数组
   * @returns 排序后的 GroupKeyTree 数组
   */
  private static sortGroupKeyTreeByOrder(trees: GroupKeyTree[]): GroupKeyTree[] {
    if (!trees || !Array.isArray(trees)) {
      return trees;
    }

    // 对当前级别按 order_index 排序
    const sorted = [...trees].sort((a, b) => {
      const orderA = a.order_index ?? 0;
      const orderB = b.order_index ?? 0;
      return orderA - orderB;
    });

    // 递归排序子级
    return sorted.map(tree => ({
      ...tree,
      children: tree.children ? this.sortGroupKeyTreeByOrder(tree.children) : []
    }));
  }
}

export default GroupKeyService; 