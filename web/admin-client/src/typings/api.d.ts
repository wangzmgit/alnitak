interface ApiListParam {
  page: number;
  pageSize: number;
}

interface ApiItemType {
  id: number;
  method: string;
  path: string;
  category: string;
  desc: string;
  createdAt: string;
}


interface ApiFormType {
  id?: number;
  method: string;
  path: string;
  category: string;
  desc: string;
}

interface ApiTreeType {
  id?: number;
  key: string;
  label: string;
  isLeaf: boolean;
  children?: ApiTreeType[];
}

interface EditRoleApiType {
  id: number;
  addIds: number[];
  removeIds: number[];
}