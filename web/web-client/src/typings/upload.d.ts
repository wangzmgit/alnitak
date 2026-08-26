interface UploadOptionsType {
  name: string;
  file: File;
  action: string;
  replaceResourceID?: number; // 替换资源时传旧资源ID
  onProgress: (percent: number) => void
  onFinish: (data?: any) => void
  onError: (error?: any) => void
}

interface FinishUploadType {
  hash: string;
  fileID: string;
  size: number;
  duration?: number;
  width?: number;
  height?: number;
  codecName?: string;
  action: string;
  replaceResourceID?: number;
  onFinish: (data?: any) => void
  onError: (error?: any) => void
}