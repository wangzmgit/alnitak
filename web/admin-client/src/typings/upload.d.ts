interface UploadOptionsType {
  name: string;
  file: File;
  action: string;
  onProgress: (percent: number) => void
  onFinish: (data?: any) => void
  onError: (error?: any) => void
}