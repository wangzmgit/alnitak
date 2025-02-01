interface OtherConfigType {
  allowOrigin: string;
  prefix: string;
}

interface EmailConfigType {
  addresser: string;
  host: string;
  port: number;
  user: string;
  pass: string;
}

interface StorageConfigType {
  maxImgSize: number;
  maxVideoSize: number;
  type: string;
  keyId: string;
  keySecret: string;
  bucket: string;
  endpoint: string;
  appId: string;
  region: string;
  domain: string;
  private: boolean;
  uploadMp4File: boolean;
}