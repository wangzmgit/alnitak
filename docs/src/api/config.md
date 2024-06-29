# 配置相关接口

## 获取邮箱配置

#### 请求URL
- `/api/v1/config/getEmailConfig`
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "config": {
      "user": "",
      "host": "",
      "port": 465,
      "addresser": "",
    },
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| user      | string | 邮箱发送者     |
| host      | string | SMTP服务器地址 |
| port      | int    | SMTP服务器端口 |
| addresser | string | 发送者的邮箱   |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 修改邮箱配置

#### 请求URL
- `/api/v1/config/setEmailConfig`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名    | 必选 | 类型   | 说明           |
| :-------- | :--- | :----- | -------------- |
| user      | 是   | string | 邮箱发送者     |
| host      | 是   | string | SMTP服务器地址 |
| port      | 是   | int    | SMTP服务器端口 |
| addresser | 是   | string | 发送者的邮箱   |
| pass      | 否   | string | 邮箱秘钥       |

#### 返回示例 
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取存储配置

#### 请求URL
- `/api/v1/config/getStorageConfig`
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "config": {
      "maxImgSize": 5,
      "maxVideoSize": 500,
      "type": "",
      "keyId": "",
      "bucket": "",
      "endpoint": "",
      "domain": "",
      "uploadMp4File": false,
    },
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名        | 类型   | 说明                                                   |
| :------------ | :----- | ------------------------------------------------------ |
| maxImgSize    | int    | 上传图片的大小限制                                     |
| maxVideoSize  | int    | 上传视频的大小限制                                     |
| type          | string | 对象存储类型，目前仅支持本地(local)和阿里云OSS(aliyun) |
| keyId         | string | 对应阿里云access_key_secret                            |
| bucket        | string | 对象存储的bucket                                       |
| endpoint      | string | 阿里云OSS的endpoint                                    |
| domain        | string | 自定义域名                                             |
| uploadMp4File | bool   | 上传对象存储时是否上传原始mp4文件                      |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 修改存储配置

#### 请求URL
- `/api/v1/config/setStorageConfig`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名        | 必选 | 类型   | 说明                                                   |
| :------------ | :--- | :----- | ------------------------------------------------------ |
| maxImgSize    | 是   | int    | 上传图片的大小限制                                     |
| maxVideoSize  | 是   | int    | 上传视频的大小限制                                     |
| type          | 是   | string | 对象存储类型，目前仅支持本地(local)和阿里云OSS(aliyun) |
| keyId         | 否   | string | 对应阿里云access_key_secret                            |
| bucket        | 否   | string | 对象存储的bucket                                       |
| endpoint      | 否   | string | 阿里云OSS的endpoint                                    |
| domain        | 否   | string | 自定义域名                                             |
| uploadMp4File | 否   | bool   | 上传对象存储时是否上传原始mp4文件                      |


#### 返回示例 
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取其他配置

#### 请求URL
- `/api/v1/config/getOtherConfig`

#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "config": {
      "allowOrigin": "",
      "prefix": "",
    },
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名      | 类型   | 说明                           |
| :---------- | :----- | ------------------------------ |
| allowOrigin | string | 跨域配置，存在多个使用`,`分隔  |
| prefix      | string | 用户注册时生成用户名的默认前缀 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 修改其他配置

#### 请求URL
- `/api/v1/config/setOtherConfig`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名      | 必选 | 类型   | 说明                           |
| :---------- | :--- | :----- | ------------------------------ |
| allowOrigin | 是   | string | 跨域配置，存在多个使用`,`分隔  |
| prefix      | 否   | string | 用户注册时生成用户名的默认前缀 |

#### 返回示例 
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无