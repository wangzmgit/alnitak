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

//获取转码配置
export const getTranscodingConfigAPI = () => {
  return request.get('v1/config/getTranscodingConfig');
}

//修改转码配置
export const setTranscodingConfigAPI = (config: TranscodingConfigType) => {
  return request.post('v1/config/setTranscodingConfig', config);
}

//获取资源清理预览
export const getCleanupPreviewAPI = () => {
  return request.get('v1/config/getCleanupPreview');
}

//执行资源清理
export const executeCleanupAPI = () => {
  return request.post('v1/config/executeCleanup');
}
