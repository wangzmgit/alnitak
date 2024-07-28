# 手动部署

:::tip 前置条件
手动部署前，请确保已经成功安装了Nginx、Go、Node.js等环境，具体请参考[环境配置](/guide/start/env)
:::

:::warning 重要信息
部署完成后项目后端、用户端分别运行在9000、9010端口上。
如果需要统一对外端口或者配置自己的域名请参考[域名配置](/guide/deploy/domain)
:::

## 后端部署

#### 1. 配置项目
按照[启动后端项目](/guide/start/start)中的步骤尝试启动后端项目
#### 2. 构建项目
在`server`目录下使用以下命令构建项目
```sh
go build -o alnitak cmd/main.go
```
#### 3. 启动项目
使用以下命令启动项目
```sh
nohup ./alnitak >logs/nohup.log 2>&1
```

## 前端部署

### 部署用户端

:::tip 提示
对于用户端项目，建议使用Docker进行部署，具体内容请参考[Docker用户端部署](/guide/deploy/docker#用户端部署)。

如果无法使用Docker进行部署，我们也提供了PM2的部署方案。
:::

#### 1.项目配置
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

#### 2.启动和构建项目
* 进入到`web/web-client`目录
* 首次运行需要执行`yarn install` 安装项目依赖
* 使命令行输入`yarn dev`启动项目并检查功能是否正常
* 使命令行输入`yarn build`构建项目
* 将生成`.output`文件复制到服务器的`web/web-client`目录下

#### 3.使用以下命令在服务器上安装pm2
```sh
# 需要先在服务器上安装Node.js
npm i pm2 -g
```

#### 4.使用以下命令启动服务
```sh
pm2 start ecosystem.config.js
```

#### 常用的pm2命令
```sh
# 查看所有pm2进程
pm2 list 
# 启动进程
pm2 start 
# 停止某个进程
pm2 stop 进程号 
# 删除某个进程
pm2 delete 进程号 
# 重启某个进程
pm2 reload 进程号 
```
<!-- 参考自 https://juejin.cn/post/7224435010072346683 -->


### 部署管理端

#### 1.项目配置
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

#### 2.启动和构建项目
* 进入到`web/admin-client`目录
* 首次运行需要执行`yarn install` 安装项目依赖
* 使命令行输入`yarn dev`启动项目并检查功能是否正常
* 使命令行输入`yarn build`构建项目
* 将生成`admin`文件复制到服务器的`/usr/share/nginx/html`目录下

#### 3.启动服务
:::warning 重要提示
该步骤将在9030端口启动管理端服务。

**如果需要统一对外端口或配置自定义域名，请直接前往[域名配置](/guide/deploy/domain)中进行统一配置。**
:::

* 在`/etc/nginx/conf.d/`目录新建`admin.conf`文件，内容如下：
```
server {
    listen       9030;
	server_name  localhost; 
	client_max_body_size 1024M;

    location / {
        rewrite ^ /admin/index.html permanent;
    }

    #后台管理
    location /admin/ {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ @admin;
    }

    # 解决后台管理history路由问题
    location @admin {
        rewrite ^.*$ /admin/index.html;
    }
}
```
* 使用以下命令重启nginx

```sh
nginx -s reload
```

### 部署移动端

#### 1.项目配置
配置文件位于目录`web/mobile-client/src/utils/global-config.ts`，配置文件内容如下：
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

#### 2.启动和构建项目
* 进入到`web/mobile-client`目录
* 首次运行需要执行`yarn install` 安装项目依赖
* 使命令行输入`yarn dev`启动项目并检查功能是否正常
* 使命令行输入`yarn build`构建项目
* 将生成`mobile`文件复制到服务器的`/usr/share/nginx/html`目录下

#### 3.启动服务
:::warning 重要提示
该步骤将在9020端口启动管理端服务。

**如果需要统一对外端口或配置自定义域名，请直接前往[域名配置](/guide/deploy/domain)中进行统一配置。**
:::

* 在`/etc/nginx/conf.d/`目录新建`mobile.conf`文件，内容如下：
```
server {
    listen       9020;
	server_name  localhost; 
	client_max_body_size 1024M;

    location / {
        rewrite ^ /mobile/index.html permanent;
    }

    #后台管理
    location /mobile/ {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ @mobile;
    }

    # 解决后台管理history路由问题
    location @mobile {
        rewrite ^.*$ /mobile/index.html;
    }
}
```
* 使用以下命令重启nginx

```sh
nginx -s reload
```


