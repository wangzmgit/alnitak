# Alnitak弹幕视频网站

Alnitak是一个基于nuxt和go开发的前后端分离的弹幕视频网站。
项目实现了视频、专栏、弹幕、评论、点赞、收藏等功能。 [演示视频](https://www.bilibili.com/video/BV1E6421Z7Zd/)

## 项目文档

`开发交流群: 909847398`

`文档：` https://alnitak.interastral-peace.com/

## 特性
* **弹幕播放器**: 基于DPlayer二次开发出弹幕播放器WPlayer并上传npm。
* **视频转码**：实现了视频码率和分辨率的调整，并将视频转码为HLS格式，以提高播放效果
* **实时通信**：基于WebSocket实现了私信的实时收发功能，同时实现了视频在线人数的实时更新
* **鉴权**：基于JWT实现双Token方案，通过对axios二次封装实现Token过期重发请求功能
* **滑块拼图库**：使用Go语言开发了滑块拼图生成库，并将其发布到GitHub上作为项目依赖使用
* **权限管理**：后台管理系统基于JWT和Casbin实现了基于角色的访问控制（RBAC），保障了系统的安全性和灵活性

## 相关项目

- wplayer-next弹幕播放器(https://github.com/wangzmgit/wplayer-next)
- GO滑块拼图生成(https://github.com/wangzmgit/jigsaw)

## 贡献

非常欢迎您参与项目的维护。如有任何意见或建议，请通过创建 Issue 或提交 Pull Request 的方式告知我们。
在提交之前，请务必阅读[贡献指南](https://alnitak.interastral-peace.com/guide/other/contribution)。


## 声明
本项目仅供学习和研究使用，请勿将本项目的任何内容用于商业或非法目的，否则后果自负。

## 项目截图
|                                 Web端                                 |                                Web端                                |
| :-------------------------------------------------------------------: | :-----------------------------------------------------------------: |
|   ![Web端登录](https://alnitak.interastral-peace.com/web_login.png)   |  ![Web端首页](https://alnitak.interastral-peace.com/web_home.png)   |
|  ![Web端上传](https://alnitak.interastral-peace.com/web_upload.png)   |  ![Web端视频](https://alnitak.interastral-peace.com/web_video.png)  |
| ![Web端个人中心](https://alnitak.interastral-peace.com/web_space.png) | ![Web端消息](https://alnitak.interastral-peace.com/web_message.png) |

|                                  后台管理端                                  |                                  后台管理端                                   |
| :--------------------------------------------------------------------------: | :---------------------------------------------------------------------------: |
|  ![后台管理端登录](https://alnitak.interastral-peace.com/manage_login.png)   | ![后台管理端视频管理](https://alnitak.interastral-peace.com/manage_video.png) |
| ![后台管理端用户管理](https://alnitak.interastral-peace.com/manage_user.png) |                                                                               |

|                                移动端                                |                                移动端                                 |
| :------------------------------------------------------------------: | :-------------------------------------------------------------------: |
| ![移动端首页](https://alnitak.interastral-peace.com/mobile_home.png) | ![移动端视频](https://alnitak.interastral-peace.com/mobile_video.png) |
