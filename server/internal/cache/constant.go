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
const REFRESH_TOKEN_EXPRIRATION_TIME = time.Hour * time.Duration(7*24)

// 刷新token缓冲时间 n 小时
const REFRESH_TOKEN_BUFFER_TIME = time.Hour * time.Duration(24)

// 验证token过期时间 n 分钟
const ACCESS_TOKEN_EXPRIRATION_TIME = time.Minute * time.Duration(60)

// 用户信息缓存标识符
const USER_INFO_KEY = "user_info_key:"

// 用户信息过期时间 n 小时
const USER_INFO_EXPIRATION_TIME = time.Hour * time.Duration(120)

// 上传文件缓存标识符
const UPLOAD_IMAGE_KEY = "upload_image_key:"

// 上传文件过期时间 n 分钟
const UPLOAD_IMAGE_EXPRIRATION_TIME = time.Minute * time.Duration(20)

// 视频分区缓存标识符
const VIDEO_PARTITION_KEY = "video_partition_key"

// 视频分区过期时间 不过期
const VIDEO_PARTITION_EXPRIRATION_TIME = 0

// 文章分区缓存标识符
const ARTICLE_PARTITION_KEY = "article_partition_key"

// 文章分区过期时间 不过期
const ARTICLE_PARTITION_EXPRIRATION_TIME = 0

// 视频切片状态标识符
const VIDEO_SLICE_STATUS = "video_slice_key:"

// 视频切片状态未使用过期时间 n 小时
const VIDEO_SLICE_EXPRIRATION_TIME = time.Hour * time.Duration(24)

// 视频点击量限制标识符
const VIDEO_CLICKS_LIMIT_KEY = "video_clicks_limit_key:"

// 视频点击量限制过期时间  n 分钟
const VIDEO_CLICKS_LIMIT_EXPRIRATION_TIME = time.Minute * time.Duration(60)

// 视频点击量标识符
const VIDEO_CLICKS_KEY = "video_clicks_key:"

// 点击量过期时间  n 小时
const VIDEO_CLICKS_EXPRIRATION_TIME = time.Hour * time.Duration(24)

// 视频信息缓存标识符
const VIDEO_INFO_KEY = "video_info_key:"

// 用户信息过期时间 n 小时
const VIDEO_INFO_EXPIRATION_TIME = time.Hour * time.Duration(120)

// 重置密码验证状态
const RESET_PWD_CHECK_KEY = "reset_pwd_check_key:"

// 重置密码验证状态过期时间 n 分钟
const RESET_PWD_CHECK_EXPRIRATION_TIME = time.Minute * time.Duration(30)

// 全部视频的ID
const ALL_VIDEO_KEY = "all_video_key:"

// 热门视频ID
const HOT_VIDEO_KEY = "hot_video_key"

// 全部文章的ID
const ALL_ARTICLE_KEY = "all_article_key:"

// 文章点击量限制标识符
const ARTICLE_CLICKS_LIMIT_KEY = "article_clicks_limit_key:"

// 文章点击量限制过期时间  n 分钟
const ARTICLE_CLICKS_LIMIT_EXPRIRATION_TIME = time.Minute * time.Duration(60)

// 文章点击量标识符
const ARTICLE_CLICKS_KEY = "article_clicks_key:"

// 文章量过期时间  n 小时
const ARTICLE_CLICKS_EXPRIRATION_TIME = time.Hour * time.Duration(24)
