package cache

import "time"

// 邮箱验证码缓存标识符
const EMAIL_CODE_KEY = "email_code_key:"

// 邮箱验证码过期时间 n 分钟
const EMIAL_CODE_EXPIRATION_TIME = time.Minute * time.Duration(5)

// 滑块x坐标缓存标识符
const SLIDER_X_KEY = "slider_x_key:"

// 滑块x坐标过期时间 n 分钟
const SLIDER_X_EXPRIRATION_TIME = time.Minute * time.Duration(5)

// 人机验证状态
const CAPTCHA_STATUS_KEY = "captcha_status_key:"

// 人机验证状态过期时间 n 分钟
const CAPTCHA_STATUS_EXPRIRATION_TIME = time.Minute * time.Duration(30)

// 登录尝试次数缓存标识符
const LOGIN_TRY_COUNT_KEY = "login_try_count_key:"

// 登录尝试次数过期时间 n 分钟
const LOGIN_TRY_COUNT_EXPRIRATION_TIME = time.Minute * time.Duration(30)

// 最大登录数量
const MAX_LOGIN_LIMIT = 3

// 刷新token缓存标识符
const REFRESH_TOKEN_KEY = "refresh_token_key:"

// 刷新token过期时间 n 小时
const REFRESH_TOKEN_EXPRIRATION_TIME = time.Hour * time.Duration(14*24) // 14 * 24

// 验证token过期时间 n 分钟
const ACCESS_TOKEN_EXPRIRATION_TIME = time.Minute * time.Duration(10)

// 用户信息缓存标识符
const USER_INFO_KEY = "user_info_key:"

// 用户信息过期时间 n 小时
const USER_INFO_EXPIRATION_TIME = time.Hour * time.Duration(24)

// 上传文件缓存标识符
const UPLOAD_IMAGE_KEY = "upload_image_key:"

// 上传文件过期时间 n 分钟
const UPLOAD_IMAGE_EXPRIRATION_TIME = time.Minute * time.Duration(20)

// 分区缓存标识符
const PARTITION_KEY = "partition_key:"

// 分区过期时间 不过期
const PARTITION_EXPRIRATION_TIME = 0
