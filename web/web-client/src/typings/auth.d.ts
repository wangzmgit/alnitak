interface UserLoginType {
  email: string
  password?: string
  code?: string //验证码
  captchaId: string // 人机验证ID
  rememberMe?: boolean // 记住登录
}

interface UserRegisterType {
  email: string
  password: string
  code: string //验证码
  captchaId: string // 人机验证ID
}
