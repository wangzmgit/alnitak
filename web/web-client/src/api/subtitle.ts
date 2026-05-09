import request from '@/utils/request';

/** 分 P 字幕列表（可选登录：作者可预览未发布分 P） */
export const getSubtitleListAPI = (resourceShortId: string) => {
  return request.get('v1/video/subtitle/list', {
    params: { resourceShortId },
  });
};

/** 上传字幕（multipart：resourceShortId、lang、label?、isDefault?、file） */
export const uploadSubtitleAPI = (formData: FormData) => {
  return request.post('v1/video/subtitle/upload', formData);
};

/** 更新字幕元数据或文件（JSON 或 multipart） */
export const updateSubtitleAPI = (id: number, data: Record<string, unknown> | FormData) => {
  if (data instanceof FormData) {
    return request.put(`v1/video/subtitle/${id}`, data);
  }
  return request.put(`v1/video/subtitle/${id}`, data);
};

export const deleteSubtitleAPI = (id: number) => {
  return request.delete(`v1/video/subtitle/${id}`);
};
