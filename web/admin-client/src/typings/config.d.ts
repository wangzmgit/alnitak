interface OtherConfigType {
  allowOrigin: string;
  prefix: string;
  generate1080p60: boolean;
  useGpu: boolean;
  useH265: boolean;
  useAv1: boolean;
  serverPort: string;
  sslEnabled: boolean;
  sslPort: string;
  sslCertFile: string;
  sslKeyFile: string;
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