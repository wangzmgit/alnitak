interface UploadOptionsType {
  name: string;
  file: File;
  action: string;
  onProgress: (percent: number) => void
  onFinish: (data?: any) => void
  onError: (error?: any) => void
}

interface FinishUploadType {
  hash: string;
  action: string;
  onFinish: (data?: any) => void
  onError: (error?: any) => void
}