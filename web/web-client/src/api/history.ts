import request from '@/utils/request';

// 上传历史记录
export const addHistoryAPI = (addHistory: AddHistoryType) => {
  return request.post('v1/history/addHistory', addHistory);
}

// 获取播放进度
export const getHistoryProgressAPI = (vid: number, part: number) => {
  return request.get(`v1/history/getProgress?vid=${vid}&part=${part}`);
}

// 获取历史记录
export const getHistoryVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/history/getHistory?page=${page}&pageSize=${pageSize}`);
}

