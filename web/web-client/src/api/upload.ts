import request from '@/utils/request';
import { getFileMD5 } from '@/utils/md5';
import type { AxiosProgressEvent } from 'axios';
import { statusCode } from '@/utils/status-code';

// 上传图片
export const uploadFileAPI = ({ name, file, action, onProgress, onFinish, onError }: UploadOptionsType) => {
  const formData = new FormData();
  formData.append(name, file)
  request.post(action, formData, {
    // 文件上传20分钟超时 1200000 = 1000*60*20
    timeout: 1200000,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent: AxiosProgressEvent) => {
      if (!progressEvent.total) {
        onProgress(0);
        return;
      }
      onProgress(Math.floor(progressEvent.loaded / progressEvent.total * 100));
    }
  }).then((res) => {
    if (res.data.code === statusCode.OK) {
      onFinish(res.data);
    } else {
      onError(res.data);
    }
  }).catch((err) => {
    onError(err);
  })
}

const getUploadedChunksAPI = async (hash: string) => {
  const res = await request.post("v1/upload/video/check", { hash }, {})
  if (res.data.code === statusCode.OK) {
    if (res.data.data.chunks) {
      return res.data.data.chunks
    }

    return []
  }

  return []
}

const mergeUploadedChunksAPI = async (hash: string) => {
  const res = await request.post("v1/upload/video/merge", { hash }, {})
  if (res.data.code === statusCode.OK) {
    return true
  }

  return false
}

const finishUploadAPI = async (hash: string, action: string, onFinish: (data?: any) => void, onError: (data?: any) => void) => {
  if (await mergeUploadedChunksAPI(hash)) {
    request.post(action, { hash }, {}).then((res) => {
      if (res.data.code === statusCode.OK) {
        onFinish(res.data);
      } else {
        onError();
      }
    })
  } else {
    onError();
  }
}

export const uploadFileChunkAPI = async (param: UploadOptionsType) => {
  param.onProgress(0);
  const hash = await getFileMD5(param.file)
  const chunkSize = 5 * 1024 * 1024; // 每个分片大小为5MB
  const totalChunks = Math.ceil(param.file.size / chunkSize);

  const tasks: number[] = []
  const uploadedChunks = await getUploadedChunksAPI(hash)
  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedChunks.includes(i)) {
      tasks.push(i)
    }
  }

  if (tasks.length === 0) {
    finishUploadAPI(hash, param.action, param.onFinish, param.onError)
  }

  // 上传进度
  let uploadedChunksCount = uploadedChunks.length;
  for (let i of tasks) {
    const start = i * chunkSize;
    const end = Math.min(start + chunkSize, param.file.size);
    const chunk = param.file.slice(start, end);

    const formData = new FormData();
    formData.append(param.name, chunk);
    formData.append('hash', hash);
    formData.append('name', param.file.name);
    formData.append('chunkIndex', i.toString());
    formData.append('totalChunks', totalChunks.toString());

    request.post("v1/upload/video/chunk", formData, {}).then(async (res) => {
      if (res.data.code === statusCode.OK) {
        uploadedChunksCount++
        // 更新进度
        const progress = Math.floor((uploadedChunksCount / totalChunks) * 100);
        if (uploadedChunksCount === totalChunks) {
          finishUploadAPI(hash, param.action, param.onFinish, param.onError)
        } else {
          param.onProgress(progress);
        }
      } else {
        param.onError(res.data);
      }
    })
  }
}