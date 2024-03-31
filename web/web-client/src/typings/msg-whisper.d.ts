interface WhisperType {
  fid: number;
  content: string;
}

interface WhisperListType {
  user: UserInfoType;
  createdAt: string | Date;
  status: boolean;
}

interface WhisperDetailsType {
  fid: number;
  fromId: number;
  content: string;
  createdAt: string;
}