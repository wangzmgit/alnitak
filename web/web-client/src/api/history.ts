import request from '@/utils/request';

// 上传历史记录
export const addHistoryAPI = (addHistory: AddHistoryType) => {
  return request.post('v1/history/video/addHistory', addHistory);
}

// 获取播放进度
export const getHistoryProgressAPI = (vid: number, part?: number) => {
  let url = `v1/history/video/getProgress?vid=${vid}`;
  if (part !== undefined && part !== null) {
    url += `&part=${part}`;
  }
  return request.get(url);
}

// 获取历史记录
export const getHistoryVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/history/video/getHistory?page=${page}&pageSize=${pageSize}`);
}

