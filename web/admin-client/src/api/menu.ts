import request from '@/utils/request';

// 获取用户菜单
export const getUserMenuAPI = () => {
  return request.get('v1/menu/getUserMenu');
}

// 获取菜单
export const getMenuAPI = () => {
  return request.get('v1/menu/getMenuTree');
}

// 新增菜单
export const addMenuAPI = (data: MenuFormType) => {
  return request.post('v1/menu/addMenu', data);
}

// 编辑菜单
export const editMenuAPI = (data: MenuFormType) => {
  return request.put('v1/menu/editMenu', data);
}

// 删除菜单
export const deleteMenuAPI = (id: number) => {
  return request.delete(`v1/menu/deleteMenu/${id}`);
}

// 获取角色菜单
export const getRoleMenuAPI = (code: string) => {
  return request.get(`v1/menu/getRoleMenu?code=${code}`);
}

// 编辑角色菜单
export const editRoleMenuAPI = (data: EditRoleMenuType) => {
  return request.put('v1/menu/editRoleMenu', data);
}