# 常见问题解答

## 端口修改
为了方便一键部署，项目没有提供配置端口的功能。如果出于某些原因需要修改端口号，
您可以前往 `server/internal/routes/router.go `文件中修改`port = "9000"`中的端口号为您想要的值。
**在修改任何与端口相关的内容之前，请确保您已经熟悉了部署流程，并清楚了解端口变更可能会对哪些内容造成影响。**

## 修改默认管理员密码
使用默认账号(账号：admin@admin.com 密码: 123456)登录后台管理，
修改邮箱为自己的邮箱，进入到前端登录页面，点击找回密码修改密码

## 轮播图不能正常显示
系统的轮播图是按照分区进行上传的。若某一分区的轮播图数量少于 3 张，则可能无法正常展示。

## 首页视频空白
首页的视频列表数据并非实时数据，系统每隔 3 小时会统计近期热门视频并展示在首页。若需立即将视频展示在首页，可以重启后端项目，每次重启时也会重新计算首页数据。

## 视频显示审核中但后台管理中不存在该视频
用户上传视频后，需要对视频进行转码等操作，视频转码需要一定的时间。在视频转码完成前，后台管理系统中不会显示该视频。