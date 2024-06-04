# 登录流程

## 整体流程
![双token登录方案](/token.png)

## 刷新token

调用[刷新TOKEN接口](/api/auth#刷新TOKEN)时会返回`RefreshToken`、`AccessToken`，同时会在Cookie中添加`user_id`和`user_id_md5`。
其中Cookie中的内容暂时没有在后端使用，如需使用请自行修改相关代码。
