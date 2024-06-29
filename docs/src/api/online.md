# 在线人数相关接口

## 视频在线人数连接

#### 请求URL
- `/api/v1/online/video?vid=视频ID&clientId=随机生成的客户端ID `

#### 返回示例 
```json
{
  "number": 1
}
```

#### 备注
请求时使用ws或者wss协议