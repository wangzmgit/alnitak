# 项目配置

## 后端项目配置
:::tip 配置文件说明
在后端项目中，根据开发环境和生产环境的不同，我们将配置文件分为了两个。默认情况下，会使用生产环境配置文件。
如果需要使用开发环境配置文件，可以在启动命令中添加参数`-env dev`。
:::
### 配置文件
配置文件位于目录`server/conf`，包括两个文件：`application.dev.yaml（用于开发环境）`和 `application.prod.yaml（用于生产环境）`。
配置文件内容如下：
```yaml
cors:
  # 跨域配置，存在多个使用 , 分隔，请根据实际状况填写，(格式：http://example.com)
  # 建议部署完成后在后台管理中配置
  allow_origin: "*"
file:
  # 上传图片文件大小限制，单位M
  max_img_size: 5
  # 上传视频文件大小限制，单位M
  max_video_size: 500
log: # 日志相关配置，不建议改动
  filename: ./logs/app.log
  max-age: 60
  max-backups: 10
  max-size: 20
  mode: prod
mail:
  # 邮箱发送者
  addresser: 网站管理员
  # debug为true只会在日志中输出验证码，不会实际发送邮件
  debug: false
  # SMTP服务器地址
  host: smtp.163.com
  # 邮箱秘钥(申请smtp得到的秘钥，不是邮箱密码)
  pass: 
  # SMTP服务器端口
  port: 465
  # 发送者的邮箱(格式：noreply@example.com)
  user: 
mysql:
  # 数据库名称
  datasource: alnitak-dev
  # 数据库地址
  host: 127.0.0.1
  # 数据库参数，不建议改动
  param: charset=utf8mb4&parseTime=True&loc=Local
  # 数据库密码
  password: 
  # 数据库端口
  port: 3306
  # 数据库用户名
  username: root
redis:
  # Redis地址
  host: 127.0.0.1
  # Redis端口
  port: 6379
  # Redis密码
  password: 
security:
  # AccessJwt秘钥，可不填写会自动生成
  access_jwt_secret: 
  # RefreshJwt秘钥，可不填写会自动生成
  refresh_jwt_secret: 
  # 是否关闭用户API调用记录
  close_record_user_operation: true
storage:
  # oss的bucket，非local类型必填
  bucket: 
  # 阿里云oss必填，例如：oss-cn-beijing.aliyuncs.com
  endpoint: 
  # oss非local类型必填，对应阿里云access_key_id
  key_id: 
  # oss非local类型必填，对应阿里云access_key_secret
  key_secret: 
  # oss类型(local服务器存储，aliyun阿里云)
  oss_type: local
  # 是否将视频的原始mp4上传到oss
  upload_mp4_file: false
user:
  # 用户注册时生成用户名的默认前缀
  prefix: user_

```

## 前端项目配置

### 用户端配置
配置文件位于目录`web/web-client/src/utils/global-config.ts`，配置文件内容如下：
```js
const title = "弹幕网站标题"; // 网站标题
const https = false; // 是否使用https
const domain = "localhost:9000"; // 后端地址
const icp = "icp备案信息"; // icp备案信息
const security = "公网安备信息"; // 公网安备信息
const keywords = "视频,弹幕"; // 网站关键词
const description = "这里是介绍"; // 网站介绍
const article = true; // 是否开启专栏模块(1.0.3新增)

//上传文件大小限制，需要先修改后端大小限制
const maxImgSize = 5;//上传图片最大大小(单位M)
const maxVideoSize = 500;//上传视频最大大小(单位M)
```

### 管理端配置
配置文件位于目录`web/admin-client/src/utils/global-config.ts`，配置文件内容如下：
```js
const title = "弹幕网站标题"; // 网站标题
const https = false; // 是否使用https
const domain = "localhost:9000"; // 后端地址
const icp = "icp备案信息"; // icp备案信息
const security = "公网安备信息"; // 公网安备信息

//上传文件大小限制，需要先修改后端大小限制
const maxImgSize = 5;//上传图片最大大小(单位M)
const maxVideoSize = 500;//上传视频最大大小(单位M)
```