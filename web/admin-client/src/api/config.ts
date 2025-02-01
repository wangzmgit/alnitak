import request from '@/utils/request';

//获取其他配置
export const getOtherConfigAPI = () => {
  return request.get('v1/config/getOtherConfig');
}

//修改其他配置
export const setOtherConfigAPI = (config: OtherConfigType) => {
  return request.post('v1/config/setOtherConfig', config);
}

//获取存储配置
export const getStorageConfigAPI = () => {
  return request.get('v1/config/getStorageConfig');
}

//修改存储配置
export const setStorageConfigAPI = (config: StorageConfigType) => {
  return request.post('v1/config/setStorageConfig', config);
}


//获取邮箱配置
export const getEmailConfigAPI = () => {
  return request.get('v1/config/getEmailConfig');
}

//修改邮箱配置
export const setEmailConfigAPI = (config: EmailConfigType) => {
  return request.post('v1/config/setEmailConfig', config);
}
