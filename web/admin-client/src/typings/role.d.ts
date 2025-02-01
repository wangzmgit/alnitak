interface RoleListParam {
  page: number;
  pageSize: number;
}

interface RoleItemType {
  id: number;
  name: string;
  code: string;
  desc: string;
  status: number;
  sort: number;
  createdAt: string;
}

interface RoleFormType {
  id?: number;
  name: string;
  code: string;
  desc: string;
}