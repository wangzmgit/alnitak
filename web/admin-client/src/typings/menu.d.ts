interface MenuType {
  id: number;
  name: string;
  path: string;
  desc: string;
  parentId: number;
  component: string;
  meta: MenuMetaType;
  children?: MenuType[];
}

interface MenuMetaType {
  title: string;
  hidden: boolean;
  disabled: boolean;
  keepAlive: boolean;
  icon: keyof typeof import("@vicons/ionicons5")
}

interface MenuTreeType {
  key: number;
  id: number;
  title: string;
  children: MenuTreeType[] | null;
}

interface MenuItemType {
  id: number;
  name: string;
  path: string;
  desc: string;
  icon?: any;
  sort: number;
  component: string;
  meta: MenuMetaType;
  children: MenuItemType[];
  parentId: number;
}

interface MenuMetaType {
  title: string;
  icon: string;
  hidden: boolean;
  keepAlive: boolean;
}

interface MenuFormType {
  id?: number;
  name: string;
  path: string;
  desc: string;
  sort: number;
  component: string;
  title: string;
  icon: string;
  hidden: boolean;
  keepAlive: boolean;
  parentId?: number;
}

interface EditRoleMenuType {
  id: number;
  menuIds: number[];
}