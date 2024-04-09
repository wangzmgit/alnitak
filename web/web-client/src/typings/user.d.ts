interface UserInfoType {
  uid: number
  name: string
  avatar: string
  spacecover?: string
  email?: string
  gender?: number
  sign?: string
  birthday?: string
  createdAt?: string
  role?: number
}

interface EditUserInfoType {
  avatar: string;
  name: string;
  gender?: number;
  sign?: string;
  birthday?: string;
}

interface ModifyPwdType {
  email: string;
  password: string;
  code: string; //验证码
  captchaId: string;
}