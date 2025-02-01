import request from '@/utils/request';

// 获取api列表
export const getApiListAPI = (data: ApiListParam) => {
  return request.post("v1/api/getApiList", data);
}

// 获取全部api列表
export const getAllApiListAPI = () => {
  return request.get("v1/api/getAllApiList");
}

// 新增API
export const addApiAPI = (data: ApiFormType) => {
  return request.post('v1/api/addApi', data);
}

// 编辑API
export const editApiAPI = (data: ApiFormType) => {
  return request.put('v1/api/editApi', data);
}

// 删除API
export const deleteApiAPI = (id: number) => {
  return request.delete(`v1/api/deleteApi/${id}`);
}

// 获取角色API
export const getRoleApiAPI = (code: string) => {
  return request.get(`v1/api/getRoleApi?code=${code}`);
}

// 编辑角色API
export const editRoleApiAPI = (data: EditRoleApiType) => {
  return request.put('v1/api/editRoleApi', data);
}