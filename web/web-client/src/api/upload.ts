import request from '@/utils/request';
import { getFileMD5 } from '@/utils/md5';
import type { AxiosProgressEvent } from 'axios';
import { statusCode } from '@/utils/status-code';

// 配置常量
const CHUNK_SIZE = 5 * 1024 * 1024; // 每个分片大小为5MB
const MAX_CONCURRENT_UPLOADS = 1; // 串行上传，避免服务器连接过载
const MAX_RETRIES = 5; // 最大重试次数
const INITIAL_RETRY_DELAY = 2000; // 初始重试延迟（毫秒）
const CHUNK_TIMEOUT = 120000; // 分片上传超时时间（2分钟）
const MERGE_TIMEOUT = 600000; // 合并超时时间（10分钟，大文件合并需要较长时间）
const OSS_TIMEOUT = 300000; // OSS 直传超时（5分钟）

// 延迟函数
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// 指数退避延迟计算
const getRetryDelay = (retryCount: number) => {
  return INITIAL_RETRY_DELAY * Math.pow(2, retryCount) + Math.random() * 1000;
};

// ═══════════════════════════════════════════════════════════════
// 前端截取视频封面（blob URL 预览）+ 元数据
// ═══════════════════════════════════════════════════════════════
export interface VideoProbe {
  coverURL: string | null;
  duration: number;
  width: number;
  height: number;
  codecName: string;
}

export const probeVideo = async (file: File): Promise<VideoProbe> => {
  const result: VideoProbe = { coverURL: null, duration: 0, width: 0, height: 0, codecName: '' };
  try {
    const video = document.createElement('video');
    video.muted = true;
    video.preload = 'auto';
    const url = URL.createObjectURL(file);
    video.src = url;

    await new Promise<void>((resolve, reject) => {
      video.onloadeddata = () => resolve();
      video.onerror = () => reject(new Error('视频加载失败'));
      setTimeout(() => reject(new Error('视频探测超时')), 10000);
    });

    result.duration = video.duration;
    result.width = video.videoWidth || 1280;
    result.height = video.videoHeight || 720;

    const tracks = (video as any).videoTracks;
    if (tracks && tracks.length > 0) {
      result.codecName = tracks[0].codec || '';
    }

    video.currentTime = Math.min(1, video.duration * 0.1);

    await new Promise<void>((resolve) => {
      video.onseeked = () => resolve();
    });

    const canvas = document.createElement('canvas');
    canvas.width = result.width;
    canvas.height = result.height;
    const ctx = canvas.getContext('2d')!;
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    URL.revokeObjectURL(url);

    const blob = await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((b) => (b ? resolve(b) : reject(new Error('canvas 转换失败'))), 'image/jpeg', 0.85);
    });
    result.coverURL = URL.createObjectURL(blob);
    console.log('【封面截取成功】仅本地预览');
  } catch (e) {
    console.warn('【视频探测失败】', e);
  }
  return result;
};

// ═══════════════════════════════════════════════════════════════
// 图片上传（直传 OSS 优先，local 模式 fallback）
// ═══════════════════════════════════════════════════════════════
export const uploadFileAPI = async ({ name, file, action, onProgress, onFinish, onError }: UploadOptionsType) => {
  onProgress(0);

  // 1. 尝试直传 OSS
  try {
    const hash = await getFileMD5(file);
    const presignRes = await request.post('v1/upload/presignImage', {
      fileName: file.name,
      fileSize: file.size,
    }, { timeout: 10000 });

    if (presignRes.data.code === statusCode.OK && presignRes.data.data?.presignURL) {
      const { presignURL, objectKey } = presignRes.data.data;
      console.log('【直传图片】使用 OSS 直传, objectKey:', objectKey);

      // PUT 直传 OSS
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('PUT', presignURL, true);
        xhr.setRequestHeader('Content-Type', file.type || 'application/octet-stream');
        xhr.upload.onprogress = (e) => {
          if (e.lengthComputable) {
            onProgress(Math.floor((e.loaded / e.total) * 100));
          }
        };
        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            resolve();
          } else {
            reject(new Error(`OSS upload failed: ${xhr.status}`));
          }
        };
        xhr.onerror = () => reject(new Error('OSS upload network error'));
        xhr.send(file);
      });

      // 通知服务器确认
      const confirmRes = await request.post('v1/upload/confirmImage', {
        objectKey,
        hash,
      }, { timeout: 10000 });

      if (confirmRes.data.code === statusCode.OK) {
        onProgress(100);
        onFinish(confirmRes.data);
        return;
      }
    }
  } catch (e) {
    console.warn('【图片直传失败，fallback 到 VPS 代理上传】', e);
  }

  // 2. Fallback: 原始 VPS 代理上传
  const formData = new FormData();
  formData.append(name, file);
  request.post(action, formData, {
    timeout: 50000,
    headers: { 'Content-Type': 'multipart/form-data' },
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
  });
}

// ═══════════════════════════════════════════════════════════════
// 视频分片上传（直传 OSS 优先，local 模式 fallback）
// ═══════════════════════════════════════════════════════════════
export const uploadFileChunkAPI = async ({ name, file, action, replaceResourceID, onProgress, onFinish, onError }: UploadOptionsType) => {
  onProgress(0);
  const hash = await getFileMD5(file);
  const size = file.size;
  const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

  // 后台探测视频（截封面 + 元数据，不阻塞上传）
  const probePromise = probeVideo(file);

  // 检查是否可以秒传
  const { chunks: uploadedChunks, instantUpload, fileID } = await getUploadedChunksAPI(hash, size);
  const probe = await probePromise;
  const onFinishWithCover = (data?: any) => onFinish(data ? { ...data, coverURL: probe.coverURL } : data);

  // 【秒传】文件已存在且转码完成，直接完成上传
  if (instantUpload) {
    console.log('【秒传】文件已存在，跳过上传直接完成, fileID:', fileID);
    onProgress(100);
    try {
      const res = await request.post(action, { hash, fileID, size, duration: probe.duration, width: probe.width, height: probe.height, codecName: probe.codecName, replaceResourceID }, {
        timeout: MERGE_TIMEOUT,
      });
      if (res.data.code === statusCode.OK) {
        onFinishWithCover(res.data);
      } else {
        onError(res.data);
      }
    } catch (error) {
      console.error('秒传完成请求失败', error);
      onError(error);
    }
    return { controllers: [], fileID };
  }

  // 尝试直传 OSS
  try {
    const directResult = await uploadDirectToOSS({
      file, hash, name, size, totalChunks, uploadedChunks, fileID, action,
      replaceResourceID, duration: probe.duration, width: probe.width, height: probe.height, codecName: probe.codecName, onProgress, onFinish: onFinishWithCover, onError,
    });
    if (directResult) return directResult;
  } catch (e) {
    console.warn('【视频直传 OSS 失败，fallback 到 VPS 代理上传】', e);
  }

  // Fallback: 原始 VPS 代理上传流程
  return uploadViaVPS({
    file, hash, name, size, totalChunks, uploadedChunks, fileID, action,
    replaceResourceID, duration: probe.duration, width: probe.width, height: probe.height, codecName: probe.codecName, onProgress, onFinish: onFinishWithCover, onError,
  });
}

// ═══════════════════════════════════════════════════════════════
// 直传 OSS 流程
// ═══════════════════════════════════════════════════════════════
const uploadDirectToOSS = async ({
  file, hash, name, size, totalChunks, uploadedChunks, fileID: _fileID, action,
  replaceResourceID, duration, width, height, codecName, onProgress, onFinish, onError,
}: {
  file: File; hash: string; name: string; size: number; totalChunks: number;
  uploadedChunks: number[]; fileID: string; action: string;
  replaceResourceID?: number; duration?: number; width?: number; height?: number; codecName?: string;
  onProgress: (p: number) => void;
  onFinish: (data: any) => void;
  onError: (err: any) => void;
}): Promise<{ controllers: AbortController[]; fileID: string } | null> => {

  // 1. 初始化分片直传（只返回第一批预签名URL）
  const initRes = await request.post('v1/upload/initVideo', {
    hash, size, fileName: file.name, totalChunks,
  }, { timeout: 15000 });

  if (initRes.data.code !== statusCode.OK || !initRes.data.data) {
    return null; // 非直传模式（oss_type == "local"），fallback
  }

  const initData = initRes.data.data;

  // 秒传
  if (initData.totalChunks === 0) {
    console.log('【直传视频秒传】fileID:', initData.fileID);
    onProgress(100);
    try {
      const res = await request.post(action, { hash, fileID: initData.fileID, size, duration, width, height, codecName, replaceResourceID }, {
        timeout: MERGE_TIMEOUT,
      });
      if (res.data.code === statusCode.OK) {
        onFinish(res.data);
      } else {
        onError(res.data);
      }
    } catch (error) {
      onError(error);
    }
    return { controllers: [], fileID: initData.fileID };
  }

  console.log('【直传视频】使用 OSS 直传, uploadID:', initData.uploadID, 'totalChunks:', totalChunks);

  // 2. 分批上传：每次处理一批，上传完后续签下一批
  const parts: { partNumber: number; etag: string }[] = [];
  let uploadedCount = uploadedChunks.length;
  const controllers: AbortController[] = [];

  // 当前批次的预签名chunks
  let currentChunks: { index: number; partNumber: number; presignURL: string }[] = initData.chunks || [];
  let nextBatchStart: number = initData.nextBatchStart;

  while (currentChunks.length > 0) {
    for (const chunk of currentChunks) {
      // 跳过已上传的分片
      if (uploadedChunks.includes(chunk.index)) {
        continue;
      }

      const start = chunk.index * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const blob = file.slice(start, end);
      const controller = new AbortController();
      controllers.push(controller);

      let success = false;
      for (let retry = 0; retry <= MAX_RETRIES; retry++) {
        try {
          const xhr = new XMLHttpRequest();
          const etag = await new Promise<string>((resolve, reject) => {
            xhr.open('PUT', chunk.presignURL, true);
            xhr.setRequestHeader('Content-Type', 'application/octet-stream');
            xhr.upload.onprogress = () => {};
            xhr.onload = () => {
              if (xhr.status >= 200 && xhr.status < 300) {
                const etagHeader = xhr.getResponseHeader('ETag');
                if (etagHeader) {
                  resolve(etagHeader);
                } else {
                  reject(new Error('Missing ETag header'));
                }
              } else {
                reject(new Error(`OSS chunk upload failed: ${xhr.status}`));
              }
            };
            xhr.onerror = () => reject(new Error('OSS chunk network error'));
            xhr.onabort = () => reject(new Error('Aborted'));
            controller.signal.addEventListener('abort', () => xhr.abort());
            xhr.send(blob);
          });

          parts.push({ partNumber: chunk.partNumber, etag });
          uploadedCount++;
          const progress = Math.floor((uploadedCount / totalChunks) * 100);
          onProgress(progress);
          success = true;
          break;
        } catch (err: any) {
          if (err?.message === 'Aborted') return { controllers, fileID: initData.fileID };
          if (retry < MAX_RETRIES) {
            console.warn(`【直传】分片 ${chunk.index} 失败，${retry + 1}/${MAX_RETRIES} 重试...`);
            await delay(getRetryDelay(retry));
          }
        }
      }

      if (!success) {
        onError({ msg: `分片 ${chunk.index} 上传失败` });
        return { controllers, fileID: initData.fileID };
      }
    }

    // 当前批次上传完毕，续签下一批
    if (nextBatchStart === -1 || nextBatchStart >= totalChunks) {
      break; // 全部签完了
    }

    console.log('【直传视频】续签下一批, start:', nextBatchStart);
    const presignRes = await request.post('v1/upload/presignChunks', {
      fileID: initData.fileID,
      start: nextBatchStart,
      count: 20, // 每批20个
    }, { timeout: 15000 });

    if (presignRes.data.code !== statusCode.OK || !presignRes.data.data) {
      onError({ msg: '续签分片失败' });
      return { controllers, fileID: initData.fileID };
    }

    currentChunks = presignRes.data.data.chunks || [];
    nextBatchStart = presignRes.data.data.nextBatchStart;
  }

  // 3. 通知服务器完成合并
  if (parts.length !== totalChunks) {
    onError({ msg: '分片数量不匹配' });
    return { controllers, fileID: initData.fileID };
  }

  const completeRes = await request.post('v1/upload/completeVideo', {
    fileID: initData.fileID,
    uploadID: initData.uploadID,
    parts,
  }, { timeout: MERGE_TIMEOUT });

  if (completeRes.data.code !== statusCode.OK) {
    onError(completeRes.data);
    return { controllers, fileID: initData.fileID };
  }

  // 4. 创建视频资源
  try {
    const res = await request.post(action, { hash, fileID: initData.fileID, size, duration, width, height, codecName, replaceResourceID }, {
      timeout: MERGE_TIMEOUT,
    });
    if (res.data.code === statusCode.OK) {
      onFinish(res.data);
    } else {
      onError(res.data);
    }
  } catch (error) {
    onError(error);
  }

  return { controllers, fileID: initData.fileID };
}

// ═══════════════════════════════════════════════════════════════
// 原始 VPS 代理上传流程（fallback）
// ═══════════════════════════════════════════════════════════════
const uploadViaVPS = async ({
  file, hash, name, size, totalChunks, uploadedChunks, fileID, action,
  replaceResourceID, duration, width, height, codecName, onProgress, onFinish, onError,
}: {
  file: File; hash: string; name: string; size: number; totalChunks: number;
  uploadedChunks: number[]; fileID: string; action: string;
  replaceResourceID?: number; duration?: number; width?: number; height?: number; codecName?: string;
  onProgress: (p: number) => void;
  onFinish: (data: any) => void;
  onError: (err: any) => void;
}): Promise<{ controllers: AbortController[]; fileID: string }> => {

  const formDataGenerator = createFormDataGenerator(hash, name, file.name, totalChunks.toString(), size);

  const tasks: number[] = [];
  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedChunks.includes(i)) {
      tasks.push(i);
    }
  }

  if (tasks.length === 0) {
    finishUploadAPI({ hash, fileID, size, duration, width, height, codecName, action, replaceResourceID, onFinish, onError });
    return { controllers: [], fileID };
  }

  // 如果服务器不存在数据则手动上传第一个分片（带重试）
  if (tasks.length === totalChunks) {
    const chunk = file.slice(0, Math.min(CHUNK_SIZE, file.size));
    const formData = formDataGenerator(chunk, '0');
    const controller = new AbortController();

    let firstChunkSuccess = false;
    for (let retry = 0; retry <= MAX_RETRIES; retry++) {
      try {
        const firstChunkRes = await uploadChunkAPI(formData, controller);
        if (firstChunkRes.data.code === statusCode.OK) {
          uploadedChunks.push(0);
          tasks.splice(tasks.indexOf(0), 1);
          if (uploadedChunks.length === totalChunks) {
            finishUploadAPI({ hash, fileID, size, duration, width, height, codecName, action, replaceResourceID, onFinish, onError });
            return { controllers: [controller], fileID };
          }
          firstChunkSuccess = true;
          break;
        }
      } catch (error) {
        if (retry < MAX_RETRIES) {
          console.warn(`首个分片上传失败，${retry + 1}/${MAX_RETRIES} 次重试中...`);
          await delay(getRetryDelay(retry));
          continue;
        }
      }
    }

    if (!firstChunkSuccess) {
      onError({ msg: '首个分片上传失败，请检查网络后重试' });
      return { controllers: [controller], fileID };
    }
  }

  let taskQueue = [...tasks];
  let currentUploads = 0;
  let uploadedChunksCount = uploadedChunks.length;
  const controllers: AbortController[] = [];

  const processQueue = () => {
    while (currentUploads < MAX_CONCURRENT_UPLOADS && taskQueue.length > 0) {
      const nextTaskIndex = taskQueue.shift();
      if (nextTaskIndex === undefined) break;
      currentUploads++;
      uploadChunk(nextTaskIndex);
    }
  };

  const uploadChunk = async (i: number, retryCount = 0): Promise<void> => {
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
        const progress = Math.floor((uploadedChunksCount / totalChunks) * 100);
        if (uploadedChunksCount === totalChunks) {
          finishUploadAPI({ hash, fileID, size, duration, width, height, codecName, action, replaceResourceID, onFinish, onError });
        } else {
          onProgress(progress);
        }
      } else {
        if (retryCount < MAX_RETRIES) {
          console.warn(`分片 ${i} 上传失败，${retryCount + 1}/${MAX_RETRIES} 次重试中...`);
          await delay(getRetryDelay(retryCount));
          return uploadChunk(i, retryCount + 1);
        }
        onError(res.data);
      }
    } catch (error) {
      if (typeof error === 'object' && error !== null && 'name' in error && (error as any).name === 'CanceledError') {
        return;
      }
      if (retryCount < MAX_RETRIES) {
        console.warn(`分片 ${i} 网络错误，${retryCount + 1}/${MAX_RETRIES} 次重试中...`, error);
        await delay(getRetryDelay(retryCount));
        return uploadChunk(i, retryCount + 1);
      }
      console.error(`分片 ${i} 上传失败，已达到最大重试次数`);
      onError(error);
    } finally {
      currentUploads--;
      processQueue();
    }
  };

  processQueue();
  return { controllers, fileID };
}

// ═══════════════════════════════════════════════════════════════
// 工具函数
// ═══════════════════════════════════════════════════════════════

const createFormDataGenerator = (hash: string, name: string, fileName: string, totalChunks: string, size: number) => {
  const savedHash = hash;
  const savedName = name;
  const savedFileName = fileName;
  const savedTotalChunks = totalChunks;
  const savedSize = size;

  return (chunk: Blob, i: string) => {
    const formData = new FormData();
    formData.append(savedName, chunk);
    formData.append('hash', savedHash);
    formData.append('name', savedFileName);
    formData.append('chunkIndex', i.toString());
    formData.append('totalChunks', savedTotalChunks);
    formData.append('size', savedSize.toString());

    return formData;
  };
}

// 检查上传进度，返回 { chunks: number[], fileID: string, instantUpload: boolean }
const getUploadedChunksAPI = async (hash: string, size: number): Promise<{ chunks: number[], fileID: string, instantUpload: boolean }> => {
  const res = await request.post('v1/upload/checkVideo', { hash, size }, {});
  if (res.data.code === statusCode.OK) {
    const chunks = res.data.data.chunks || [];
    const fileID: string = res.data.data.fileID || '';
    if (chunks.length === 1 && chunks[0] === -1) {
      return { chunks: [], fileID, instantUpload: true };
    }
    return { chunks, fileID, instantUpload: false };
  }

  return { chunks: [], fileID: '', instantUpload: false };
}

const uploadChunkAPI = (formData: FormData, controller?: AbortController) => {
  return request.post('v1/upload/chunkVideo', formData, {
    timeout: CHUNK_TIMEOUT,
    signal: controller?.signal,
  });
}

const mergeUploadedChunksAPI = async (hash: string, fileID: string, size: number, retryCount = 0): Promise<boolean> => {
  try {
    const res = await request.post('v1/upload/mergeVideo', { hash, fileID, size }, {
      timeout: MERGE_TIMEOUT,
    });
    if (res.data.code === statusCode.OK) {
      return true;
    }
    if (retryCount < MAX_RETRIES) {
      console.warn(`合并分片失败，${retryCount + 1}/${MAX_RETRIES} 次重试中...`);
      await delay(getRetryDelay(retryCount));
      return mergeUploadedChunksAPI(hash, fileID, size, retryCount + 1);
    }
    return false;
  } catch (error) {
    if (retryCount < MAX_RETRIES) {
      console.warn(`合并分片网络错误，${retryCount + 1}/${MAX_RETRIES} 次重试中...`);
      await delay(getRetryDelay(retryCount));
      return mergeUploadedChunksAPI(hash, fileID, size, retryCount + 1);
    }
    console.error('合并分片失败，已达到最大重试次数');
    return false;
  }
}

const finishUploadAPI = async ({ hash, fileID, size, duration, width, height, codecName, action, replaceResourceID, onFinish, onError }: FinishUploadType) => {
  if (await mergeUploadedChunksAPI(hash, fileID, size)) {
    try {
      const res = await request.post(action, { hash, fileID, size, duration, width, height, codecName, replaceResourceID }, {
        timeout: MERGE_TIMEOUT,
      });
      if (res.data.code === statusCode.OK) {
        onFinish(res.data);
      } else {
        onError();
      }
    } catch (error) {
      console.error('完成上传请求失败', error);
      onError();
    }
  } else {
    onError();
  }
}
