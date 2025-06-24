import { useQuery } from '@tanstack/react-query';
import { apiQcDepartmentQuery } from '@/services';
import { mapTreeList } from '@/utils/tree';

export interface DeptItem {
  id: number;
  department_name: string;
  department_desc: string;
  order: number;
  children?: DeptItem[];
}

export const useDepartmentQuery = () => {
  return useQuery({
    queryKey: [apiQcDepartmentQuery.requestConfig.path],
    queryFn: apiQcDepartmentQuery,
    select: (data) => {
      return mapTreeList(data, (item) => ({
        department_name: item.department_name,
        department_desc: '',
        id: item.id,
        order: item.order,
      })) as DeptItem[];
    },
  });
}; 