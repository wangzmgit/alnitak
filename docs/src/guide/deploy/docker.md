# Docker部署

:::tip 前置条件
当准备使用 Docker 进行部署前，请确保已经成功安装了 Docker、Docker Compose 和 Node.js 环境。
:::

:::warning 重要信息
部署完成后项目后端、用户端、管理端分别运行在9000、9010、9030端口上。
如果需要统一对外端口或者配置自己的域名请参考[域名配置](/guide/deploy/domain)
:::

## 后端部署

#### 1.配置后端文件
打开 `data/config/application.yml`文件，根据其中的注释修改配置文件并保存

#### 2.构建
执行以下命令
```
docker-compose build
```

#### 3.启动后端
执行以下命令
```
docker-compose up -d
```

如需停止项目执行
```
docker-compose stop 
```

## 前端部署

:::tip 前置条件
前端需要使用yarn，如果尚未安装yarn，可以使用命令 `npm install -g yarn` 来安装yarn。 
:::

### 用户端部署

:::tip 打包说明
用户端项目打包时性能消耗较大，建议选择内存大于4GB的机器进行打包。如果打包和部署不在同一台机器上，可使用`yarn build`命令完成打包，
然后将`.output`文件复制到`web/web-client`目录下即可。
:::

#### 1. 配置项目
配置文件位于目录`web/web-client/src/utils/global-config.ts`，配置文件内容如下：
```js
const title = "弹幕网站标题"; // 网站标题
const https = false; // 是否使用https
const domain = "localhost:9000"; // 后端地址
const icp = "icp备案信息"; // icp备案信息
const security = "公网安备信息"; // 公网安备信息
const keywords = "视频,弹幕"; // 网站关键词
const description = "这里是介绍"; // 网站介绍

//上传文件大小限制，需要先修改后端大小限制
const maxImgSize = 5;//上传图片最大大小(单位M)
const maxVideoSize = 500;//上传视频最大大小(单位M)
```
#### 2. 构建项目
回到`web/web-client`目录下，使用一下命令构建项目
```sh
# 先安装项目依赖
yarn install
# 然后对项目进行打包
yarn build
# 最后使用docker部署项目
docker build -t alnitak-web .
```

#### 3. 启动项目
```sh
docker run -itd --name alnitakWeb -p 9010:9010 alnitak-web
```

### 管理端部署
:::tip 建议
如果你已经熟悉Nginx的配置，建议直接采取[手动部署](/guide/deploy/manual#部署管理端)来部署管理端，
因为额外运行一个Docker和Nginx会增加性能开销。
:::

#### 1. 配置项目
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
#### 2. 构建项目
回到`web/admin-client`目录下，使用一下命令构建项目
```sh
# 先安装项目依赖
yarn install
# 然后对项目进行打包
yarn build
# 最后使用docker部署项目
docker build -t alnitak-admin .
```

#### 3. 启动项目
```sh
docker run -itd --name alnitakAdmin -p 9030:9030 alnitak-admin
```


