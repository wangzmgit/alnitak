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

export const uploadFileChunkAPI = async ({ name, file, action, onProgress, onFinish, onError }: UploadOptionsType) => {
  onProgress(0);
  const hash = await getFileMD5(file)
  const chunkSize = 5 * 1024 * 1024; // 每个分片大小为5MB
  const totalChunks = Math.ceil(file.size / chunkSize);
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
    return;
  }

  // 如果服务器不存在数据则手动上传第一个分片
  if (tasks.length === totalChunks) {
    const chunk = file.slice(0, Math.min(chunkSize, file.size));
    const formData = formDataGenerator(chunk, "0");
    const firstChunkRes = await uploadChunkAPI(formData)
    if (firstChunkRes.data.code === statusCode.OK) {
      uploadedChunks.push(0);
      tasks.splice(tasks.indexOf(0), 1);
    } else {
      onError(firstChunkRes.data)
      return;
    }
  }

  // 上传进度
  let uploadedChunksCount = uploadedChunks.length;
  for (let i of tasks) {
    const start = i * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const chunk = file.slice(start, end);
    const formData = formDataGenerator(chunk, i.toString());
    uploadChunkAPI(formData).then(async (res) => {
      if (res.data.code === statusCode.OK) {
        uploadedChunksCount++
        // 更新进度
        const progress = Math.floor((uploadedChunksCount / totalChunks) * 100);
        if (uploadedChunksCount === totalChunks) {
          finishUploadAPI({ hash, action, onFinish, onError })
        } else {
          onProgress(progress);
        }
      } else {
        onError(res.data);
      }
    })

  }
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