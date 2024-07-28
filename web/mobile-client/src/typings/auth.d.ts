interface UserLoginType {
  email: string
  password?: string
  code?: string //验证码
  captchaId: string // 人机验证ID
}

interface UserRegisterType {
  email: string
  password: string
  code: string //验证码
  captchaId: string // 人机验证ID
}
