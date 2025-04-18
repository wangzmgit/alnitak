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
  const hash = await getFileMD5(file);
  const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
  const formDataGenerator = createFormDataGenerator(hash, name, file.name, totalChunks.toString());

  // 构造待上传的分片索引队列
  const uploadedChunks = await getUploadedChunksAPI(hash);
  const tasks: number[] = [];
  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedChunks.includes(i)) tasks.push(i);
  }

  // 全部已上传，直接完成
  if (tasks.length === 0) {
    finishUploadAPI({ hash, action, onFinish, onError });
    return;
  }

  // 如果一次都没上传过，先上传第一个分片
  if (tasks.length === totalChunks) {
    const chunk = file.slice(0, Math.min(CHUNK_SIZE, file.size));
    const firstRes = await uploadChunkAPI(formDataGenerator(chunk, '0'));
    if (firstRes.data.code === statusCode.OK) {
      uploadedChunks.push(0);
      tasks.shift(); // 移除 0
    } else {
      onError(firstRes.data);
      return;
    }
  }

  // 并发上传控制
  let taskQueue = [...tasks];
  let currentUploads = 0;
  let uploadedCount = uploadedChunks.length;

  const processQueue = () => {
    while (currentUploads < MAX_CONCURRENT_UPLOADS && taskQueue.length) {
      const idx = taskQueue.shift()!;
      currentUploads++;
      uploadChunk(idx);
    }
  };

  // 单分片上传 + 重试 + 收尾
  const uploadChunk = async (i: number) => {
    const MAX_RETRIES = 3;
    let success = false;

    for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
      const start = i * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const chunk = file.slice(start, end);
      const formData = formDataGenerator(chunk, i.toString());

      try {
        const res = await uploadChunkAPI(formData);
        if (res.data.code === statusCode.OK) {
          // 上传成功
          success = true;
          uploadedCount++;
          const progress = Math.floor((uploadedCount / totalChunks) * 100);
          if (uploadedCount === totalChunks) {
            finishUploadAPI({ hash, action, onFinish, onError });
          } else {
            onProgress(progress);
          }
          break; // 跳出重试循环
        } else {
          // 业务层面失败，不再重试
          onError(res.data);
          break;
        }
      } catch (err) {
        if (attempt < MAX_RETRIES) {
          console.warn(`分片 ${i} 上传失败，第 ${attempt} 次重试...`);
        } else {
          // 重试耗尽
          onError(err);
        }
      }
    }

    // 重试循环结束后（不论成功或失败）才执行收尾
    currentUploads--;
    processQueue();
  };

  // 启动队列
  processQueue();
};

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

const uploadChunkAPI = (formData: FormData) => {
  return request.post("v1/upload/chunkVideo", formData)
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
