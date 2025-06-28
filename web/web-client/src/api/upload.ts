import request from '@/utils/request';
import { getFileMD5 } from '@/utils/md5';
import type { AxiosProgressEvent } from 'axios';
import { statusCode } from '@/utils/status-code';

// 上传图片
export const uploadFileAPI = ({ name, file, action, onProgress, onFinish, onError }: UploadOptionsType) => {
  const formData = new FormData();
  formData.append(name, file)
  request.post(action, formData, {
    timeout: 50000,
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

const CHUNK_SIZE = 5 * 1024 * 1024; // 每个分片大小为5MB
const MAX_CONCURRENT_UPLOADS = 5; // 最大并发上传数量
export const uploadFileChunkAPI = async ({ name, file, action, onProgress, onFinish, onError }: UploadOptionsType) => {
  onProgress(0);
  const hash = await getFileMD5(file)

  const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
  const formDataGenerator = createFormDataGenerator(hash, name, file.name, totalChunks.toString());


  const tasks: number[] = []
  const uploadedChunks = await getUploadedChunksAPI(hash)
  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedChunks.includes(i)) {
      tasks.push(i)
    }
  }

  if (tasks.length === 0) {
    finishUploadAPI({ hash, action, onFinish, onError })
    return { controllers: [] };
  }

  // 如果服务器不存在数据则手动上传第一个分片
  if (tasks.length === totalChunks) {
    const chunk = file.slice(0, Math.min(CHUNK_SIZE, file.size));
    const formData = formDataGenerator(chunk, "0");
    const controller = new AbortController();
    const firstChunkRes = await uploadChunkAPI(formData, controller)
    if (firstChunkRes.data.code === statusCode.OK) {
      uploadedChunks.push(0);
      tasks.splice(tasks.indexOf(0), 1);
    } else {
      onError(firstChunkRes.data)
      return { controllers: [controller] };
    }
  }
  // 上传进度
 let taskQueue = [...tasks]
 let currentUploads = 0;
 let uploadedChunksCount = uploadedChunks.length;
 const controllers: AbortController[] = [];

  const processQueue = () => {
    while (currentUploads < MAX_CONCURRENT_UPLOADS && taskQueue.length > 0) {
      const nextTaskIndex = taskQueue.shift();// 从队列中取出下一个任务
      if (nextTaskIndex === undefined) break;
      currentUploads++;
      uploadChunk(nextTaskIndex); // 上传该块
    }
  };

  const uploadChunk = async (i: number) => {
    const start = i * CHUNK_SIZE;
    const end = Math.min(start + CHUNK_SIZE, file.size);
    const chunk = file.slice(start, end);
    const formData = formDataGenerator(chunk, i.toString());
    const controller = new AbortController();
    controllers.push(controller);
    try {
      const res = await uploadChunkAPI(formData, controller);
      if (res.data.code === statusCode.OK) {
        uploadedChunksCount++;
        // 更新进度
        const progress = Math.floor((uploadedChunksCount / totalChunks) * 100);
        if (uploadedChunksCount === totalChunks) {
          finishUploadAPI({ hash, action, onFinish, onError });
        } else {
          onProgress(progress);
        }
      } else {
        onError(res.data);
      }
    } catch (error) {
      if (typeof error === 'object' && error !== null && 'name' in error && (error as any).name === 'CanceledError') {
        // 被取消
      } else {
        onError(error);
      }
    } finally {
      currentUploads--;
      processQueue(); // 尝试处理下一个任务
    }
  };

  // 开始处理队列
  processQueue();
  return { controllers };
}

const createFormDataGenerator = (hash: string, name: string, fileName: string, totalChunks: string) => {
  const savedHash = hash;
  const savedName = name;
  const savedFileName = fileName;
  const savedTotalChunks = totalChunks;

  return (chunk: Blob, i: string) => {
    const formData = new FormData();
    formData.append(savedName, chunk);
    formData.append('hash', savedHash);
    formData.append('name', savedFileName);
    formData.append('chunkIndex', i.toString());
    formData.append('totalChunks', savedTotalChunks);

    return formData;
  };
}

const getUploadedChunksAPI = async (hash: string) => {
  const res = await request.post("v1/upload/checkVideo", { hash }, {})
  if (res.data.code === statusCode.OK) {
    if (res.data.data.chunks) {
      return res.data.data.chunks
    }

    return []
  }

  return []
}

const uploadChunkAPI = (formData: FormData, controller?: AbortController) => {
  return request.post("v1/upload/chunkVideo", formData, controller ? { signal: controller.signal } : undefined)
}

const mergeUploadedChunksAPI = async (hash: string) => {
  const res = await request.post("v1/upload/mergeVideo", { hash }, {})
  if (res.data.code === statusCode.OK) {
    return true
  }

  return false
}

const finishUploadAPI = async ({ hash, action, onFinish, onError }: FinishUploadType) => {
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