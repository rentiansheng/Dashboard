export type GroupKeyItem = {
  id: number;
  department_name: string;
  department_desc: string;
  order: number;
  children?: GroupKeyItem[];
};
