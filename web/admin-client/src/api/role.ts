import request from '@/utils/request';

// 获取个人角色信息
export const getRoleInfoAPI = () => {
  return request.get("v1/role/getRoleInfo");
}

// 获取角色列表
export const getRoleListAPI = (data: RoleListParam) => {
  return request.post("v1/role/getRoleList", data);
}

// 新增角色
export const addRoleAPI = (data: RoleFormType) => {
  return request.post('v1/role/addRole', data);
}

// 编辑角色
export const editRoleAPI = (data: RoleFormType) => {
  return request.put('v1/role/editRole', data);
}

// 编辑角色首页
export const editRoleHomeAPI = (id: number, home: string) => {
  return request.put('v1/role/editRoleHome', { id, home });
}

// 删除角色
export const deleteRoleAPI = (id: number) => {
  return request.delete(`v1/role/deleteRole/${id}`);
}

// 获取全部角色列表
export const getAllRoleListAPI = () => {
  return request.get("v1/role/getAllRoleList");
}
