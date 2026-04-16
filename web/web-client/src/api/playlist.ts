import request from '@/utils/request';

// 创建合集
export const addPlaylistAPI = (data: { title: string, cover: string, desc: string }) => {
  return request.post('v1/playlist/add', data);
}

// 编辑合集
export const editPlaylistAPI = (data: { id: number, title: string, cover: string, desc: string }) => {
  return request.put('v1/playlist/edit', data);
}

// 删除合集
export const deletePlaylistAPI = (id: number) => {
  return request.delete(`v1/playlist/del/${id}`);
}

// 获取自己的合集列表
export const getMyPlaylistAPI = () => {
  return request.get('v1/playlist/myList');
}

// 添加视频到合集
export const addPlaylistVideoAPI = (data: { playlistId: number, vids: number[] }) => {
  return request.post('v1/playlist/video/add', data);
}

// 从合集移除视频
export const delPlaylistVideoAPI = (data: { playlistId: number, vids: number[] }) => {
  return request.post('v1/playlist/video/del', data);
}

// 调整合集视频排序
export const sortPlaylistVideoAPI = (data: { playlistId: number, vids: number[] }) => {
  return request.post('v1/playlist/video/sort', data);
}

// 获取合集视频列表
export const getPlaylistVideoListAPI = (id: number, page: number, pageSize: number) => {
  return request.get(`v1/playlist/video/list?playlistId=${id}&page=${page}&pageSize=${pageSize}`);
}

// 获取合集审核记录
export const getPlaylistReviewRecordAPI = (id: number) => {
  return request.get(`v1/playlist/getPlaylistReviewRecord?id=${id}`);
}

// 获取用户所有合集中的视频ID映射
export const getMyPlaylistVideoIdsAPI = () => {
  return request.get('v1/playlist/video/myVideoIds');
}

// 获取视频所属的合集列表
export const getVideoPlaylistsAPI = (vid: number | string) => {
  return request.get(`v1/playlist/video/playlists?vid=${vid}`);
}

// 获取合集视频列表（包含多分P展开）
export const getPlaylistVideoListWithPartsAPI = (vid: number | string) => {
  return request.get(`v1/playlist/video/listWithParts?vid=${vid}`);
}
