import request from '@/utils/request';
import { getFileMD5 } from '@/utils/md5';
import type { AxiosProgressEvent } from 'axios';
import { statusCode } from '@/utils/status-code';
import { AxiosError } from 'axios'; // 导入 AxiosError 类型
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

// 合并视频分片
const mergeUploadedChunksAPI = async (hash: string) => {
  try {
    const res = await request.post("v1/upload/video/merge", { hash }, {});
    if (res.data.code === statusCode.OK) {
      return true;
    } else {
      throw new Error('合成失败');
    }
  } catch (error) {
    // 使用类型断言，告诉 TypeScript error 是 AxiosError 类型
    if ((error as AxiosError).isAxiosError) {
      const axiosError = error as AxiosError; // 类型断言为 AxiosError
      console.error(`Axios 错误: ${axiosError.message}, 状态码: ${axiosError.response?.status}`);
    } else {
      console.error(`未知错误: ${(error as Error).message}`);
    }

    // 直接抛出错误，不再重试
    console.error('合成分片失败，未进行重试');
    return false;
  }
};

const finishUploadAPI = async (hash: string, action: string, onFinish: (data?: any) => void, onError: (data?: any) => void) => {
  const mergeSuccess = await mergeUploadedChunksAPI(hash); // 调用合成分片函数并等待结果
  if (mergeSuccess) {
    request.post(action, { hash }, {}).then((res) => {
      if (res.data.code === statusCode.OK) {
        onFinish(res.data);
      } else {
        onError();
      }
    }).catch((err) => onError(err));
  } else {
    onError();
  }
}

const uploadLock = new Map<number, Promise<void>>(); // 用于存储每个分片的锁

// 上传视频分片
export const uploadFileChunkAPI = async (param: UploadOptionsType) => {
  param.onProgress(0);
  const hash = await getFileMD5(param.file);
  const chunkSize = 5 * 1024 * 1024; // 每个分片大小为5MB
  const totalChunks = Math.ceil(param.file.size / chunkSize);
  const tasks: number[] = [];
  const uploadedChunks = await getUploadedChunksAPI(hash);

  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedChunks.includes(i)) {
      tasks.push(i);
    }
  }

  if (tasks.length === 0) {
    finishUploadAPI(hash, param.action, param.onFinish, param.onError);
    return;
  }

  // 上传进度
  let uploadedChunksCount = uploadedChunks.length;

  // 上传分片
  const uploadChunk = async (i: number) => {
    // 检查是否有锁
    if (uploadLock.has(i)) {
      await uploadLock.get(i); // 等待当前分片的上传完成
    }

    // 创建锁
    const uploadPromise = new Promise<void>(async (resolve, reject) => {
      const start = i * chunkSize;
      const end = Math.min(start + chunkSize, param.file.size);
      const chunk = param.file.slice(start, end);

      const formData = new FormData();
      formData.append(param.name, chunk);
      formData.append('hash', hash);
      formData.append('name', param.file.name);
      formData.append('chunkIndex', i.toString());
      formData.append('totalChunks', totalChunks.toString());

      try {
        const res = await request.post("v1/upload/video/chunk", formData, {
          onUploadProgress: (progressEvent: AxiosProgressEvent) => {
            if (!progressEvent.total) return;
            const progress = Math.floor((progressEvent.loaded / progressEvent.total) * 100);
            console.log(`Chunk ${i} upload progress: ${progress}%`);
          },
        });

        if (res.data.code === statusCode.OK) {
          uploadedChunksCount++;
          const progress = Math.floor((uploadedChunksCount / totalChunks) * 100);
          if (uploadedChunksCount === totalChunks) {
            finishUploadAPI(hash, param.action, param.onFinish, param.onError);
          } else {
            param.onProgress(progress);
          }
          resolve(); // 成功上传，解除锁
        } else {
          param.onError(res.data);
          reject(new Error(`分片 ${i} 上传失败`)); // 直接抛出错误，不再重试
        }
      } catch (error) {
        console.error(`分片 ${i} 上传失败:`, error);
        param.onError(error);
        reject(new Error(`分片 ${i} 上传失败`)); // 直接抛出错误，不再重试
      }
    });

    uploadLock.set(i, uploadPromise); // 设置锁
    await uploadPromise; // 等待上传完成
    uploadLock.delete(i); // 上传完成后解除锁
  };

  // 批量上传
  const uploadBatch = async (batch: number[]) => {
    const uploadPromises = batch.map((i) => uploadChunk(i));
    await Promise.all(uploadPromises); // 等待当前批次上传完成
  };

  // 分批上传
  const batchSize = 3; // 每批次上传的分片数量
  for (let i = 0; i < tasks.length; i += batchSize) {
    const currentBatch = tasks.slice(i, i + batchSize);
    try {
      await uploadBatch(currentBatch); // 上传当前批次
    } catch (error) {
      console.error('上传失败', error);
      param.onError(error);
      break;
    }
  }
};